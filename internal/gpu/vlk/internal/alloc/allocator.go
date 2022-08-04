package alloc

import (
	"fmt"

	"github.com/vulkan-go/vma"
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type (
	Allocator struct {
		logger           config.Logger
		ref              *vma.Allocator
		allocatedBuffers []allocatedBuffer
	}

	allocatedBuffer struct {
		buffer vulkan.Buffer
		alloc  vma.Allocation
	}
)

func NewAllocator(logger config.Logger, inst *instance.Instance, pd *physical.Device, ld *logical.Device) *Allocator {
	ref, err := vma.NewAllocator(&vma.AllocatorCreateInfo{
		VulkanProcAddr:              inst.ProcAddr(),
		Instance:                    inst.Ref(),
		PhysicalDevice:              pd.PrimaryGPU().Ref,
		Device:                      ld.Ref(),
		PreferredLargeHeapBlockSize: 0,   // todo:?
		FrameInUseCount:             0,   // todo:?
		HeapSizeLimit:               nil, // todo:?
		VulkanAPIVersion:            def.VKApiVersion,
	})
	if err != nil {
		panic(fmt.Errorf("failed create vulkan memory allocator: %w", err))
	}

	logger.Debug("memory allocator created")

	return &Allocator{
		logger: logger,
		ref:    ref,
	}
}

func (a *Allocator) Free() {
	for _, obj := range a.allocatedBuffers {
		a.ref.DestroyBuffer(obj.buffer, obj.alloc)
	}

	a.ref.Destroy()
	a.logger.Debug("freed: memory allocator")
}

func (a *Allocator) CreateVertexBuffer() (vulkan.Buffer, vma.Allocation) {
	return a.createBuffer(def.BufferVertexSizeBytes, vulkan.BufferUsageFlags(
		vulkan.BufferUsageVertexBufferBit,
	))
}

func (a *Allocator) createBuffer(size uint64, usage vulkan.BufferUsageFlags) (vulkan.Buffer, vma.Allocation) {
	info := &vulkan.BufferCreateInfo{
		SType:       vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(size),
		Usage:       usage,
		SharingMode: vulkan.SharingModeExclusive,
	}

	allocInfo := &vma.AllocationCreateInfo{
		Usage: vma.MemoryUsageCPUToGPU,
	}

	buffer, alloc, _, err := a.ref.CreateBuffer(info, allocInfo, false)
	if err != nil {
		panic(fmt.Errorf("failed alloc vulkan buffer: %w", err))
	}

	a.allocatedBuffers = append(a.allocatedBuffers, allocatedBuffer{
		buffer: buffer,
		alloc:  alloc,
	})

	return buffer, alloc
}
