package vlk

import (
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/frame"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

func (c *Container) commandPool() *command.Pool {
	return dynamic(c, &c.vlkCommandPool,
		func(x *command.Pool) { x.Free() },
		func() *command.Pool {
			return command.NewPool(
				c.physicalDevice(),
				c.logicalDevice(),
			)
		},
	)
}

func (c *Container) frameManager() *frame.Manager {
	return dynamic(c, &c.vlkFrameManager,
		func(x *frame.Manager) { x.Free() },
		func() *frame.Manager {
			return frame.NewManager(
				c.logicalDevice(),
				c.commandPool(),
				c.swapChain(),
				c.renderPassMain(),
				c.rebuild,
			)
		},
	)
}

func (c *Container) swapChain() *swapchain.Chain {
	return dynamic(c, &c.vlkSwapChain,
		func(x *swapchain.Chain) { x.Free() },
		func() *swapchain.Chain {
			wWidth, wHeight := c.wm.GetFramebufferSize()
			return swapchain.NewChain(
				uint32(wWidth),
				uint32(wHeight),
				c.physicalDevice(),
				c.logicalDevice(),
				c.surface(),
				c.renderPassMain(),
				c.cfg.HasGPUVSync(),
			)
		},
	)
}

func (c *Container) renderPassMain() *renderpass.Pass {
	return dynamic(c, &c.vlkRenderPassMain,
		func(x *renderpass.Pass) { x.Free() },
		func() *renderpass.Pass {
			return renderpass.NewMain(
				c.physicalDevice(),
				c.logicalDevice(),
			)
		},
	)
}
