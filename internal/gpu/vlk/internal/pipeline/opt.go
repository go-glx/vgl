package pipeline

import "github.com/vulkan-go/vulkan"

type Initializer = func(*vulkan.GraphicsPipelineCreateInfo)

func WithStages(stages []vulkan.PipelineShaderStageCreateInfo) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.StageCount = uint32(len(stages))
		info.PStages = stages
	}
}

func WithTopology(topology vulkan.PrimitiveTopology) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.PInputAssemblyState = &vulkan.PipelineInputAssemblyStateCreateInfo{
			SType:                  vulkan.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology:               topology,
			PrimitiveRestartEnable: vulkan.False,
		}
	}
}

func WithVertexInput(
	bindings []vulkan.VertexInputBindingDescription,
	attributes []vulkan.VertexInputAttributeDescription,
) Initializer {
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
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
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
		info.PRasterizationState = &vulkan.PipelineRasterizationStateCreateInfo{
			SType:                   vulkan.StructureTypePipelineRasterizationStateCreateInfo,
			DepthClampEnable:        vulkan.False,
			RasterizerDiscardEnable: vulkan.False,
			PolygonMode:             mode,

			// todo: vulkan.CullModeBackBit, after faces is debugged and correct
			CullMode:  vulkan.CullModeFlags(vulkan.CullModeNone),
			FrontFace: vulkan.FrontFaceClockwise,

			DepthBiasEnable:         vulkan.False,
			DepthBiasConstantFactor: 0.0,
			DepthBiasClamp:          0.0,
			DepthBiasSlopeFactor:    0.0,
			LineWidth:               1.0, // todo: require ext
		}
	}
}

func WithMultisampling() Initializer {
	// todo: options
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
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
	// todo: options
	return func(info *vulkan.GraphicsPipelineCreateInfo) {
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
