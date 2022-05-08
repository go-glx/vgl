package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

type (
	renderPassType string

	renderPassFactory = func(crp *container) renderPassFactoryParams

	renderPassFactoryParams struct {
		attachments  []vulkan.AttachmentDescription
		subPasses    []vulkan.SubpassDescription
		dependencies []vulkan.SubpassDependency
	}
)

const (
	renderPassTypeDefault renderPassType = "default"
)

var passMap = map[renderPassType]renderPassFactory{
	renderPassTypeDefault: func(c *container) renderPassFactoryParams {
		return renderPassFactoryParams{
			attachments:  []vulkan.AttachmentDescription{c.defaultColorAttachment()},
			subPasses:    []vulkan.SubpassDescription{c.defaultSubPass()},
			dependencies: []vulkan.SubpassDependency{c.defaultDependency()},
		}
	},
}

func (c *container) lazyRenderPass(id renderPassType) vulkan.RenderPass {
	if rp, exist := c.vkRenderPassHandlesLazyCache[id]; exist {
		return rp
	}

	paramsFactory, ok := passMap[id]
	if !ok {
		panic(fmt.Errorf("failed find render pass factory '%s'", id))
	}

	params := paramsFactory(c)
	renderPassInfo := &vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(params.attachments)),
		PAttachments:    params.attachments,
		SubpassCount:    uint32(len(params.subPasses)),
		PSubpasses:      params.subPasses,
		DependencyCount: uint32(len(params.dependencies)),
		PDependencies:   params.dependencies,
	}

	var renderPass vulkan.RenderPass
	vkAssert(
		vulkan.CreateRenderPass(c.provideVkLogicalDevice().ref, renderPassInfo, nil, &renderPass),
		fmt.Errorf("failed create render pass '%s'", id),
	)

	c.vkRenderPassHandlesLazyCache[id] = renderPass
	return c.vkRenderPassHandlesLazyCache[id]
}

func (c *container) defaultRenderPass() vulkan.RenderPass {
	return c.lazyRenderPass(renderPassTypeDefault)
}

func (c *container) defaultSubPass() vulkan.SubpassDescription {
	colorAttachmentRef := vulkan.AttachmentReference{
		Attachment: 0,
		Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
	}

	return vulkan.SubpassDescription{
		PipelineBindPoint:       vulkan.PipelineBindPointGraphics,
		InputAttachmentCount:    0,
		PInputAttachments:       nil,
		ColorAttachmentCount:    1,
		PColorAttachments:       []vulkan.AttachmentReference{colorAttachmentRef},
		PResolveAttachments:     nil,
		PDepthStencilAttachment: nil,
		PreserveAttachmentCount: 0,
		PPreserveAttachments:    nil,
	}
}

func (c *container) defaultDependency() vulkan.SubpassDependency {
	return vulkan.SubpassDependency{
		SrcSubpass:    vulkan.SubpassExternal,
		DstSubpass:    0,
		SrcStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
		DstStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
		SrcAccessMask: 0,
		DstAccessMask: vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit),
	}
}

func (c *container) defaultColorAttachment() vulkan.AttachmentDescription {
	return vulkan.AttachmentDescription{
		Format:         c.provideVkPhysicalDevice().surfaceProps.richColorSpaceFormat().Format,
		Samples:        vulkan.SampleCount1Bit,
		LoadOp:         vulkan.AttachmentLoadOpClear,
		StoreOp:        vulkan.AttachmentStoreOpStore,
		StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
		InitialLayout:  vulkan.ImageLayoutUndefined,
		FinalLayout:    vulkan.ImageLayoutPresentSrc,
	}
}
