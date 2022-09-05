package pipeline

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/shared/vlkext"
)

type Cache struct {
	ref vulkan.PipelineCache

	logger vlkext.Logger
	ld     *logical.Device
}

func NewCache(logger vlkext.Logger, ld *logical.Device) *Cache {
	info := &vulkan.PipelineCacheCreateInfo{
		SType: vulkan.StructureTypePipelineCacheCreateInfo,
	}

	var cache vulkan.PipelineCache
	vulkan.CreatePipelineCache(ld.Ref(), info, nil, &cache)
	logger.Debug("pipeline cache created")

	return &Cache{
		ref:    cache,
		logger: logger,
		ld:     ld,
	}
}

func (c *Cache) Free() {
	vulkan.DestroyPipelineCache(c.ld.Ref(), c.ref, nil)
	c.logger.Debug("freed: pipeline cache")
}

func (c *Cache) Ref() vulkan.PipelineCache {
	return c.ref
}
