package pipeline

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func (f *Factory) withDefaultViewport() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.PViewportState = &vulkan.PipelineViewportStateCreateInfo{
			SType:         vulkan.StructureTypePipelineViewportStateCreateInfo,
			ViewportCount: 1,
			PViewports:    []vulkan.Viewport{f.swapChain.Viewport()},
			ScissorCount:  1,
			PScissors:     []vulkan.Rect2D{f.swapChain.Scissor()},
		}
	}
}

func (f *Factory) withDefaultMainRenderPass() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.RenderPass = f.mainRenderPass.Ref()
		info.Subpass = 0
	}
}

func (f *Factory) withDefaultLayout() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.Layout = f.defaultPipelineLayout
	}
}

func (f *Factory) newDefaultPipelineLayout() vulkan.PipelineLayout {
	info := &vulkan.PipelineLayoutCreateInfo{
		SType:                  vulkan.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount:         1,
		PSetLayouts:            []vulkan.DescriptorSetLayout{f.defaultDescriptionSetLayout},
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}

	var pipelineLayout vulkan.PipelineLayout
	must.Work(vulkan.CreatePipelineLayout(f.ld.Ref(), info, nil, &pipelineLayout))

	return pipelineLayout
}

func (f *Factory) newDefaultDescriptorSetLayout() vulkan.DescriptorSetLayout {
	// 0 = global UBO matrix (model * view * projection)
	bindingGlobal := vulkan.DescriptorSetLayoutBinding{
		Binding:         0,
		DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit), // ubo used only in vert shaders
	}

	info := &vulkan.DescriptorSetLayoutCreateInfo{
		SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: 1,
		PBindings: []vulkan.DescriptorSetLayoutBinding{
			bindingGlobal,
		},
	}

	var layout vulkan.DescriptorSetLayout
	must.Work(vulkan.CreateDescriptorSetLayout(f.ld.Ref(), info, nil, &layout))

	return layout
}
