package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

const (
	BufferTypeVertex  BufferType = iota // will be mapped to GPU with vertex buffer specifications
	BufferTypeIndex                     // will be mapped to GPU with vertex index specifications
	BufferTypeUniform                   // will be mapped to GPU with vertex uniform specifications
)

const (
	StorageTargetCoherent  StorageTarget = iota // all allocations will be written only on MEM
	StorageTargetWritable                       // all allocations will be written to MEM and copied to device automatically
	StorageTargetImmutable                      // all allocations will be written to device once. Cannot be changes, only deleted
)

const (
	FlagsNone      Flags = 1 << iota // not use any special logic
	FlagsTemporary                   // memory page will be freed automatic at next frame
)

type (
	pageID  uint32 // unique page ID (each page has own ID)
	allocID uint32 // its area offset (node ID)

	BufferType    uint8
	StorageTarget uint8
	Flags         uint8

	Heap struct {
		allocator      *Allocator                  // physical binding to vulkan
		nextPageID     pageID                      // counter for inc pageID
		pages          map[pageID]*h3Page          // allocated logical pages
		features       map[pageFeatures]pageID     // pages features
		featuresPtr    map[pageID]pageFeatures     // pages features (back ptr)
		buffers        map[bufferID]internalBuffer // allocated buffers from this heap
		stagingBuffers map[bufferID]internalBuffer // staging buffer for bufferID
		bufferOwner    map[bufferID]pageID         // page that own buffer with this ID
	}

	Allocation struct {
		Valid  bool
		Buffer vulkan.Buffer
		Offset vulkan.DeviceSize
		Size   vulkan.DeviceSize

		pageID  pageID
		buffID  bufferID
		allocID allocID
	}

	Stats struct {
		Grouped       []GroupedStats
		TotalCapacity uint32
		TotalSize     uint32
	}

	GroupedStats struct {
		BufferType BufferType
		TotalPages uint32
		TotalAreas uint32
		Capacity   uint32
		Size       uint32
	}
)

func NewHeap(allocator *Allocator) *Heap {
	return &Heap{
		allocator:      allocator,
		nextPageID:     0,
		pages:          make(map[pageID]*h3Page),
		features:       make(map[pageFeatures]pageID),
		featuresPtr:    make(map[pageID]pageFeatures),
		buffers:        make(map[bufferID]internalBuffer),
		stagingBuffers: make(map[bufferID]internalBuffer),
		bufferOwner:    make(map[bufferID]pageID),
	}
}

// GarbageCollect should be called after every frame
// it will automatically free some unused buffer memory
// not calling this function will create memory leaks
// specially in vertex buffers
func (h *Heap) GarbageCollect() {
	for _, page := range h.pages {
		page.garbageCollect()
	}
}

func (h *Heap) Stats() Stats {
	grouped := make(map[BufferType]*GroupedStats)

	for _, page := range h.pages {
		features := h.featuresPtr[page.id]
		gStats, exist := grouped[features.bufferType]
		if !exist {
			grouped[features.bufferType] = &GroupedStats{
				BufferType: features.bufferType,
			}
			gStats = grouped[features.bufferType]
		}

		gStats.TotalPages++

		for _, area := range page.areas {
			gStats.TotalAreas++
			gStats.Capacity += area.capacity
			gStats.Size += area.size
		}
	}

	stats := Stats{}
	for _, gStats := range grouped {
		stats.TotalCapacity += gStats.Capacity
		stats.TotalSize += gStats.Size
		stats.Grouped = append(stats.Grouped, *gStats)
	}

	return stats
}

// Write will find/create buffer with specified features, automatic write
// data to free space and return Allocation object with offset, buffer ptr
// and other details
//
// Write very love big chunks of data, so good idea to use some staging buffer
// and write as many bytes as possible per one call (but try to limit data size to some
// reasonable amount, like 32MB or default buffer capacity)
//
// When FlagsTemporary used, buffer will automatic garbage collect all this data
// in next frame, so no need to manually call Free
func (h *Heap) Write(data []byte, bType BufferType, target StorageTarget, flags Flags) Allocation {
	features := pageFeatures{
		bufferType:    bType,
		storageTarget: target,
		flags:         flags,
	}

	// find page that support this features
	page := h.pageWithFeatures(features)

	// write data to page
	buffID, allocID := page.write(data)

	// find underlying vk buffer
	internalBuffer := h.bufferByID(buffID)

	// assemble allocation info
	return Allocation{
		Valid:   true,
		Buffer:  internalBuffer.ref,
		Offset:  vulkan.DeviceSize(allocID),
		Size:    vulkan.DeviceSize(len(data)),
		pageID:  page.id,
		buffID:  buffID,
		allocID: allocID,
	}
}

// Free will clean allocated memory
// Calling this function many times (or with invalid Allocation object) will
// panic
func (h *Heap) Free(alloc Allocation) {
	page, exist := h.pages[alloc.pageID]
	if !exist {
		panic(fmt.Errorf("failed free mem: page with id %d not exist in heap", alloc.pageID))
	}

	page.free(alloc.buffID, alloc.allocID)
}

func (h *Heap) bufferByID(id bufferID) internalBuffer {
	return h.buffers[id]
}

func (h *Heap) pageWithFeatures(features pageFeatures) *h3Page {
	if pageID, exist := h.features[features]; exist {
		return h.pages[pageID]
	}

	createdPageID := h.createPage(features)
	return h.pages[createdPageID]
}

func (h *Heap) createPage(features pageFeatures) pageID {
	// create ID for new page
	h.nextPageID++
	pageID := h.nextPageID

	// create page
	page := newH3Page(
		pageID,
		h,
		features.defaultBufferCapacity(),
		features.defaultBufferAlign(),
		features.storageTarget == StorageTargetImmutable,
		features.flags&FlagsTemporary != 0,
	)

	// save it for usage
	h.features[features] = pageID
	h.featuresPtr[pageID] = features
	h.pages[pageID] = page

	return pageID
}

func (h *Heap) createBuffer(pageID pageID, size uint32) bufferID {
	features, ok := h.featuresPtr[pageID]
	if !ok {
		panic(fmt.Errorf("unexpected empty features for page %d", pageID))
	}

	buff := h.allocator.createBuffer(
		size,
		features.vulkanBufferUsage(),
		features.vulkanMemoryFlags(),
	)

	h.buffers[buff.id] = buff
	h.bufferOwner[buff.id] = pageID
	return buff.id
}

func (h *Heap) destroyBuffer(buffID bufferID) {
	buff, ok := h.buffers[buffID]
	if !ok {
		panic(fmt.Errorf("unexpected empty buffer by id %d", buffID))
	}

	h.allocator.destroyBuffer(buff)
}

func (h *Heap) writeAt(buffID bufferID, offset uint32, data []byte) {
	pageID, exist := h.bufferOwner[buffID]
	if !exist {
		panic(fmt.Errorf("failed find pageID by buffID %d", buffID))
	}

	features, exist := h.featuresPtr[pageID]
	if !exist {
		panic(fmt.Errorf("failed find features of pageID %d", pageID))
	}

	buff, exist := h.buffers[buffID]
	if !exist {
		panic(fmt.Errorf("failed find buffer with ID %d", buffID))
	}

	switch features.storageTarget {
	case StorageTargetImmutable:
		h.writeImmutableToDevice(buff, offset, data)
		return
	case StorageTargetWritable:
		h.writeToDevice(buff, offset, data)
		return
	case StorageTargetCoherent:
		h.writeToCoherent(buff, offset, data)
		return
	default:
		panic(fmt.Errorf("unknown storageTarget %d", features.storageTarget))
	}
}

func (h *Heap) writeImmutableToDevice(buff internalBuffer, offset uint32, data []byte) {
	size := uint32(len(data))

	// create tmp buffer, visible from CPU/GPU side
	tmpBuffer := h.allocator.createBuffer(
		size,
		vulkan.BufferUsageTransferSrcBit,
		vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit,
	)

	// copy data to it
	h.allocator.writeBuffer(tmpBuffer, 0, data)

	// copy data from tmp to fast device memory
	h.allocator.copyBuffer(tmpBuffer, buff, 0, offset, size)

	// drop tmp buffer, because immutable mean "we not want change it again"
	// so tmp buffer not need
	h.allocator.destroyBuffer(tmpBuffer)
}

func (h *Heap) writeToDevice(buff internalBuffer, offset uint32, data []byte) {
	size := uint32(len(data))

	// get staging buffer for this real device buffer
	stagingBuffer, exist := h.stagingBuffers[buff.id]
	if !exist {
		stagingBuffer = h.allocator.createBuffer(
			size,
			vulkan.BufferUsageTransferSrcBit,
			vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit,
		)

		h.stagingBuffers[buff.id] = stagingBuffer
	}

	// copy data to it
	h.allocator.writeBuffer(stagingBuffer, 0, data)

	// copy data from src to fast device memory
	h.allocator.copyBuffer(stagingBuffer, buff, 0, offset, size)
}

func (h *Heap) writeToCoherent(buff internalBuffer, offset uint32, data []byte) {
	h.allocator.writeBuffer(buff, offset, data)
}
