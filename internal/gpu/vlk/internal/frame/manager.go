package frame

import (
	"fmt"
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
	"github.com/go-glx/vgl/shared/vlkext"
)

// todo: need some refactoring with available state management

// how many continues errors is ok, before crash
const maxWaitImageErrorsCount = 16

type (
	frameID uint32
	imageID uint32

	Context struct {
		isAvailable bool
		frameID     frameID
		imageID     imageID
	}
)

func (c Context) FrameID() uint32 {
	return uint32(c.frameID)
}

type Manager struct {
	logger         vlkext.Logger
	chain          *swapchain.Chain
	mainRenderPass *renderpass.Pass
	ld             *logical.Device
	onSuboptimal   func()

	count uint32

	semRenderAvailable  map[frameID]vulkan.Semaphore
	semPresentAvailable map[frameID]vulkan.Semaphore
	syncFrameBusy       map[frameID]vulkan.Fence
	commandBuffers      map[frameID]vulkan.CommandBuffer
}

func NewManager(
	logger vlkext.Logger,
	ld *logical.Device,
	pool *command.Pool,
	chain *swapchain.Chain,
	renderToScreenPass *renderpass.Pass,
	onSuboptimal func(),
) *Manager {
	m := &Manager{
		logger:         logger,
		chain:          chain,
		mainRenderPass: renderToScreenPass,
		ld:             ld,
		onSuboptimal:   onSuboptimal,

		count: uint32(pool.MainBuffersCount()),

		semRenderAvailable:  make(map[frameID]vulkan.Semaphore),
		semPresentAvailable: make(map[frameID]vulkan.Semaphore),
		syncFrameBusy:       make(map[frameID]vulkan.Fence),
		commandBuffers:      make(map[frameID]vulkan.CommandBuffer),
	}

	for fID := uint32(0); fID < m.count; fID++ {
		m.commandBuffers[frameID(fID)] = pool.MainCommandBuffer(int(fID))
		m.semRenderAvailable[frameID(fID)] = allocateSemaphore(ld)
		m.semPresentAvailable[frameID(fID)] = allocateSemaphore(ld)
		m.syncFrameBusy[frameID(fID)] = allocateFence(ld)
	}

	logger.Debug("frame manager created")
	return m
}

func (m *Manager) Free() {
	for fID := uint32(0); fID < m.count; fID++ {
		vulkan.DestroyFence(m.ld.Ref(), m.syncFrameBusy[frameID(fID)], nil)
		vulkan.DestroySemaphore(m.ld.Ref(), m.semPresentAvailable[frameID(fID)], nil)
		vulkan.DestroySemaphore(m.ld.Ref(), m.semRenderAvailable[frameID(fID)], nil)
	}

	m.logger.Debug("freed: frames manager")
}

func (m *Manager) FrameBegin(prev Context) (Context, bool) {
	// get next frame id
	frmID := m.nextFrameID(prev.frameID)

	// wait until this frame is ready for rendering to (frame = logical)
	m.waitUntilFrameReady(frmID)

	// acquire new imageID for rendering into (image = physical)
	imgID, available := m.acquireNextImage(frmID)
	if !available {
		return Context{
			isAvailable: false,
			frameID:     0,
			imageID:     0,
		}, false
	}

	// create new context
	ctx := Context{
		isAvailable: true,
		frameID:     frmID,
		imageID:     imgID,
	}

	// run start commands
	m.commandBufferBegin(ctx)

	// start render pass
	m.FrameApplyCommands(ctx, func(cb vulkan.CommandBuffer) {
		m.renderPassMainBegin(uint32(imgID), cb)
	})

	// start user space command listening
	return ctx, true
}

func (m *Manager) FrameApplyCommands(ctx Context, apply func(cb vulkan.CommandBuffer)) {
	if !ctx.isAvailable {
		return
	}

	apply(m.commandBuffers[ctx.frameID])
}

func (m *Manager) FrameEnd(ctx Context) {
	if !ctx.isAvailable {
		return
	}

	// end render pass
	m.FrameApplyCommands(ctx, func(cb vulkan.CommandBuffer) {
		m.renderPassMainEnd(cb)
	})

	// end buffer
	m.commandBufferEnd(ctx)

	// submit rendering on GPU
	m.submit(ctx.frameID, ctx.imageID)
}

func (m *Manager) nextFrameID(current frameID) frameID {
	return (current + 1) % frameID(m.count)
}

func (m *Manager) waitUntilFrameReady(frameID frameID) {
	m.waitUntilFenceReady(m.syncFrameBusy[frameID])
}

func (m *Manager) waitUntilFenceReady(fence vulkan.Fence) {
	timeout := uint64(def.FrameAcquireTimeout.Nanoseconds())
	triesCount := maxWaitImageErrorsCount

	for triesCount > 0 {
		// wait for rendering/logic in current frame is done
		ok := m.notice(vulkan.WaitForFences(m.ld.Ref(), 1, []vulkan.Fence{fence}, vulkan.True, timeout))
		if ok {
			vulkan.ResetFences(m.ld.Ref(), 1, []vulkan.Fence{fence})
			return
		}

		triesCount--
		m.logger.Notice(fmt.Sprintf("busy: wait for device rendering done (%d tries left)..", triesCount))
	}

	panic(fmt.Errorf("render image not available after %d tries", maxWaitImageErrorsCount))
	return
}

func (m *Manager) acquireNextImage(frameID frameID) (imageID, bool) {
	timeout := uint64(def.FrameAcquireTimeout.Nanoseconds())

	id := uint32(0)
	result := vulkan.AcquireNextImage(m.ld.Ref(), m.chain.Ref(), timeout, m.semRenderAvailable[frameID], nil, &id)

	if result == vulkan.ErrorOutOfDate || result == vulkan.Suboptimal {
		// buffer size changes (window rebuildGraphicsPipeline, minimize, etc..)
		// and not more valid
		m.logger.Notice("Suboptimal or outOfDate image (maybe resolution changed). Will enforce pipeline recreation..")
		m.onSuboptimal()
		return 0, false
	}

	if result != vulkan.Success {
		m.notice(result)
		return 0, false
	}

	return imageID(id), true
}

func (m *Manager) submit(frameID frameID, imageID imageID) {
	if !m.render(frameID) {
		return
	}

	if !m.present(frameID, imageID) {
		return
	}

	vulkan.QueueWaitIdle(m.ld.QueuePresent())
}

func (m *Manager) render(frameID frameID) bool {
	info := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{m.semRenderAvailable[frameID]},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{m.commandBuffers[frameID]},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{m.semPresentAvailable[frameID]},
	}

	return m.notice(vulkan.QueueSubmit(m.ld.QueueGraphics(), 1, []vulkan.SubmitInfo{info}, m.syncFrameBusy[frameID]))
}

func (m *Manager) present(frameID frameID, imageID imageID) bool {
	info := &vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{m.semPresentAvailable[frameID]},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{m.chain.Ref()},
		PImageIndices:      []uint32{uint32(imageID)},
	}

	return m.notice(vulkan.QueuePresent(m.ld.QueuePresent(), info))
}
