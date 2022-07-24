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

func (f *Factory) withMainRenderPass() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.RenderPass = f.mainRenderPass.Ref()
		info.Subpass = 0
	}
}

func (f *Factory) newDefaultPipelineLayout() vulkan.PipelineLayout {
	info := &vulkan.PipelineLayoutCreateInfo{
		SType: vulkan.StructureTypePipelineLayoutCreateInfo,
		// SetLayoutCount:         1,
		// PSetLayouts:            []vulkan.DescriptorSetLayout{ubo},
		SetLayoutCount:         0,   // todo ^
		PSetLayouts:            nil, // todo ^
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}

	var pipelineLayout vulkan.PipelineLayout
	must.Work(vulkan.CreatePipelineLayout(f.ld.Ref(), info, nil, &pipelineLayout))

	return pipelineLayout
}
