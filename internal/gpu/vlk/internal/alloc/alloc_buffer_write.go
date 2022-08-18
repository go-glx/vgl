package alloc

import (
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

func (a *Allocator) writeBuffer(buff internalBuffer, offset uint32, data []byte) {
	// map memory region
	var ptr unsafe.Pointer
	vulkan.MapMemory(
		a.ld.Ref(),
		buff.memory,
		vulkan.DeviceSize(offset),
		vulkan.DeviceSize(len(data)),
		0,
		&ptr,
	)

	// coherent is host memory
	// we can just copy data directly to it
	vulkan.Memcopy(ptr, data)

	// unmap host memory
	vulkan.UnmapMemory(a.ld.Ref(), buff.memory)
}
