package descriptors

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func layoutGlobal(ld vulkan.Device) BlueprintLayout {
	layout, bindings := createLayoutGlobal(ld)

	return BlueprintLayout{
		index:    layoutIndexGlobal,
		title:    "Global UBO Layout",
		usage:    "Has two binding for vert={[view, projection] matrix} flag={surface.size.xy}, used in every frame as global UBO",
		layout:   layout,
		bindings: bindings,
	}
}

func createLayoutGlobal(ld vulkan.Device) (vulkan.DescriptorSetLayout, layoutBindings) {
	bindings := layoutBindings{
		bindingGlobalUniforms: {
			Binding:         uint32(bindingGlobalUniforms),
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit), // ubo used only in vert shaders
		},
		bindingGlobalSurfaceSize: {
			Binding:         uint32(bindingGlobalSurfaceSize),
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageFragmentBit), // ubo used only in frag shaders
		},
	}

	info := &vulkan.DescriptorSetLayoutCreateInfo{
		SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings: []vulkan.DescriptorSetLayoutBinding{
			bindings[bindingGlobalUniforms],
			bindings[bindingGlobalSurfaceSize],
		},
	}

	var layout vulkan.DescriptorSetLayout
	must.Work(vulkan.CreateDescriptorSetLayout(ld, info, nil, &layout))

	return layout, bindings
}
