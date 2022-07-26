package command

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type Pool struct {
	logger config.Logger
	pd     *physical.Device
	ld     *logical.Device

	ref     vulkan.CommandPool
	buffers []vulkan.CommandBuffer
}

func NewPool(logger config.Logger, pd *physical.Device, ld *logical.Device) *Pool {
	pool, buffers := createPool(pd, ld)
	return &Pool{
		logger:  logger,
		pd:      pd,
		ld:      ld,
		ref:     pool,
		buffers: buffers,
	}
}

func (p *Pool) Free() {
	vulkan.FreeCommandBuffers(p.ld.Ref(), p.ref, uint32(len(p.buffers)), p.buffers)
	vulkan.DestroyCommandPool(p.ld.Ref(), p.ref, nil)

	p.logger.Debug("freed: command pool")
}

func (p *Pool) BuffersCount() int {
	return len(p.buffers)
}

func (p *Pool) CommandBuffer(ind int) vulkan.CommandBuffer {
	return p.buffers[ind]
}

func createPool(pd *physical.Device, ld *logical.Device) (vulkan.CommandPool, []vulkan.CommandBuffer) {
	createInfo := &vulkan.CommandPoolCreateInfo{
		SType:            vulkan.StructureTypeCommandPoolCreateInfo,
		QueueFamilyIndex: pd.PrimaryGPU().Families.GraphicsFamilyId,
		Flags:            vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateResetCommandBufferBit),
	}

	// create pool
	var pool vulkan.CommandPool
	must.Work(vulkan.CreateCommandPool(ld.Ref(), createInfo, nil, &pool))

	// create buffers
	buffersCount := pd.PrimaryGPU().SurfaceProps.ConcurrentBuffersCount()
	allocInfo := &vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        pool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: buffersCount,
	}

	buffers := make([]vulkan.CommandBuffer, buffersCount)
	must.Work(vulkan.AllocateCommandBuffers(ld.Ref(), allocInfo, buffers))

	// ok
	return pool, buffers
}
