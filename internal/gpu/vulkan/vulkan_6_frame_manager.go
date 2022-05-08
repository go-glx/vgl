package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"
)

func newFrameManager(ld *vkLogicalDevice, pd *vkPhysicalDevice, onSwapOutOfDate func()) *vkFrameManager {
	maxFrames := pd.surfaceProps.imageBuffersCount()
	frameManager := &vkFrameManager{
		maxFrames:           maxFrames,
		currentFrameID:      0,
		currentImageID:      0,
		muxRenderAvailable:  make(map[frameID]vulkan.Semaphore),
		muxPresentAvailable: make(map[frameID]vulkan.Semaphore),
		fence:               make(map[frameID]vulkan.Fence),
		imagesInFlight:      make(map[imageID]vulkan.Fence),
		ld:                  ld,
		onSwapOutOfDate:     onSwapOutOfDate,
	}

	for i := frameID(0); i < maxFrames; i++ {
		frameManager.muxRenderAvailable[i] = allocateSemaphore(ld)
		frameManager.muxPresentAvailable[i] = allocateSemaphore(ld)
		frameManager.fence[i] = allocateFence(ld)
	}

	return frameManager
}

func (fm *vkFrameManager) free() {
	for i := frameID(0); i < fm.maxFrames; i++ {
		vulkan.DestroyFence(fm.ld.ref, fm.fence[i], nil)
		vulkan.DestroySemaphore(fm.ld.ref, fm.muxPresentAvailable[i], nil)
		vulkan.DestroySemaphore(fm.ld.ref, fm.muxRenderAvailable[i], nil)
	}

	log.Printf("vk: freed: frame manager\n")
}

func allocateSemaphore(ld *vkLogicalDevice) vulkan.Semaphore {
	var ref vulkan.Semaphore

	muxCreateInfo := &vulkan.SemaphoreCreateInfo{
		SType: vulkan.StructureTypeSemaphoreCreateInfo,
	}
	vkAssert(
		vulkan.CreateSemaphore(ld.ref, muxCreateInfo, nil, &ref),
		fmt.Errorf("failed create image available mux"),
	)

	return ref
}

func allocateFence(ld *vkLogicalDevice) vulkan.Fence {
	var fence vulkan.Fence

	fenceCreateInfo := &vulkan.FenceCreateInfo{
		SType: vulkan.StructureTypeFenceCreateInfo,
		Flags: vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit),
	}
	vkAssert(
		vulkan.CreateFence(ld.ref, fenceCreateInfo, nil, &fence),
		fmt.Errorf("failed create present available mux"),
	)

	return fence
}

func (fm *vkFrameManager) frameStart(swapChain *vkSwapChain) (imageIndex uint32, successful bool) {
	// todo: check this code, fences logic

	timeout := uint64(swapChainTimeout.Nanoseconds())
	muxRender := fm.muxRenderAvailable[fm.currentFrameID]
	fence := fm.fence[fm.currentFrameID]

	// wait current ops
	vulkan.WaitForFences(fm.ld.ref, 1, []vulkan.Fence{fence}, vulkan.True, timeout)

	// acquire new image
	fm.currentImageID = uint32(0)
	result := vulkan.AcquireNextImage(fm.ld.ref, swapChain.ref, timeout, muxRender, nil, &fm.currentImageID)
	if result == vulkan.ErrorOutOfDate || result == vulkan.Suboptimal {
		// buffer size changes (window rebuildGraphicsPipeline, minimize, etc..)
		// and not more valid
		fm.onSwapOutOfDate()
		return 0, false
	}

	vkAssert(result, fmt.Errorf("failed acquire next image for rendering"))

	// wait when image will be available
	if imageFence, inFlight := fm.imagesInFlight[fm.currentImageID]; inFlight {
		vulkan.WaitForFences(fm.ld.ref, 1, []vulkan.Fence{imageFence}, vulkan.True, timeout)
	}
	fm.imagesInFlight[fm.currentImageID] = fence

	vulkan.ResetFences(fm.ld.ref, 1, []vulkan.Fence{fence})

	return fm.currentImageID, true
}

func (fm *vkFrameManager) frameEnd(swapChain *vkSwapChain, frameCommandBuffer vulkan.CommandBuffer) {
	muxRender := fm.muxRenderAvailable[fm.currentFrameID]
	muxPresent := fm.muxPresentAvailable[fm.currentFrameID]
	fence := fm.fence[fm.currentFrameID]

	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{muxRender},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{frameCommandBuffer},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{muxPresent},
	}

	vkAssert(
		vulkan.QueueSubmit(fm.ld.queueGraphics, 1, []vulkan.SubmitInfo{submitInfo}, fence),
		fmt.Errorf("failed submit graphics queue"),
	)

	presentInfo := &vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{muxPresent},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{swapChain.ref},
		PImageIndices:      []uint32{fm.currentImageID},
	}

	result := vulkan.QueuePresent(fm.ld.queuePresent, presentInfo)
	if result != vulkan.Success {
		fm.presentFailsCount++
	} else {
		fm.presentFailsCount = 0
	}

	if fm.presentFailsCount > maxPresetFails {
		panic(fmt.Errorf("failed present rendered graphics into screen, more that %d times", maxPresetFails))
	}

	vulkan.QueueWaitIdle(fm.ld.queuePresent)

	// done
	fm.currentFrameID = (fm.currentFrameID + 1) % fm.maxFrames
}
