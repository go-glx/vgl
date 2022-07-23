package swapchain

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

type Chain struct {
	props     ChainProps
	swapChain vulkan.Swapchain
	images    []vulkan.Image
	views     []vulkan.ImageView
	buffers   []vulkan.Framebuffer

	ld *logical.Device
}

func NewChain(width, height uint32, pd *physical.Device, ld *logical.Device, surface *surface.Surface, mainRenderPass *renderpass.Pass, mobileFriendly bool) *Chain {
	props := newProps(width, height, pd, mobileFriendly)
	sharingMode := deviceSharingMode(pd)
	swapChain := newSwapChain(pd, ld, surface, props, sharingMode)

	images := createImages(swapChain, ld)
	views := createViews(images, ld, props)
	buffers := createFrameBuffers(ld, mainRenderPass.Ref(), props, views)

	log.Printf("vk: swapchain created, images=%d, props=(%s)\n", len(images), props.String())

	return &Chain{
		props:     props,
		swapChain: swapChain,
		images:    images,
		views:     views,
		buffers:   buffers,

		ld: ld,
	}
}

func (c *Chain) Free() {
	for _, buffer := range c.buffers {
		vulkan.DestroyFramebuffer(c.ld.Ref(), buffer, nil)
	}

	for _, view := range c.views {
		vulkan.DestroyImageView(c.ld.Ref(), view, nil)
	}

	vulkan.DestroySwapchain(c.ld.Ref(), c.swapChain, nil)
	log.Printf("vk: freed: swapchain\n")
}

func (c *Chain) Ref() vulkan.Swapchain {
	return c.swapChain
}

func (c *Chain) Props() ChainProps {
	return c.props
}

func (c *Chain) FrameBuffer(index int) vulkan.Framebuffer {
	return c.buffers[index]
}

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
