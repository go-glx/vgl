package descriptors

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func layoutGlobal(ld vulkan.Device) BlueprintLayout {
	globalLayout, globalBindings := createLayoutGlobal(ld)

	return BlueprintLayout{
		index:          layoutIndexGlobal,
		title:          "Global UBO Layout",
		usage:          "Has one binding for [view, projection] matrix, used in every frame as global UBO",
		layout:         globalLayout,
		layoutBindings: globalBindings,
	}
}

func createLayoutGlobal(ld vulkan.Device) (vulkan.DescriptorSetLayout, layoutBindings) {
	bindings := layoutBindings{
		bindingGlobalIndexUniforms: {
			Binding:         uint32(bindingGlobalIndexUniforms),
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit), // ubo used only in vert shaders
		},
	}

	info := &vulkan.DescriptorSetLayoutCreateInfo{
		SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings: []vulkan.DescriptorSetLayoutBinding{
			bindings[bindingGlobalIndexUniforms],
		},
	}

	var layout vulkan.DescriptorSetLayout
	must.Work(vulkan.CreateDescriptorSetLayout(ld, info, nil, &layout))

	return layout, bindings
}
