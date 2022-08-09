package alloc

import (
	"github.com/vulkan-go/vulkan"
)

const (
	memoryRegionIDEmpty = 0
)

type (
	pageID uint32

	heapV2 struct {
		lastAllocationID AllocationID
		lastPageID       pageID
		allocations      map[AllocationID]heapPagePointer
		exclusive        map[pageID]internalBuffer
	}

	heapPagePointer struct {
		isExclusive bool
		pageID      pageID
		allocRange  Range
	}
)

func newHeapV2() *heapV2 {
	return &heapV2{
		allocations: make(map[AllocationID]heapPagePointer),
		exclusive:   make(map[pageID]internalBuffer),
	}
}

func (b *Buffers) GetMemoryPointer(id AllocationID) Allocation {
	if allocPtr, exist := b.heapV2.allocations[id]; exist {
		if allocPtr.allocRange.Size == 0 {
			return Allocation{
				HasData: false,
			}
		}

		var buff vulkan.Buffer
		if allocPtr.isExclusive {
			buff = b.heapV2.exclusive[allocPtr.pageID].ref
		}

		return Allocation{
			HasData: true,
			Buffer:  buff,
			Range:   allocPtr.allocRange,
		}
	}

	return Allocation{
		HasData: false,
	}
}

func (b *Buffers) AllocateIndexMemory(data []byte) AllocationID {
	return b.allocateExclusiveMemory(data, vulkan.BufferUsageIndexBufferBit)
}

func (b *Buffers) allocateExclusiveMemory(data []byte, dstType vulkan.BufferUsageFlagBits) AllocationID {
	bufferSize := uint32(len(data))
	if bufferSize == 0 {
		return memoryRegionIDEmpty
	}

	// create tmp buffer, visible from CPU/GPU side
	tmpBuffer := b.alloc.createBuffer(
		bufferSize,
		vulkan.BufferUsageTransferSrcBit,
		memRequirementHostVisible,
	)

	// copy data to it
	vulkan.Memcopy(tmpBuffer.dataPtr, data)
	vulkan.UnmapMemory(b.alloc.ld.Ref(), tmpBuffer.memory)

	// create fast-persistence buffer, visible only from GPU side
	fastBuffer := b.alloc.createBuffer(
		bufferSize,
		vulkan.BufferUsageTransferDstBit|dstType,
		memRequirementFastGPU,
	)

	// copy data from tmp to fast device memory
	b.alloc.copyBuffer(tmpBuffer, fastBuffer, 0, 0, bufferSize)

	// clean tmp data, drop tmp buffers, etc..
	b.alloc.destroyBuffer(tmpBuffer)

	// store allocated buffer and return reference to it
	return b.storeAllocated(fastBuffer, Range{
		// exclusive memory use all buffer space
		PositionFrom: 0,
		Size:         bufferSize,
	})
}

func (b *Buffers) storeAllocated(buff internalBuffer, memRange Range) AllocationID {
	// write mem page
	b.heapV2.lastPageID++
	b.heapV2.exclusive[b.heapV2.lastPageID] = buff

	// write pointer to this page
	b.heapV2.lastAllocationID++
	b.heapV2.allocations[b.heapV2.lastAllocationID] = heapPagePointer{
		isExclusive: true,
		pageID:      b.heapV2.lastPageID,
		allocRange:  memRange,
	}

	// return pointer ID
	return b.heapV2.lastAllocationID
}
