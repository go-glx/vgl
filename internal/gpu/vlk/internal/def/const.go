package def

import (
	"time"

	"github.com/vulkan-go/vulkan"
)

// ------------------------------------------------------
// -- Instance
// ------------------------------------------------------

// VKApiVersion is main vulkan API version
const VKApiVersion = vulkan.ApiVersion11

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
//
//	Min: 2 (double buffering)
//	Max: 3 (triple buffering)
//
// Recommended value 3. If GPU not support X Buffers, this will be
// automatic clamped as clamp(X, gpuMin, gpuMax)
const OptimalSwapChainBuffersCount = 3

// -- Format

// What format/color space we want for rendering
// If GPU not support this formats, render will panic
// on initialization.
const (
	SurfaceFormat     = vulkan.FormatB8g8r8a8Unorm
	SurfaceColorSpace = vulkan.ColorSpaceSrgbNonlinear
)

// ------------------------------------------------------
// -- Rendering
// ------------------------------------------------------

// FrameAcquireTimeout how much time CPU will wait
// for latest frame<n = OptimalSwapChainBuffersCount - 1>
// at frameStart. If GPU hang/lag and not present this N frame
// on screen, after FrameAcquireTimeout CPU will panic and crash
// application (before panic this will be retries X times)
const FrameAcquireTimeout = time.Second

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
//   - too small = more draw calls, less performance in intensive applications
//   - too big   = less draw calls, slower copy speed cpu->gpu = less performance in simple applications
//   - 16MB      = good in most cases
const BufferVertexSizeBytes = 16 * 1024 * 1024

// BufferIndexSizeBytes used for initial pre-generated shader indexes
// Count of indexed (batched) draw calls primary depend on this value
//
// Recommended value:
//   - too small = more draw calls, less performance in intensive applications
//   - too big   = just more GPU local immutable memory usage
//   - 4MB       = good in most cases
const BufferIndexSizeBytes = 4 * 1024 * 1024

// BufferUniformSizeBytes
// 16KB is minimum guaranteed on any device
// Recommended value:
// - <16KB
// - equal of real buffer usage * 2
const BufferUniformSizeBytes = 1024

// BufferStorageSizeBytes
// Common use storage buffer for shaders
// Recommended value:
// - big enough to store good chunk of random common use data
// - 32MB good in most cases
const BufferStorageSizeBytes = 32 * 1024 * 1024

// BufferIndexMaxInstances How many index data will be generated
// and saved to fast-persistent GPU buffer memory
// bufferSize = BufferIndexMaxInstances * instanceSize
// instanceSize = vertexSize * vertexesCount
//
// example:
//
//	for triangle, we need 3 vertexes
//	each vertex has 24 bytes (vec2(xy) + vec4(rgba))
//	total is 72 bytes per instance
//	total is 4_718_592 bytes (or 4.71 MB)
//
// this is maximum instance count, that can be drawn
// in one draw-call. So, when we want draw 100k instances
// library will use 2 draw calls.
const BufferIndexMaxInstances = 65536
