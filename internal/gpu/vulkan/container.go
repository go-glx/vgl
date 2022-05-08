package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/arch"
	"github.com/fe3dback/govgl/config"
	"github.com/fe3dback/govgl/internal/closer"
)

type (
	container struct {
		// dependencies
		windowManager arch.WindowManager
		cfg           *config.Config
		closer        *closer.Closer

		// internal
		vkRenderPassHandlesLazyCache map[renderPassType]vulkan.RenderPass
		vkPipelineHandlesLazyCache   map[shaderModuleID]vulkan.Pipeline

		// vk handle wrappers
		vk                               *Vk
		vkInstance                       *vkInstance
		vkSurface                        *vkSurface
		vkPhysicalDevice                 *vkPhysicalDevice
		vkLogicalDevice                  *vkLogicalDevice
		vkCommandPool                    *vkCommandPool
		vkFrameManager                   *vkFrameManager
		vkSwapChain                      *vkSwapChain
		vkFrameBuffers                   *vkFrameBuffers
		vkVertexBuffers                  *vkDataBuffersManager
		vkShaderManager                  *vkShaderManager
		vkPipelineManager                *vkPipelineManager
		vkPipelineLayout                 vulkan.PipelineLayout
		vkPipelineLayoutUBODescriptorSet vulkan.DescriptorSetLayout
	}
)

func newContainer(wm arch.WindowManager, cfg *config.Config) *container {
	contCloser := closer.NewCloser()
	contCloser.EnqueueClose(wm.Close)

	return &container{
		windowManager: wm,
		cfg:           cfg,
		closer:        contCloser,

		// internal
		vkRenderPassHandlesLazyCache: map[renderPassType]vulkan.RenderPass{},
		vkPipelineHandlesLazyCache:   map[shaderModuleID]vulkan.Pipeline{},
	}
}

func (c *container) renderer() *Vk {
	if c.vk != nil {
		return c.vk
	}

	// init vk proc addr with wm
	c.windowManager.InitVulkanProcAddr()

	// init gpu
	err := vulkan.Init()
	if err != nil {
		panic(fmt.Errorf("failed init vulkan: %w", err))
	}

	log.Printf("vk: lib initialized: [%#v]\n", c.cfg)

	// main
	c.vk = &Vk{
		renderQueue: make(map[string][]shaderProgram),
	}
	c.vk.container = c
	c.vk.inst = c.provideVkInstance()
	c.vk.surface = c.provideVkSurface()
	c.vk.pd = c.provideVkPhysicalDevice()
	c.vk.ld = c.provideVkLogicalDevice()

	// render utils
	c.vk.commandPool = c.provideVkCommandPool()
	c.vk.frameManager = c.provideFrameManager(c.vk.rebuildGraphicsPipeline)
	c.vk.swapChain = c.provideSwapChain()
	c.vk.frameBuffers = c.provideFrameBuffers()
	c.vk.dataBuffersManager = c.provideDataBuffersManager()
	c.vk.shaderManager = c.provideShaderManager()
	c.vk.pipelineManager = c.providePipelineManager()
	c.vk.pipelineLayout = c.providePipelineLayout()

	for shader, pipelineFactory := range buildInShaders {
		c.vk.shaderManager.preloadShader(shader)
		c.vkPipelineManager.preloadPipelineFor(shader, pipelineFactory(c, shader))
	}

	// render

	return c.vk
}

func (c *container) provideVkInstance() *vkInstance {
	if c.vkInstance != nil {
		return c.vkInstance
	}

	// required ext
	requiredExt := c.windowManager.GetRequiredInstanceExtensions()

	// todo: debug callbacks

	// init
	c.vkInstance = newVkInstance(c.windowManager, requiredExt, c.cfg.InDebug())
	return c.vkInstance
}

func (c *container) provideVkSurface() *vkSurface {
	if c.vkSurface != nil {
		return c.vkSurface
	}

	c.vkSurface = newSurfaceFromWindow(
		c.provideVkInstance(),
		c.windowManager,
	)
	return c.vkSurface
}

func (c *container) provideVkPhysicalDevice() *vkPhysicalDevice {
	if c.vkPhysicalDevice != nil {
		return c.vkPhysicalDevice
	}

	finder := newPhysicalDeviceFinder(
		c.provideVkInstance(),
		c.provideVkSurface(),
	)

	c.vkPhysicalDevice = finder.physicalDevicePick()
	return c.vkPhysicalDevice
}

func (c *container) provideVkLogicalDevice() *vkLogicalDevice {
	if c.vkLogicalDevice != nil {
		return c.vkLogicalDevice
	}

	c.vkLogicalDevice = newLogicalDevice(
		c.provideVkPhysicalDevice(),
	)
	return c.vkLogicalDevice
}

func (c *container) provideVkCommandPool() *vkCommandPool {
	if c.vkCommandPool != nil {
		return c.vkCommandPool
	}

	c.vkCommandPool = newCommandPool(
		c.provideVkPhysicalDevice(),
		c.provideVkLogicalDevice(),
	)
	return c.vkCommandPool
}

func (c *container) provideFrameManager(onSwapOutOfDate func()) *vkFrameManager {
	if c.vkFrameManager != nil {
		return c.vkFrameManager
	}

	c.vkFrameManager = newFrameManager(
		c.provideVkLogicalDevice(),
		c.provideVkPhysicalDevice(),
		onSwapOutOfDate,
	)
	return c.vkFrameManager
}

func (c *container) provideSwapChain() *vkSwapChain {
	if c.vkSwapChain != nil {
		return c.vkSwapChain
	}

	wWidth, wHeight := c.windowManager.GetFramebufferSize()

	c.vkSwapChain = newSwapChain(
		uint32(wWidth), uint32(wHeight),
		c.provideVkPhysicalDevice(),
		c.provideVkLogicalDevice(),
		c.provideVkSurface(),
		c.cfg,
	)
	return c.vkSwapChain
}

func (c *container) provideFrameBuffers() *vkFrameBuffers {
	if c.vkFrameBuffers != nil {
		return c.vkFrameBuffers
	}

	c.vkFrameBuffers = newFrameBuffers(
		c.provideVkLogicalDevice(),
		c.provideSwapChain(),
		c.defaultRenderPass(),
	)
	return c.vkFrameBuffers
}

func (c *container) provideDataBuffersManager() *vkDataBuffersManager {
	if c.vkVertexBuffers != nil {
		return c.vkVertexBuffers
	}

	c.vkVertexBuffers = newDataBuffersManager(
		c.provideVkLogicalDevice(),
		c.provideVkPhysicalDevice(),
	)
	return c.vkVertexBuffers
}

func (c *container) provideShaderManager() *vkShaderManager {
	if c.vkShaderManager != nil {
		return c.vkShaderManager
	}

	c.vkShaderManager = newShaderManager(
		c.provideVkLogicalDevice(),
	)
	return c.vkShaderManager
}

func (c *container) providePipelineManager() *vkPipelineManager {
	if c.vkPipelineManager != nil {
		return c.vkPipelineManager
	}

	c.vkPipelineManager = newPipelineManager(
		c.provideVkLogicalDevice(),
	)
	return c.vkPipelineManager
}

func (c *container) providePipelineLayout() vulkan.PipelineLayout {
	if c.vkPipelineLayout != nil {
		return c.vkPipelineLayout
	}

	c.vkPipelineLayout = newPipeLineLayout(
		c.provideVkLogicalDevice(),
		c.providePipelineLayoutUBODescriptorSet(),
	)
	return c.vkPipelineLayout
}

func (c *container) providePipelineLayoutUBODescriptorSet() vulkan.DescriptorSetLayout {
	if c.vkPipelineLayoutUBODescriptorSet != nil {
		return c.vkPipelineLayoutUBODescriptorSet
	}

	c.vkPipelineLayoutUBODescriptorSet = newPipeLineLayoutUBODescriptorSet(
		c.provideVkLogicalDevice(),
	)
	return c.vkPipelineLayoutUBODescriptorSet
}
