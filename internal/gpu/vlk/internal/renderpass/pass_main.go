package renderpass

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/shared/vlkext"
)

// NewMain return main render pass that used for rendering
// buffers to window screen surface
func NewMain(logger vlkext.Logger, pd *physical.Device, ld *logical.Device) *Pass {
	return newPass(
		ld,
		createPass(
			"main",
			logger,
			ld,
			mainAttachments(pd),
			mainSubPasses(),
			mainDependencies(),
		),
	)
}

func mainAttachments(pd *physical.Device) []vulkan.AttachmentDescription {
	return []vulkan.AttachmentDescription{
		{
			Format:         pd.PrimaryGPU().SurfaceProps.RichColorSpaceFormat().Format,
			Samples:        vulkan.SampleCount1Bit,
			LoadOp:         vulkan.AttachmentLoadOpClear,
			StoreOp:        vulkan.AttachmentStoreOpStore,
			StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
			StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
			InitialLayout:  vulkan.ImageLayoutUndefined,
			FinalLayout:    vulkan.ImageLayoutPresentSrc,
		},
	}
}

func mainSubPasses() []vulkan.SubpassDescription {
	return []vulkan.SubpassDescription{
		{
			PipelineBindPoint:    vulkan.PipelineBindPointGraphics,
			InputAttachmentCount: 0,
			PInputAttachments:    nil,
			ColorAttachmentCount: 1,
			PColorAttachments: []vulkan.AttachmentReference{{
				Attachment: 0,
				Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
			}},
			PResolveAttachments:     nil,
			PDepthStencilAttachment: nil,
			PreserveAttachmentCount: 0,
			PPreserveAttachments:    nil,
		},
	}
}

func mainDependencies() []vulkan.SubpassDependency {
	return []vulkan.SubpassDependency{
		{
			SrcSubpass:    vulkan.SubpassExternal,
			DstSubpass:    0,
			SrcStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
			DstStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
			SrcAccessMask: 0,
			DstAccessMask: vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit),
		},
	}
}
