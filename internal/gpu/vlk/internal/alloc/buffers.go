package alloc

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type Chunk struct {
	Buffer         vulkan.Buffer
	BufferOffset   uint64
	InstancesCount uint32
	IndexCount     uint32
}

type Buffers struct {
	alloc *Allocator
	heap  *bufferHeap
}

type (
	bufferHeap struct {
		pages              []*bufferPage
		shaderIndexBuffers map[string]internalBuffer
	}

	bufferPage struct {
		buff   internalBuffer
		staged []byte
	}
)

func NewBuffers(alloc *Allocator) *Buffers {
	return &Buffers{
		alloc: alloc,
		heap: &bufferHeap{
			pages:              make([]*bufferPage, 0),
			shaderIndexBuffers: make(map[string]internalBuffer),
		},
	}
}

func (b *Buffers) Clear() {
	for _, page := range b.heap.pages {
		page.staged = make([]byte, 0, page.buff.capacity)
	}
}

func (b *Buffers) IndexBufferOf(shader *shader.Shader, indexes []uint16) vulkan.Buffer {
	if buff, exist := b.heap.shaderIndexBuffers[shader.Meta().ID()]; exist {
		return buff.ref
	}

	newBuff := b.alloc.createIndexBuffer()

	staging := make([]byte, 0, len(indexes)*2)
	for _, index := range indexes {
		staging = append(staging, uint8(index&0xff), uint8(index>>8))
	}

	vulkan.Memcopy(newBuff.dataPtr, staging)
	b.heap.shaderIndexBuffers[shader.Meta().ID()] = newBuff
	return newBuff.ref
}

func (b *Buffers) Write(instances []shader.InstanceData) []Chunk {
	if len(instances) == 0 {
		return nil
	}

	chunks := make([]Chunk, 0)
	instSize := uint64(len(instances[0].BindingData()))
	instIndexCount := uint32(len(instances[0].Indexes()))

	instCount := uint32(0)

	var page *bufferPage
	currentPageID := 0
	currentPageOffset := uint64(0)
	currentChunkPageOffsetSpecified := false

	// todo: capacity is max, do not copy all 64mb every frame
	// for one instance in 12 bytes..

	for _, inst := range instances {
		if exist := currentPageID <= len(b.heap.pages)-1; !exist {
			b.extendHeap()
		}

		page = b.heap.pages[currentPageID]
		freeSpace := uint64(cap(page.staged) - len(page.staged))

		if !currentChunkPageOffsetSpecified {
			currentPageOffset = uint64(len(page.staged))
			currentChunkPageOffsetSpecified = true
		}

		if instSize > freeSpace {
			// add chunk and reset counters
			chunks = append(chunks, Chunk{
				Buffer:         page.buff.ref,
				BufferOffset:   currentPageOffset,
				InstancesCount: instCount,
				IndexCount:     instIndexCount,
			})
			instCount = 0

			// extend heap and reset pointers
			currentPageID++
			if exist := currentPageID <= len(b.heap.pages)-1; !exist {
				b.extendHeap()
			}

			page = b.heap.pages[currentPageID]
			freeSpace = uint64(cap(page.staged) - len(page.staged))
			currentPageOffset = 0
		}

		// stage data to current page
		instCount++
		page.staged = append(page.staged, inst.BindingData()...)
	}

	chunks = append(chunks, Chunk{
		Buffer:         page.buff.ref,
		BufferOffset:   currentPageOffset,
		InstancesCount: instCount,
		IndexCount:     instIndexCount,
	})

	return chunks
}

func (b *Buffers) Flush() {
	for _, page := range b.heap.pages {
		if len(page.staged) <= 0 {
			break
		}

		vulkan.Memcopy(page.buff.dataPtr, page.staged)
	}
}

func (b *Buffers) extendHeap() {
	newBuffer := b.alloc.createVertexBuffer()

	b.heap.pages = append(b.heap.pages, &bufferPage{
		buff:   newBuffer,
		staged: make([]byte, 0, newBuffer.capacity),
	})
}
