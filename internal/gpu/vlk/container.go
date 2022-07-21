package vlk

import (
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

type closer interface {
	EnqueueFree(fn func())
	EnqueueClose(fn func() error)
}

type Container struct {
	closer    closer
	rebuilder *rebuilder
	wm        arch.WindowManager
	cfg       *config.Config

	// static
	vlkRef            *VLK
	vlkInstance       *instance.Instance
	vlkSurface        *surface.Surface
	vlkPhysicalDevice *physical.Device
	vlkLogicalDevice  *logical.Device

	// dynamic
	vlkCommandPool *command.Pool
}

func NewContainer(
	closer closer,
	wm arch.WindowManager,
	cfg *config.Config,
) *Container {
	cont := &Container{
		closer:    closer,
		rebuilder: newRebuilder(),
		wm:        wm,
		cfg:       cfg,
	}

	wm.OnWindowResized(func(_, _ int) {
		cont.rebuild()
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
