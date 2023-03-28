package vlk

import (
	"github.com/go-glx/vgl/shared/config"
	"github.com/go-glx/vgl/shared/vlkext"
)

type closer interface {
	EnqueueFree(fn func())
	EnqueueBackFree(fn func())
}

type Container struct {
	closer    closer
	logger    vlkext.Logger
	rebuilder *rebuilder
	wm        vlkext.WindowManager
	cfg       *config.Config
}

func NewContainer(
	closer closer,
	wm vlkext.WindowManager,
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
