package frame

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

// frame is single [start -> submit -> render -> present] operation
// usually we have at least two frames, that can process in ||
//
//  Frame is [
//     startFrame
//     lock (cpuSync)
//       start commandBuffer
//       api calls [ renderPass, drawCalls, any commands to buffer]
//       end commandBuffer
//       submit render (semaphore to present)
//       submit present (will unlock main fence when done)
//     unlock (cpuSync)
//  ]
type frame struct {
	chain           *swapchain.Chain
	commandBuffer   vulkan.CommandBuffer
	fenceBufferFree vulkan.Fence
	semImageReady   vulkan.Semaphore
	semRenderDone   vulkan.Semaphore
	onSuboptimal    func()
	freed           bool
	alive           bool

	ld *logical.Device
}

func newFrame(ld *logical.Device, chain *swapchain.Chain, commandBuffer vulkan.CommandBuffer, onSuboptimal func()) *frame {
	return &frame{
		chain:           chain,
		commandBuffer:   commandBuffer,
		fenceBufferFree: allocateFence(ld),
		semImageReady:   allocateSemaphore(ld),
		semRenderDone:   allocateSemaphore(ld),
		onSuboptimal:    onSuboptimal,
		freed:           false,
		alive:           true,

		ld: ld,
	}
}

func (f *frame) free() {
	if f.freed {
		return
	}

	f.freed = true
	vulkan.DestroyFence(f.ld.Ref(), f.fenceBufferFree, nil)
	vulkan.DestroySemaphore(f.ld.Ref(), f.semImageReady, nil)
	vulkan.DestroySemaphore(f.ld.Ref(), f.semRenderDone, nil)
}

func (f *frame) frameBegin() {
	if f.freed {
		f.alive = false
		return
	}

	timeout := uint64(def.FrameAcquireTimeout.Nanoseconds())
	isFree := must.NotCare(vulkan.WaitForFences(f.ld.Ref(), 1, []vulkan.Fence{f.fenceBufferFree}, vulkan.True, timeout))
	must.NotCare(vulkan.ResetFences(f.ld.Ref(), 1, []vulkan.Fence{f.fenceBufferFree}))

	if !isFree {
		f.alive = false
		return
	}

	f.commandBufferBegin()
	f.alive = true
}

func (f *frame) frameWrite(exec func(cb vulkan.CommandBuffer)) {
	if !f.alive {
		return
	}

	exec(f.commandBuffer)
}

func (f *frame) frameEnd() {
	if !f.alive {
		return
	}

	f.commandBufferEnd()

	imageID, imageExist := f.acquireNextImage()
	if !imageExist {
		return
	}

	if submitted := f.submit(); !submitted {
		return
	}

	if queued := f.present(imageID); !queued {
		return
	}
}

func (f *frame) acquireNextImage() (uint32, bool) {
	timeout := uint64(def.FrameAcquireTimeout.Nanoseconds())
	imageID := uint32(256)

	result := vulkan.AcquireNextImage(f.ld.Ref(), f.chain.Ref(), timeout, f.semImageReady, nil, &imageID)
	if result == vulkan.ErrorOutOfDate || result == vulkan.Suboptimal {
		// buffer size changes (window rebuildGraphicsPipeline, minimize, etc..)
		// and not more valid
		f.onSuboptimal()
		return imageID, false
	}

	if result != vulkan.Success {
		must.NotCare(result)
		return imageID, false
	}

	return imageID, true
}

func (f *frame) submit() bool {
	info := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{f.semImageReady},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{f.commandBuffer},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{f.semRenderDone},
	}

	return must.NotCare(vulkan.QueueSubmit(f.ld.QueueGraphics(), 1, []vulkan.SubmitInfo{info}, f.fenceBufferFree))
}

func (f *frame) present(targetImageID uint32) bool {
	info := &vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{f.semRenderDone},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{f.chain.Ref()},
		PImageIndices:      []uint32{targetImageID},
	}

	return must.NotCare(vulkan.QueuePresent(f.ld.QueuePresent(), info))
}

func (f *frame) commandBufferBegin() {
	must.Work(
		vulkan.BeginCommandBuffer(f.commandBuffer, &vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
		}),
	)
}

func (f *frame) commandBufferEnd() {
	must.Work(
		vulkan.EndCommandBuffer(f.commandBuffer),
	)
}
