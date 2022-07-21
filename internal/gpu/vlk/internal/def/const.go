package def

import "github.com/vulkan-go/vulkan"

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
