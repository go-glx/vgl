package alloc

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/shared/vlkext"
)

type (
	bufferID uint32

	Allocator struct {
		logger vlkext.Logger
		inst   *instance.Instance
		pd     *physical.Device
		ld     *logical.Device
		pool   *command.Pool

		internalBufferLastID bufferID
		allocatedBuffers     map[bufferID]internalBuffer
	}

	internalBuffer struct {
		id       bufferID
		ref      vulkan.Buffer
		memory   vulkan.DeviceMemory
		capacity vulkan.DeviceSize
	}
)

func NewAllocator(
	logger vlkext.Logger,
	inst *instance.Instance,
	pd *physical.Device,
	ld *logical.Device,
	pool *command.Pool,
) *Allocator {
	return &Allocator{
		logger: logger,
		inst:   inst,
		pd:     pd,
		ld:     ld,
		pool:   pool,

		internalBufferLastID: 0,
		allocatedBuffers:     make(map[bufferID]internalBuffer),
	}
}

func (a *Allocator) Free() {
	for _, buff := range a.allocatedBuffers {
		a.destroyBuffer(buff)
	}

	a.logger.Debug("freed: memory allocator")
}

func (a *Allocator) destroyBuffer(buff internalBuffer) {
	vulkan.DestroyBuffer(a.ld.Ref(), buff.ref, nil)
	vulkan.FreeMemory(a.ld.Ref(), buff.memory, nil)

	delete(a.allocatedBuffers, buff.id)
	a.logger.Debug(fmt.Sprintf("freed: buffer %d", buff.id))
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

func (a *Allocator) copyBuffer(src internalBuffer, dst internalBuffer, srcOffset, dstOffset, size uint32) {
	a.pool.TemporaryBuffer(func(cb vulkan.CommandBuffer) {
		copyRegion := vulkan.BufferCopy{
			SrcOffset: vulkan.DeviceSize(srcOffset),
			DstOffset: vulkan.DeviceSize(dstOffset),
			Size:      vulkan.DeviceSize(size),
		}

		vulkan.CmdCopyBuffer(cb, src.ref, dst.ref, 1, []vulkan.BufferCopy{copyRegion})

		a.logger.Debug(fmt.Sprintf("buffer data copied (%d->%d) offsets=[src=%d, dst=%d], size=%.2fKB",
			src.id,
			dst.id,
			srcOffset,
			dstOffset,
			float32(size)/1024,
		))
	})
}

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
