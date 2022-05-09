package vulkan

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/arch"
	"github.com/fe3dback/govgl/config"
)

func NewVulkanApi(wm arch.WindowManager, cfg *config.Config) *Vk {
	cont := newContainer(wm, cfg)
	renderer := cont.renderer()

	// subscribe to system events
	wm.OnWindowResized(func(_, _ int) {
		renderer.rebuildGraphicsPipeline()
	})

	return renderer
}

func (vk *Vk) Close() error {
	vk.free()

	return nil
}

func (vk *Vk) free() {
	if vk.pipelineLayout != nil {
		vulkan.DestroyPipelineLayout(vk.ld.ref, vk.pipelineLayout, nil)
		vk.pipelineLayout = nil
		log.Printf("vk: freed: pipeline layout\n")
	}

	if vk.container.vkPipelineLayoutUBODescriptorSet != nil {
		vulkan.DestroyDescriptorSetLayout(vk.ld.ref, vk.container.vkPipelineLayoutUBODescriptorSet, nil)
		vk.container.vkPipelineLayoutUBODescriptorSet = nil
		log.Printf("vk: freed: descriptor set layout\n")
	}

	if vk.pipelineManager != nil {
		vk.pipelineManager.free()
		vk.pipelineManager = nil
	}

	if vk.shaderManager != nil {
		vk.shaderManager.free()
		vk.shaderManager = nil
	}

	if vk.dataBuffersManager != nil {
		vk.dataBuffersManager.free()
		vk.dataBuffersManager = nil
	}

	if vk.frameBuffers != nil {
		vk.frameBuffers.free()
		vk.frameBuffers = nil
	}

	for _, renderPass := range vk.container.vkRenderPassHandlesLazyCache {
		vulkan.DestroyRenderPass(vk.ld.ref, renderPass, nil)
	}
	log.Printf("vk: freed: all render pases\n")

	if vk.swapChain != nil {
		vk.swapChain.free()
		vk.swapChain = nil
	}

	if vk.frameManager != nil {
		vk.frameManager.free()
		vk.frameManager = nil
	}

	if vk.commandPool != nil {
		vk.commandPool.free()
		vk.commandPool = nil
	}

	if vk.ld != nil {
		vk.ld.free()
		vk.ld = nil
	}

	vk.pd = nil
	log.Printf("vk: freed: physical device\n")

	if vk.surface != nil {
		vk.surface.free()
		vk.surface = nil
	}

	if vk.inst != nil {
		vk.inst.free()
		vk.inst = nil
	}

	log.Printf("vk: freed: renderer complete freed\n")
}

func (vk *Vk) rebuildGraphicsPipeline() {
	// resize mode
	// ----------------------------

	if vk.inResizing {
		// already in rebuildGraphicsPipeline mode
		// wait for end
		return
	}

	vk.inResizing = true

	// wait for all current render is done
	// ----------------------------
	vulkan.DeviceWaitIdle(vk.ld.ref)

	// free all pipeline staff
	// ----------------------------

	if vk.pipelineManager != nil {
		vk.pipelineManager.free()
		vk.pipelineManager = nil
		vk.container.vkPipelineManager = nil
	}

	for _, renderPass := range vk.container.vkRenderPassHandlesLazyCache {
		vulkan.DestroyRenderPass(vk.ld.ref, renderPass, nil)
	}

	vk.container.vkRenderPassHandlesLazyCache = make(map[renderPassType]vulkan.RenderPass)

	if vk.frameBuffers != nil {
		vk.frameBuffers.free()
		vk.frameBuffers = nil
		vk.container.vkFrameBuffers = nil
	}

	if vk.swapChain != nil {
		vk.swapChain.free()
		vk.swapChain = nil
		vk.container.vkSwapChain = nil
	}

	if vk.commandPool != nil {
		vk.commandPool.free()
		vk.commandPool = nil
		vk.container.vkCommandPool = nil
	}

	// minimization handle
	// ----------------------------

	vk.isMinimized = false
	wWidth, wHeight := vk.container.windowManager.GetFramebufferSize()
	if wWidth == 0 || wHeight == 0 {
		// window is minimized now, just wait for next resize for
		// swapChain recreate

		vk.inResizing = false
		vk.isMinimized = true
		return
	}

	// recreate vk objects
	// ----------------------------

	vk.commandPool = vk.container.provideVkCommandPool()
	vk.swapChain = vk.container.provideSwapChain()
	vk.frameBuffers = vk.container.provideFrameBuffers()
	vk.pipelineManager = vk.container.providePipelineManager()

	for shader, pipelineFactory := range buildInShaders {
		vk.pipelineManager.preloadPipelineFor(shader, pipelineFactory(vk.container, shader))
	}

	// finalize
	// ----------------------------

	vk.inResizing = false
}
