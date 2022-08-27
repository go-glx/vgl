package vlk

import (
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/frame"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

func (c *Container) commandPool() *command.Pool {
	return dynamic(c, func() *command.Pool {
		return command.NewPool(
			c.logger,
			c.physicalDevice(),
			c.logicalDevice(),
		)
	})
}

func (c *Container) frameManager() *frame.Manager {
	return dynamic(c, func() *frame.Manager {
		return frame.NewManager(
			c.logger,
			c.logicalDevice(),
			c.commandPool(),
			c.swapChain(),
			c.renderPassMain(),
			c.rebuild,
			&c.vlkFrameRenderingAvailable,
		)
	})
}

func (c *Container) swapChain() *swapchain.Chain {
	return dynamic(c, func() *swapchain.Chain {
		wWidth, wHeight := c.wm.GetFramebufferSize()
		return swapchain.NewChain(
			c.logger,
			uint32(wWidth),
			uint32(wHeight),
			c.physicalDevice(),
			c.logicalDevice(),
			c.surface(),
			c.renderPassMain(),
			c.cfg.IsMobileFriendly(),
		)
	})
}

func (c *Container) renderPassMain() *renderpass.Pass {
	return dynamic(c, func() *renderpass.Pass {
		return renderpass.NewMain(
			c.logger,
			c.physicalDevice(),
			c.logicalDevice(),
		)
	})
}

func (c *Container) pipelineFactory() *pipeline.Factory {
	return dynamic(c, func() *pipeline.Factory {
		return pipeline.NewFactory(
			c.logger,
			c.logicalDevice(),
			c.swapChain(),
			c.renderPassMain(),
			c.descriptorsManager(),
			c.pipelineCache(),
		)
	})
}
