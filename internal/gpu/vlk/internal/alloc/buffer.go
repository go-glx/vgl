package alloc

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type internalBuffer struct {
	ref      vulkan.Buffer
	dataPtr  unsafe.Pointer
	memory   vulkan.DeviceMemory
	capacity vulkan.DeviceSize
}

func (a *Allocator) createBuffer(size int, buffType vulkan.BufferUsageFlags) internalBuffer {
	// create new buffer page
	info := &vulkan.BufferCreateInfo{
		SType:       vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(size),
		Usage:       buffType,
		SharingMode: vulkan.SharingModeExclusive,
	}

	var buffer vulkan.Buffer
	must.Work(vulkan.CreateBuffer(a.ld.Ref(), info, nil, &buffer))

	// get device memory requirements for it
	var memoryReq vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(a.ld.Ref(), buffer, &memoryReq)
	memoryReq.Deref()

	memoryTypeIndex := findVertexBufferMemoryType(
		a.pd,
		memoryReq,
		vulkan.MemoryPropertyFlags(
			vulkan.MemoryPropertyHostVisibleBit|
				vulkan.MemoryPropertyHostCoherentBit,
		),
	)

	memAllocInfo := &vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memoryReq.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	var bufferMemory vulkan.DeviceMemory
	must.Work(vulkan.AllocateMemory(a.ld.Ref(), memAllocInfo, nil, &bufferMemory))

	vulkan.BindBufferMemory(a.ld.Ref(), buffer, bufferMemory, 0)

	var data unsafe.Pointer
	vulkan.MapMemory(a.ld.Ref(), bufferMemory, 0, info.Size, 0, &data)

	internalBuff := internalBuffer{
		ref:      buffer,
		dataPtr:  data,
		memory:   bufferMemory,
		capacity: info.Size,
	}

	a.logger.Debug(fmt.Sprintf("Buffer %.3fMB capacity - allocated", float64(info.Size/1024)))
	a.allocatedBuffers = append(a.allocatedBuffers, internalBuff)

	return internalBuff
}

func findVertexBufferMemoryType(pd *physical.Device, memoryReq vulkan.MemoryRequirements, memFlags vulkan.MemoryPropertyFlags) uint32 {
	typeFilter := memoryReq.MemoryTypeBits

	var memProperties vulkan.PhysicalDeviceMemoryProperties
	vulkan.GetPhysicalDeviceMemoryProperties(pd.PrimaryGPU().Ref, &memProperties)
	memProperties.Deref()

	for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
		memType := memProperties.MemoryTypes[i]
		memType.Deref()

		if (typeFilter&(1<<i) != 0) && ((memType.PropertyFlags & memFlags) == memFlags) {
			return i
		}
	}

	panic(fmt.Errorf("failed find suitable GPU memory for vertex buffer"))
}
