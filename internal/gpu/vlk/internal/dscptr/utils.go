package dscptr

import "github.com/vulkan-go/vulkan"

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
