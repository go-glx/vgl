package vlk

import (
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
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
		c.VulkanRenderer().surfacesSize[0] = [2]float32{float32(wWidth), float32(wHeight)}
	})
}

func (c *Container) VulkanRenderer() *VLK {
	return static(c, func() *VLK {
		return newVLK(c)
	})
}
