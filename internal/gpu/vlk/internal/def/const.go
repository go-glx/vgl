package def

import (
	"time"

	"github.com/vulkan-go/vulkan"
)

// ------------------------------------------------------
// -- Instance
// ------------------------------------------------------

const (
	vkVersion11 = 0x00401000
	// const vkVersion12 = 0x00402000
	// const vkVersion13 = 0x00403000
)

// VKApiVersion is main vulkan API version
const VKApiVersion = vkVersion11
const VKVMAVersion = vkVersion11

// RequiredValidationLayers will be used, only when renderer in debug mode
var RequiredValidationLayers = []string{
	"VK_LAYER_KHRONOS_validation",
}

// ------------------------------------------------------
// -- Device
// ------------------------------------------------------

// RequiredDeviceExtensions list of required ext that GPU should support
// if system not have GPU with listed extensions, VK will return error on
// initialization
var RequiredDeviceExtensions = []string{
	"VK_KHR_swapchain", // require for display buffer to screen
}

// ------------------------------------------------------
// -- SwapChain
// ------------------------------------------------------

// OptimalSwapChainBuffersCount defines how many images/buffers we want
// for optimal rendering
//  Min: 2 (double buffering)
//  Max: 3 (triple buffering)
// Recommended value 3. If GPU not support X Buffers, this will be
// automatic clamped as clamp(X, gpuMin, gpuMax)
const OptimalSwapChainBuffersCount = 3

// -- Format

// What format/color space we want for rendering
// If GPU not support this formats, render will panic
// on initialization.
const (
	SurfaceFormat     = vulkan.FormatB8g8r8a8Srgb
	SurfaceColorSpace = vulkan.ColorSpaceSrgbNonlinear
)

// ------------------------------------------------------
// -- Rendering
// ------------------------------------------------------

// FrameAcquireTimeout how much time CPU will wait
// for latest frame<n = OptimalSwapChainBuffersCount - 1>
// at frameStart. If GPU hang/lag and not present this N frame
// on screen, after FrameAcquireTimeout CPU will panic and crash
// application
const FrameAcquireTimeout = time.Second * 3

// ShaderEntryPoint is entry point in shader bytecode
// where GPU start executing shader code
// do not change from "main"
const ShaderEntryPoint = "main"

// BufferVertexSizeBytes used for transport vertex data from cpu to gpu
// vertex data mostly is [positions, colors]
//
// maximum buffer size for each drawCall
// if API want draw more objects, that fit into one buffer
// it will be split into few draw Calls and buffer flush
//
// Recommended value:
//  - too small = more draw calls, less performance in intensive applications
//  - too big   = less draw calls, slower copy speed cpu->gpu = less performance in simple applications
//  - 65536     = good in most cases
const BufferVertexSizeBytes = 2048 // todo:65536

// BufferIndexSizeBytes used for transport index data from cpu to gpu
// index is simple uint16 array. Each index is pointer to specific offset
// in vertex buffer inside BufferVertexSizeBytes
const BufferIndexSizeBytes = 2048
