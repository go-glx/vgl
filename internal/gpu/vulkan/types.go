package vulkan

import (
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

type (
	Vk struct {
		inst               *vkInstance
		surface            *vkSurface
		pd                 *vkPhysicalDevice
		ld                 *vkLogicalDevice
		commandPool        *vkCommandPool
		frameManager       *vkFrameManager
		swapChain          *vkSwapChain
		frameBuffers       *vkFrameBuffers
		dataBuffersManager *vkDataBuffersManager
		shaderManager      *vkShaderManager
		pipelineManager    *vkPipelineManager
		pipelineLayout     vulkan.PipelineLayout

		// back link to container
		container *container

		// render variables
		renderQueue                    map[string][]shaderProgram
		currentFrameImageID            uint32
		currentFrameAvailableForRender bool
		inResizing                     bool
		isMinimized                    bool
	}

	// --------------------------------------
	// SYSTEM
	// --------------------------------------

	vkInstance struct {
		ref vulkan.Instance
	}

	vkSurface struct {
		inst *vkInstance
		ref  vulkan.Surface
	}

	// --------------------------------------
	// DEVICES
	// --------------------------------------

	vkPhysicalDeviceFinder struct {
		surface *vkSurface
		inst    *vkInstance
	}

	vkPhysicalDevice struct {
		ref          vulkan.PhysicalDevice
		props        vulkan.PhysicalDeviceProperties
		features     vulkan.PhysicalDeviceFeatures
		extensions   []vulkan.ExtensionProperties
		families     vkPhysicalDeviceFamilies
		surfaceProps vkPhysicalDeviceSurfaceProps
	}

	vkPhysicalDeviceFamilies struct {
		graphicsFamilyId uint32
		presentFamilyId  uint32
		supportGraphics  bool
		supportPresent   bool
	}

	vkPhysicalDeviceSurfaceProps struct {
		capabilities   vulkan.SurfaceCapabilities
		formats        []vulkan.SurfaceFormat
		presentModes   []vulkan.PresentMode
		surfaceSupport bool
	}

	vkLogicalDevice struct {
		ref           vulkan.Device
		queueGraphics vulkan.Queue
		queuePresent  vulkan.Queue
	}

	// --------------------------------------
	// Render Requirements
	// --------------------------------------

	vkCommandPool struct {
		ref     vulkan.CommandPool
		buffers []vulkan.CommandBuffer
		ld      *vkLogicalDevice
	}

	frameID = uint32
	imageID = uint32

	vkFrameManager struct {
		maxFrames         frameID
		currentFrameID    frameID
		currentImageID    imageID
		presentFailsCount int

		muxRenderAvailable  map[frameID]vulkan.Semaphore
		muxPresentAvailable map[frameID]vulkan.Semaphore
		fence               map[frameID]vulkan.Fence
		imagesInFlight      map[imageID]vulkan.Fence
		onSwapOutOfDate     func()

		ld *vkLogicalDevice
	}

	// --------------------------------------
	// Graphics
	// --------------------------------------

	vkSwapChain struct {
		ref          vulkan.Swapchain
		surfaceProps vkPhysicalDeviceSurfaceProps
		images       []vulkan.Image
		imagesView   []vulkan.ImageView
		info         vkSwapChainInfo

		ld *vkLogicalDevice
	}

	vkSwapChainInfo struct {
		imageFormat     vulkan.Format
		imageColorSpace vulkan.ColorSpace
		bufferSize      vulkan.Extent2D
		presentMode     vulkan.PresentMode
	}

	vkFrameBuffers struct {
		buffers []vulkan.Framebuffer

		clearColor vulkan.ClearValue
		ld         *vkLogicalDevice
		swapChain  *vkSwapChain
	}

	vkDataBuffersManager struct {
		vertex  vkBufferTable
		index   vkBufferUnion
		uniform []vkBufferUnion

		ld *vkLogicalDevice
		pd *vkPhysicalDevice
	}

	vkBufferUnion struct {
		offset  uint16
		buffer  vkBuffer
		staging []byte
	}

	vkBufferTable struct {
		totalCapacity       uint64
		buffers             []vkBuffer
		framePageID         int16
		framePageCapacity   uint64
		frameStagedData     [][]byte
		frameInstanceCounts []uint32
	}

	vkBuffer struct {
		capacity vulkan.DeviceSize
		dataPtr  unsafe.Pointer
		handle   vulkan.Buffer
		memory   vulkan.DeviceMemory

		// todo: writeAt (time) - clear unused memory buffers, if not used > 1m
	}

	vkPipelineManager struct {
		shaderPipelines map[shaderModuleID]vulkan.Pipeline

		ld *vkLogicalDevice
	}

	// --------------------------------------
	// Shaders
	// --------------------------------------

	shaderModuleID  = string
	vkShaderManager struct {
		modules map[shaderModuleID]*vkShaderModule

		ld *vkLogicalDevice
	}

	vkShaderModule struct {
		module    vulkan.ShaderModule
		stageInfo vulkan.PipelineShaderStageCreateInfo
	}

	shaderProgram interface {
		// shaderPipelines bytecode

		ID() string
		ProgramFrag() []byte
		ProgramVert() []byte

		// attributes

		Size() uint64
		TriangleCount() uint32
		VertexCount() uint32
		Topology() vulkan.PrimitiveTopology
		Indexes() []uint16
		Data() []byte
		Bindings() []vulkan.VertexInputBindingDescription
		Attributes() []vulkan.VertexInputAttributeDescription
	}

	shaderPipelineFactory = func(c *container, sp shaderProgram) vulkan.Pipeline
)
