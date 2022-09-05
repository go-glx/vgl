package dscptr

import (
	"fmt"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/shared/vlkext"
)

const framesCount = def.OptimalSwapChainBuffersCount

type Pool struct {
	logger vlkext.Logger
	ld     *logical.Device

	ref vulkan.DescriptorPool
}

func NewPool(logger vlkext.Logger, ld *logical.Device) *Pool {
	return &Pool{
		logger: logger,
		ld:     ld,

		ref: createPool(logger, ld),
	}
}

func (p *Pool) Free() {
	vulkan.DestroyDescriptorPool(p.ld.Ref(), p.ref, nil)
	p.logger.Debug("freed: descriptor pool")
}

func (p *Pool) Pool() vulkan.DescriptorPool {
	return p.ref
}

func createPool(logger vlkext.Logger, ld *logical.Device) vulkan.DescriptorPool {
	bufferTypes := map[vulkan.DescriptorType]uint32{}
	for _, layout := range blueprint {
		for _, binding := range layout.bindings {
			if _, exist := bufferTypes[binding.descriptorType]; !exist {
				bufferTypes[binding.descriptorType] = 0
			}

			bufferTypes[binding.descriptorType]++
		}
	}

	sizes := make([]vulkan.DescriptorPoolSize, 0, len(bufferTypes))
	sizesLogs := make([]string, 0, len(bufferTypes))
	for descriptorType, size := range bufferTypes {
		sizes = append(sizes, vulkan.DescriptorPoolSize{
			Type:            descriptorType,
			DescriptorCount: size,
		})
		sizesLogs = append(sizesLogs, fmt.Sprintf("%dx%s", size, nameOfDescriptorType(descriptorType)))
	}

	info := vulkan.DescriptorPoolCreateInfo{
		SType:         vulkan.StructureTypeDescriptorPoolCreateInfo,
		MaxSets:       uint32(framesCount * len(blueprint)),
		PoolSizeCount: uint32(len(sizes)),
		PPoolSizes:    sizes,
	}

	var pool vulkan.DescriptorPool
	must.Work(vulkan.CreateDescriptorPool(ld.Ref(), &info, nil, &pool))

	logger.Debug(fmt.Sprintf(
		"descriptor pool created with %d sets: %s",
		len(bufferTypes),
		strings.Join(sizesLogs, ", "),
	))

	return pool
}
