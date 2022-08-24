package vlk

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/descriptors"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

func (c *Container) instance() *instance.Instance {
	return static(c, &c.vlkInstance,
		func(x *instance.Instance) { x.Free() },
		func() *instance.Instance {
			// init proc addr
			procAddr := c.wm.InitVulkanProcAddr()

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
					procAddr,
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

func (c *Container) pipelineCache() *pipeline.Cache {
	return static(c, &c.vlkPipelineCache,
		func(x *pipeline.Cache) { x.Free() },
		func() *pipeline.Cache {
			return pipeline.NewCache(
				c.logger,
				c.logicalDevice(),
			)
		},
	)
}

func (c *Container) shaderManager() *shader.Manager {
	return static(c, &c.vlkShaderManager,
		func(x *shader.Manager) { x.Free() },
		func() *shader.Manager {
			return shader.NewManager(
				c.logger,
				c.logicalDevice(),
			)
		},
	)
}

func (c *Container) memoryAllocator() *alloc.Allocator {
	return static(c, &c.vlkMemoryAllocator,
		func(x *alloc.Allocator) { x.Free() },
		func() *alloc.Allocator {
			return alloc.NewAllocator(
				c.logger,
				c.instance(),
				c.physicalDevice(),
				c.logicalDevice(),
				c.commandPool(),
			)
		},
	)
}

func (c *Container) allocBuffers() *alloc.Buffers {
	return static(c, &c.vlkAllocBuffers,
		func(x *alloc.Buffers) {},
		func() *alloc.Buffers {
			return alloc.NewBuffers(
				c.allocHeap(),
			)
		},
	)
}

func (c *Container) allocHeap() *alloc.Heap {
	return static(c, &c.vlkAllocHeap,
		func(x *alloc.Heap) {},
		func() *alloc.Heap {
			return alloc.NewHeap(
				c.memoryAllocator(),
			)
		},
	)
}

func (c *Container) descriptorsManager() *descriptors.Manager {
	return static(c, &c.vlkDescriptorsManager,
		func(x *descriptors.Manager) {},
		func() *descriptors.Manager {
			return descriptors.NewManager(
				c.logger,
				c.logicalDevice(),
				c.physicalDevice(),
				c.descriptorsPool(),
				c.allocHeap(),
				c.descriptorsBlueprint(),
			)
		},
	)
}

func (c *Container) descriptorsPool() *descriptors.Pool {
	return static(c, &c.vlkDescriptorsPool,
		func(x *descriptors.Pool) { x.Free() },
		func() *descriptors.Pool {
			return descriptors.NewPool(
				c.logger,
				c.logicalDevice(),
				c.descriptorsBlueprint(),
			)
		},
	)
}

func (c *Container) descriptorsBlueprint() *descriptors.Blueprint {
	return static(c, &c.vlkDescriptorsBlueprint,
		func(x *descriptors.Blueprint) { x.Free() },
		func() *descriptors.Blueprint {
			return descriptors.NewBlueprint(
				c.logger,
				c.logicalDevice(),
			)
		},
	)
}
