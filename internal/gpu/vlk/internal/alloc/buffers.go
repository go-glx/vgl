package alloc

import (
	"github.com/vulkan-go/vulkan"
)

type Chunk struct {
	Buffer         vulkan.Buffer
	BufferOffset   uint64
	InstancesCount uint32
	IndexCount     uint32
}

type Buffers struct {
	alloc  *Allocator
	heapV1 *bufferHeap // todo: rewrite to v2
	heapV2 *heapV2     // todo: rename to heap, after v1 refactored
}

func NewBuffers(alloc *Allocator) *Buffers {
	return &Buffers{
		alloc:  alloc,
		heapV1: newHeapV1(),
		heapV2: newHeapV2(),
	}
}
