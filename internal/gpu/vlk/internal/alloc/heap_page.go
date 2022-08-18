package alloc

import (
	"fmt"
	"math"
)

const garbageListCapacity = 32

type (
	heapCtl interface {
		createBuffer(pageID pageID, size uint32) bufferID
		destroyBuffer(buffID bufferID)
		writeAt(buffID bufferID, offset uint32, data []byte)
	}
)

// Responsibility:
// - create new areas when needed
// - physically map vulkan buffers to logical areas
// - contain list of flags and requirements
type h3Page struct {
	id                  pageID
	ctl                 heapCtl
	defaultAreaCapacity uint32
	defaultAreaAlign    uint32
	isImmutable         bool
	autoGarbageCollect  bool

	garbageList [][2]uint32 // list of garbage pair (buffID, allocID), that should be cleaned each GC call
	areas       map[bufferID]*h3Area
}

func newH3Page(
	id pageID,
	ctl heapCtl,
	defaultAreaCapacity uint32,
	defaultAreaAlign uint32,
	isImmutable bool,
	autoGarbageCollect bool,
) *h3Page {
	return &h3Page{
		id:                  id,
		ctl:                 ctl,
		defaultAreaCapacity: defaultAreaCapacity,
		defaultAreaAlign:    defaultAreaAlign,
		isImmutable:         isImmutable,
		autoGarbageCollect:  autoGarbageCollect,

		areas:       make(map[bufferID]*h3Area),
		garbageList: make([][2]uint32, 0, garbageListCapacity),
	}
}

func (p *h3Page) garbageCollect() {
	if !p.autoGarbageCollect {
		return
	}

	for _, pair := range p.garbageList {
		p.free(bufferID(pair[0]), allocID(pair[1]))
	}

	p.garbageList = make([][2]uint32, 0, garbageListCapacity)
}

func (p *h3Page) write(data []byte) (bufferID, allocID) {
	size := uint32(len(data))
	area, buffID := p.areaThatCanFit(size)

	// mark logical area space as claimed
	node, ok := area.claim(size)
	if !ok {
		// this is not possible, because we check size in areaThatCanFit
		panic(fmt.Errorf("logical area not have space to write"))
	}

	if p.autoGarbageCollect {
		// node space will be freed automatic in next GC tick
		p.garbageList = append(p.garbageList, [2]uint32{uint32(buffID), node.offset})
	}

	// write data to real buffer
	p.ctl.writeAt(buffID, node.offset, data)

	return buffID, allocID(node.offset)
}

func (p *h3Page) free(buffID bufferID, allocationID allocID) {
	if p.isImmutable {
		panic(fmt.Errorf("cannot free immutable data block at %d:%d", buffID, allocationID))
	}

	area, exist := p.areas[buffID]
	if !exist {
		panic(fmt.Errorf("not found area for buffID %d in page", buffID))
	}

	_, freed := area.free(uint32(allocationID))
	if !freed {
		panic(fmt.Errorf("failed free page memory at %d:%d (maybe invalid offset?)", buffID, allocationID))
	}
}

func (p *h3Page) areaThatCanFit(size uint32) (*h3Area, bufferID) {
	// get first area, that can fit this data
	for buffID, area := range p.areas {
		if area.freeSize() > size {
			return area, buffID
		}
	}

	// if not area found, we need to extend page memory
	// be creating new area
	alignedSize := uint32(math.Ceil(float64(size)/float64(p.defaultAreaAlign)) * float64(p.defaultAreaAlign))
	capacity := p.max(p.defaultAreaCapacity, alignedSize)

	buffID := p.ctl.createBuffer(p.id, capacity)
	area := newArea(capacity, p.defaultAreaAlign)
	p.areas[buffID] = area

	// return created area
	return area, buffID
}

func (p *h3Page) max(a, b uint32) uint32 {
	if a > b {
		return a
	}

	return b
}
