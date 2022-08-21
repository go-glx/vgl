package descriptors

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

type Pool struct {
	logger config.Logger
	ld     *logical.Device

	ref vulkan.DescriptorPool
}

func NewPool(logger config.Logger, ld *logical.Device, blueprint *Blueprint) *Pool {
	pool := createPool(ld, blueprint)
	logger.Debug("descriptor pool created")

	return &Pool{
		logger: logger,
		ld:     ld,

		ref: pool,
	}
}

func (p *Pool) Free() {
	vulkan.DestroyDescriptorPool(p.ld.Ref(), p.ref, nil)
	p.logger.Debug("freed: descriptor pool")
}

func (p *Pool) Pool() vulkan.DescriptorPool {
	return p.ref
}

func createPool(ld *logical.Device, blueprint *Blueprint) vulkan.DescriptorPool {
	sizes := make([]vulkan.DescriptorPoolSize, 0)

	for _, layout := range blueprint.layouts {
		for _, binding := range layout.bindings {
			sizes = append(sizes, vulkan.DescriptorPoolSize{
				Type:            binding.DescriptorType,
				DescriptorCount: binding.DescriptorCount,
			})
		}
	}

	info := vulkan.DescriptorPoolCreateInfo{
		SType:         vulkan.StructureTypeDescriptorPoolCreateInfo,
		MaxSets:       uint32(framesCount * len(blueprint.layouts)),
		PoolSizeCount: uint32(len(sizes)),
		PPoolSizes:    sizes,
	}

	var pool vulkan.DescriptorPool
	must.Work(vulkan.CreateDescriptorPool(ld.Ref(), &info, nil, &pool))

	return pool
}
