package swapchain

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func createFrameBuffers(ld *logical.Device, mainRenderPass vulkan.RenderPass, props ChainProps, views []vulkan.ImageView) []vulkan.Framebuffer {
	buffers := make([]vulkan.Framebuffer, 0, len(views))

	for _, view := range views {
		buffers = append(buffers, createFrameBuffer(ld, mainRenderPass, props, view))
	}

	return buffers
}

func createFrameBuffer(ld *logical.Device, mainRenderPass vulkan.RenderPass, props ChainProps, view vulkan.ImageView) vulkan.Framebuffer {
	info := &vulkan.FramebufferCreateInfo{
		SType:           vulkan.StructureTypeFramebufferCreateInfo,
		RenderPass:      mainRenderPass,
		AttachmentCount: 1,
		PAttachments: []vulkan.ImageView{
			view,
		},
		Width:  props.BufferSize.Width,
		Height: props.BufferSize.Height,
		Layers: 1,
	}

	var buffer vulkan.Framebuffer
	must.Work(vulkan.CreateFramebuffer(ld.Ref(), info, nil, &buffer))

	return buffer
}
