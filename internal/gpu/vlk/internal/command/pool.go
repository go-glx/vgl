package command

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/shared/vlkext"
)

type Pool struct {
	logger vlkext.Logger
	pd     *physical.Device
	ld     *logical.Device

	ref vulkan.CommandPool

	// main command buffers, used for draw commands
	// one buffer for each swapChain image
	mainBuffers []vulkan.CommandBuffer
}

func NewPool(logger vlkext.Logger, pd *physical.Device, ld *logical.Device) *Pool {
	pool, buffers := createPool(pd, ld)
	return &Pool{
		logger:      logger,
		pd:          pd,
		ld:          ld,
		ref:         pool,
		mainBuffers: buffers,
	}
}

func (p *Pool) Free() {
	vulkan.FreeCommandBuffers(p.ld.Ref(), p.ref, uint32(len(p.mainBuffers)), p.mainBuffers)
	vulkan.DestroyCommandPool(p.ld.Ref(), p.ref, nil)

	p.logger.Debug("freed: command pool")
}

func (p *Pool) MainBuffersCount() int {
	return len(p.mainBuffers)
}

func (p *Pool) MainCommandBuffer(ind int) vulkan.CommandBuffer {
	return p.mainBuffers[ind]
}

// TemporaryBuffer will create temporary one time command buffer
// and give it to exec for execution
// right after exec is completed, this buffer will be destroyed
// this useful for one time GPU commands, like data uploading to GPU
//
// All written commands will be automatically executed in GPU
// after exec
func (p *Pool) TemporaryBuffer(exec func(cb vulkan.CommandBuffer)) {
	// create tmp command buffer
	buffers := createBuffers(p.ld.Ref(), p.ref, 1)
	tmpBuffer := buffers[0]

	// execute some user command on it
	vulkan.BeginCommandBuffer(tmpBuffer, &vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: vulkan.CommandBufferUsageFlags(vulkan.CommandBufferUsageOneTimeSubmitBit),
	})

	exec(tmpBuffer)

	vulkan.EndCommandBuffer(tmpBuffer)

	// submit written commands to GPU execution
	submitInfo := vulkan.SubmitInfo{
		SType:              vulkan.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    []vulkan.CommandBuffer{buffers[0]},
	}
	vulkan.QueueSubmit(p.ld.QueueGraphics(), 1, []vulkan.SubmitInfo{submitInfo}, nil)

	// wait for GPU execute it
	vulkan.QueueWaitIdle(p.ld.QueuePresent())

	// now command buffer is not used anymore
	// and can be destroyed safely
	vulkan.FreeCommandBuffers(p.ld.Ref(), p.ref, uint32(len(buffers)), buffers)
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
	buffers := createBuffers(ld.Ref(), pool, buffersCount)

	// ok
	return pool, buffers
}

func createBuffers(ld vulkan.Device, pool vulkan.CommandPool, buffersCount uint32) []vulkan.CommandBuffer {
	allocInfo := &vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        pool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: buffersCount,
	}

	buffers := make([]vulkan.CommandBuffer, buffersCount)
	must.Work(vulkan.AllocateCommandBuffers(ld, allocInfo, buffers))

	return buffers
}
