package vlk

import (
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/frame"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

type closer interface {
	EnqueueFree(fn func())
	EnqueueBackFree(fn func())
}

type Container struct {
	closer    closer
	logger    config.Logger
	rebuilder *rebuilder
	wm        arch.WindowManager
	cfg       *config.Config

	// global state

	// this only for use inside of Frame Manager
	// it cannot be part of fm struct, because FM
	// will be recreated every time, when GPU suboptimal (window resize)
	vlkFrameRenderingAvailable bool

	// static

	vlkRef             *VLK
	vlkInstance        *instance.Instance
	vlkSurface         *surface.Surface
	vlkPhysicalDevice  *physical.Device
	vlkLogicalDevice   *logical.Device
	vlkPipelineCache   *pipeline.Cache
	vlkShaderManager   *shader.Manager
	vlkMemoryAllocator *alloc.Allocator
	vlkAllocBuffers    *alloc.Buffers

	// dynamic

	vlkCommandPool     *command.Pool
	vlkSwapChain       *swapchain.Chain
	vlkFrameManager    *frame.Manager
	vlkRenderPassMain  *renderpass.Pass
	vlkPipelineFactory *pipeline.Factory
}

func NewContainer(
	closer closer,
	wm arch.WindowManager,
	cfg *config.Config,
) *Container {
	cont := &Container{
		closer:    closer,
		logger:    cfg.Logger(),
		rebuilder: newRebuilder(),
		wm:        wm,
		cfg:       cfg,
	}

	closer.EnqueueBackFree(func() {
		cont.rebuilder.free()
		cont.logger.Debug("freed: dynamic resources")
	})

	return cont
}

func (c *Container) rebuild() {
	c.VulkanRenderer().maintenance(func() {
		// free all dynamic resources
		c.rebuilder.free()

		// after maintenance is end
		// all of these resources will be automatic
		// lazy recreated when needed by graphic pipeline

		wWidth, wHeight := c.wm.GetFramebufferSize()
		c.vlkRef.surfacesSize[0] = [2]uint32{uint32(wWidth), uint32(wHeight)}
	})
}

func (c *Container) VulkanRenderer() *VLK {
	return static(c, &c.vlkRef,
		func(x *VLK) {},
		func() *VLK {
			return newVLK(c)
		},
	)
}
