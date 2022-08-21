package pipeline

import (
	"github.com/vulkan-go/vulkan"
)

func withDefaultViewport() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, f *Factory) {
		info.PViewportState = &vulkan.PipelineViewportStateCreateInfo{
			SType:         vulkan.StructureTypePipelineViewportStateCreateInfo,
			ViewportCount: 1,
			PViewports:    []vulkan.Viewport{f.swapChain.Viewport()},
			ScissorCount:  1,
			PScissors:     []vulkan.Rect2D{f.swapChain.Scissor()},
		}
	}
}

func withDefaultMainRenderPass() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, f *Factory) {
		info.RenderPass = f.mainRenderPass.Ref()
		info.Subpass = 0
	}
}
