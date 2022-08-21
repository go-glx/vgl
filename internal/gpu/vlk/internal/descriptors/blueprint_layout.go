package descriptors

import (
	"fmt"
	"strings"

	"github.com/vulkan-go/vulkan"
)

type (
	BlueprintLayout struct {
		index          layoutIndex
		title          string
		usage          string
		layout         vulkan.DescriptorSetLayout
		layoutBindings layoutBindings
	}
)

func (lay *BlueprintLayout) Ref() vulkan.DescriptorSetLayout {
	return lay.layout
}

func (lay *BlueprintLayout) String() string {
	return fmt.Sprintf("layout #%d (%s) with %d sets: %s",
		lay.index,
		lay.title,
		len(lay.layoutBindings),
		lay.layoutBindings.String(),
	)
}

func (lb layoutBindings) String() string {
	dst := make([]string, 0, len(lb))

	for _, bindingSet := range lb {
		dst = append(dst, fmt.Sprintf("[set#%d descriptors=%d type=%s]",
			bindingSet.Binding,
			bindingSet.DescriptorCount,
			nameOfDescriptorType(bindingSet.DescriptorType),
		))
	}

	return strings.Join(dst, ", ")
}

func nameOfDescriptorType(dt vulkan.DescriptorType) string {
	switch dt {
	case vulkan.DescriptorTypeSampler:
		return "Sampler"
	case vulkan.DescriptorTypeCombinedImageSampler:
		return "CombinedImageSampler"
	case vulkan.DescriptorTypeSampledImage:
		return "SampledImage"
	case vulkan.DescriptorTypeStorageImage:
		return "StorageImage"
	case vulkan.DescriptorTypeUniformTexelBuffer:
		return "UniformTexelBuffer"
	case vulkan.DescriptorTypeStorageTexelBuffer:
		return "StorageTexelBuffer"
	case vulkan.DescriptorTypeUniformBuffer:
		return "UniformBuffer"
	case vulkan.DescriptorTypeStorageBuffer:
		return "StorageBuffer"
	case vulkan.DescriptorTypeUniformBufferDynamic:
		return "UniformBufferDynamic"
	case vulkan.DescriptorTypeStorageBufferDynamic:
		return "StorageBufferDynamic"
	case vulkan.DescriptorTypeInputAttachment:
		return "InputAttachment"
	case vulkan.DescriptorTypeInlineUniformBlock:
		return "InlineUniformBlock"
	case vulkan.DescriptorTypeAccelerationStructureNvx:
		return "AccelerationStructureNvx"
	}

	return "Unknown"
}
