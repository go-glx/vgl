package swapchain

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
	"github.com/go-glx/vgl/shared/vlkext"
)

type Chain struct {
	logger    vlkext.Logger
	props     ChainProps
	swapChain vulkan.Swapchain
	images    []vulkan.Image
	views     []vulkan.ImageView
	buffers   []vulkan.Framebuffer

	ld *logical.Device
}

func NewChain(logger vlkext.Logger, width, height uint32, pd *physical.Device, ld *logical.Device, surface *surface.Surface, mainRenderPass *renderpass.Pass, mobileFriendly bool) *Chain {
	props := newProps(width, height, pd, mobileFriendly)
	sharingMode := deviceSharingMode(pd)
	swapChain := newSwapChain(pd, ld, surface, props, sharingMode)

	images := createImages(swapChain, ld)
	views := createViews(images, ld, props)
	buffers := createFrameBuffers(ld, mainRenderPass.Ref(), props, views)

	logger.Debug(fmt.Sprintf("swapchain created, images=%d, props=(%s)", len(images), props.String()))

	return &Chain{
		logger:    logger,
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
	c.logger.Debug("freed: swapchain")
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

func (c *Chain) Viewport() vulkan.Viewport {
	return vulkan.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(c.Props().BufferSize.Width),
		Height:   float32(c.Props().BufferSize.Height),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}
}

func (c *Chain) Scissor() vulkan.Rect2D {
	return vulkan.Rect2D{
		Offset: vulkan.Offset2D{
			X: 0,
			Y: 0,
		},
		Extent: c.Props().BufferSize,
	}
}
