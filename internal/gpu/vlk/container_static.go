package vlk

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

func (c *Container) instance() *instance.Instance {
	return static(c, &c.vlkInstance,
		func(x *instance.Instance) { x.Free() },
		func() *instance.Instance {
			// init proc addr
			c.wm.InitVulkanProcAddr()

			// init vulkan driver
			err := vulkan.Init()
			if err != nil {
				panic(fmt.Errorf("failed init vulkan: %w", err))
			}

			c.logger.Info(fmt.Sprintf("lib initialized: [%#v]", c.cfg))

			// create instance
			return instance.NewInstance(
				instance.NewCreateOptions(
					c.logger,
					c.wm.AppName(),
					c.wm.EngineName(),
					c.wm.GetRequiredInstanceExtensions(),
					c.cfg.InDebug(),
				),
			)
		},
	)
}

func (c *Container) surface() *surface.Surface {
	return static(c, &c.vlkSurface,
		func(x *surface.Surface) { x.Free() },
		func() *surface.Surface {
			return surface.NewSurface(
				c.logger,
				c.instance(),
				c.wm,
			)
		},
	)
}

func (c *Container) physicalDevice() *physical.Device {
	return static(c, &c.vlkPhysicalDevice,
		func(x *physical.Device) {},
		func() *physical.Device {
			return physical.NewDevice(
				c.logger,
				c.instance(),
				c.surface(),
			)
		},
	)
}

func (c *Container) logicalDevice() *logical.Device {
	return static(c, &c.vlkLogicalDevice,
		func(x *logical.Device) { x.Free() },
		func() *logical.Device {
			return logical.NewDevice(
				c.logger,
				c.physicalDevice(),
			)
		},
	)
}

func (c *Container) shaderManager() *shader.Manager {
	return static(c, &c.vlkShaderManager,
		func(x *shader.Manager) { x.Free() },
		func() *shader.Manager {
			mng := shader.NewManager(
				c.logger,
				c.logicalDevice(),
			)

			// register build-in shaders
			mng.RegisterShader(defaultShaderTriangle())

			//
			return mng
		},
	)
}
