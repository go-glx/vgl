package dscptr

import (
	"github.com/vulkan-go/vulkan"
)

const (
	LayoutIndexGlobal layoutIndex = 0
	LayoutIndexObject layoutIndex = 1
	LayoutIndexLocal  layoutIndex = 2
)

type (
	frameID      = uint32
	layoutIndex  = uint32
	bindingIndex = uint32

	blueprintLayout struct {
		title       string
		description string
		bindings    blueprintBindingsMap
	}

	blueprintBinding struct {
		descriptorType vulkan.DescriptorType
		flags          vulkan.ShaderStageFlagBits
	}

	blueprintLayoutMap   = map[layoutIndex]blueprintLayout
	blueprintBindingsMap = map[bindingIndex]blueprintBinding
)

var blueprint = blueprintLayoutMap{
	LayoutIndexGlobal: {
		title:       "Global",
		description: "Has two binding for vert={[view, projection] matrix} frag={surface.size.xy}, used in every frame as global UBO",
		bindings: blueprintBindingsMap{
			0: {
				descriptorType: vulkan.DescriptorTypeUniformBuffer,
				flags:          vulkan.ShaderStageVertexBit,
			},
			1: {
				descriptorType: vulkan.DescriptorTypeUniformBuffer,
				flags:          vulkan.ShaderStageFragmentBit,
			},
		},
	},
	// LayoutIndexObject: {
	// 	title:       "Object",
	// 	description: "general purpose custom object buffers",
	// 	bindings: blueprintBindingsMap{
	// 		0: {
	// 			descriptorType: vulkan.DescriptorTypeStorageBuffer,
	// 			flags:          vulkan.ShaderStageVertexBit | vulkan.ShaderStageFragmentBit,
	// 		},
	// 		1: {
	// 			descriptorType: vulkan.DescriptorTypeStorageBuffer,
	// 			flags:          vulkan.ShaderStageVertexBit | vulkan.ShaderStageFragmentBit,
	// 		},
	// 		2: {
	// 			descriptorType: vulkan.DescriptorTypeStorageBuffer,
	// 			flags:          vulkan.ShaderStageVertexBit | vulkan.ShaderStageFragmentBit,
	// 		},
	// 		3: {
	// 			descriptorType: vulkan.DescriptorTypeStorageBuffer,
	// 			flags:          vulkan.ShaderStageVertexBit | vulkan.ShaderStageFragmentBit,
	// 		},
	// 	},
	// },
	// LayoutIndexLocal: {
	// 	title:       "Local",
	// 	description: "Small uniform with object transform matrix, used only for 3D objects",
	// 	bindings: blueprintBindingsMap{
	// 		0: {
	// 			descriptorType: vulkan.DescriptorTypeUniformBuffer,
	// 			flags:          vulkan.ShaderStageVertexBit,
	// 		},
	// 	},
	// },
}
