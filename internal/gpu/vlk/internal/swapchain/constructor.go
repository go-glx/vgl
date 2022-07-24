package swapchain

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

func newSwapChain(pd *physical.Device, ld *logical.Device, surface *surface.Surface, props ChainProps, sharingMode vulkan.SharingMode) vulkan.Swapchain {
	families := pd.PrimaryGPU().Families.UniqueIDs()

	info := &vulkan.SwapchainCreateInfo{
		SType:                 vulkan.StructureTypeSwapchainCreateInfo,
		Surface:               surface.Ref(),
		MinImageCount:         props.BuffersCount,
		ImageFormat:           props.ImageFormat,
		ImageColorSpace:       props.ImageColorSpace,
		ImageExtent:           props.BufferSize,
		ImageArrayLayers:      1,
		ImageUsage:            vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		ImageSharingMode:      sharingMode,
		QueueFamilyIndexCount: uint32(len(families)),
		PQueueFamilyIndices:   families,
		PreTransform:          pd.PrimaryGPU().SurfaceProps.Capabilities().CurrentTransform,
		CompositeAlpha:        vulkan.CompositeAlphaOpaqueBit,
		PresentMode:           props.PresentMode,
		Clipped:               vulkan.True,
	}

	var swapChain vulkan.Swapchain
	must.Work(vulkan.CreateSwapchain(ld.Ref(), info, nil, &swapChain))

	return swapChain
}
