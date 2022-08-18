package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

const memRequirementHostVisible = vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit
const memRequirementFastGPU = vulkan.MemoryPropertyDeviceLocalBit

type internalBuffer struct {
	id       bufferID
	ref      vulkan.Buffer
	memory   vulkan.DeviceMemory
	capacity vulkan.DeviceSize
}

func (a *Allocator) createBuffer(size uint32, buffType vulkan.BufferUsageFlagBits, memoryFlags vulkan.MemoryPropertyFlagBits) internalBuffer {
	// create new buffer page
	info := &vulkan.BufferCreateInfo{
		SType:       vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(size),
		Usage:       vulkan.BufferUsageFlags(buffType),
		SharingMode: vulkan.SharingModeExclusive,
	}

	var buffer vulkan.Buffer
	must.Work(vulkan.CreateBuffer(a.ld.Ref(), info, nil, &buffer))

	// get device memory requirements for it
	var memoryReq vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(a.ld.Ref(), buffer, &memoryReq)
	memoryReq.Deref()

	memoryTypeIndex := findBufferWithMemoryType(
		a.pd,
		memoryReq,
		vulkan.MemoryPropertyFlags(memoryFlags),
	)

	memAllocInfo := &vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memoryReq.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	var bufferMemory vulkan.DeviceMemory
	must.Work(vulkan.AllocateMemory(a.ld.Ref(), memAllocInfo, nil, &bufferMemory))

	vulkan.BindBufferMemory(a.ld.Ref(), buffer, bufferMemory, 0)

	a.internalBufferLastID++
	internalBuff := internalBuffer{
		id:       a.internalBufferLastID,
		ref:      buffer,
		memory:   bufferMemory,
		capacity: info.Size,
	}

	a.allocatedBuffers[internalBuff.id] = internalBuff
	a.logger.Debug(fmt.Sprintf("Buffer %d with %.3fKB capacity - allocated",
		internalBuff.id,
		float64(info.Size/1024)),
	)

	return internalBuff
}

func findBufferWithMemoryType(pd *physical.Device, memoryReq vulkan.MemoryRequirements, memFlags vulkan.MemoryPropertyFlags) uint32 {
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
