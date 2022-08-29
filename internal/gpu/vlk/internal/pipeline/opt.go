package pipeline

import (
	"github.com/vulkan-go/vulkan"
)

type Initializer = func(*vulkan.GraphicsPipelineCreateInfo, *Factory)

func WithStages(stages []vulkan.PipelineShaderStageCreateInfo) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.StageCount = uint32(len(stages))
		info.PStages = stages
	}
}

func WithTopology(topology vulkan.PrimitiveTopology, restartEnable bool) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.PInputAssemblyState = &vulkan.PipelineInputAssemblyStateCreateInfo{
			SType:                  vulkan.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology:               topology,
			PrimitiveRestartEnable: castToVKBool(restartEnable),
		}
	}
}

func WithVertexInput(
	bindings []vulkan.VertexInputBindingDescription,
	attributes []vulkan.VertexInputAttributeDescription,
) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.PVertexInputState = &vulkan.PipelineVertexInputStateCreateInfo{
			SType:                           vulkan.StructureTypePipelineVertexInputStateCreateInfo,
			VertexBindingDescriptionCount:   uint32(len(bindings)),
			PVertexBindingDescriptions:      bindings,
			VertexAttributeDescriptionCount: uint32(len(attributes)),
			PVertexAttributeDescriptions:    attributes,
		}
	}
}

func WithRasterization(mode vulkan.PolygonMode) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.PRasterizationState = &vulkan.PipelineRasterizationStateCreateInfo{
			SType:                   vulkan.StructureTypePipelineRasterizationStateCreateInfo,
			DepthClampEnable:        vulkan.False,
			RasterizerDiscardEnable: vulkan.False, // todo: on and check on all apis
			PolygonMode:             mode,
			CullMode:                vulkan.CullModeFlags(vulkan.CullModeBackBit),
			FrontFace:               vulkan.FrontFaceClockwise,
			DepthBiasEnable:         vulkan.False,
			DepthBiasConstantFactor: 0.0,
			DepthBiasClamp:          0.0,
			DepthBiasSlopeFactor:    0.0,
			LineWidth:               1.0,
		}
	}
}

func WithMultisampling() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.PMultisampleState = &vulkan.PipelineMultisampleStateCreateInfo{
			SType:                 vulkan.StructureTypePipelineMultisampleStateCreateInfo,
			RasterizationSamples:  vulkan.SampleCount1Bit,
			SampleShadingEnable:   vulkan.False,
			MinSampleShading:      1.0,
			PSampleMask:           nil,
			AlphaToCoverageEnable: vulkan.False,
			AlphaToOneEnable:      vulkan.False,
		}
	}
}

func WithColorBlend() Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo, _ *Factory) {
		info.PColorBlendState = &vulkan.PipelineColorBlendStateCreateInfo{
			SType:           vulkan.StructureTypePipelineColorBlendStateCreateInfo,
			LogicOpEnable:   vulkan.False,
			LogicOp:         vulkan.LogicOpCopy,
			AttachmentCount: 1,
			PAttachments: []vulkan.PipelineColorBlendAttachmentState{{
				BlendEnable:         vulkan.True,
				SrcColorBlendFactor: vulkan.BlendFactorSrcAlpha,
				DstColorBlendFactor: vulkan.BlendFactorOneMinusSrcAlpha,
				ColorBlendOp:        vulkan.BlendOpAdd,
				SrcAlphaBlendFactor: vulkan.BlendFactorOne,
				DstAlphaBlendFactor: vulkan.BlendFactorZero,
				AlphaBlendOp:        vulkan.BlendOpAdd,
				ColorWriteMask: vulkan.ColorComponentFlags(
					vulkan.ColorComponentRBit | vulkan.ColorComponentGBit | vulkan.ColorComponentBBit | vulkan.ColorComponentABit,
				),
			}},
			BlendConstants: [4]float32{0, 0, 0, 0},
		}
	}
}

func castToVKBool(b bool) vulkan.Bool32 {
	if b {
		return vulkan.True
	}

	return vulkan.False
}
