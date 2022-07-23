package renderpass

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

type Pass struct {
	ref vulkan.RenderPass

	ld *logical.Device
}

func newPass(ld *logical.Device, renderPass vulkan.RenderPass) *Pass {
	return &Pass{
		ref: renderPass,

		ld: ld,
	}
}

func (p *Pass) Free() {
	vulkan.DestroyRenderPass(p.ld.Ref(), p.ref, nil)
}

func (p *Pass) Ref() vulkan.RenderPass {
	return p.ref
}

func createPass(
	name string,
	ld *logical.Device,
	attachments []vulkan.AttachmentDescription,
	subPasses []vulkan.SubpassDescription,
	dependencies []vulkan.SubpassDependency,
) vulkan.RenderPass {
	info := &vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(attachments)),
		PAttachments:    attachments,
		SubpassCount:    uint32(len(subPasses)),
		PSubpasses:      subPasses,
		DependencyCount: uint32(len(dependencies)),
		PDependencies:   dependencies,
	}

	var renderPass vulkan.RenderPass
	must.Work(vulkan.CreateRenderPass(ld.Ref(), info, nil, &renderPass))

	log.Printf("vk: render pass '%s' created\n", name)

	return renderPass
}
