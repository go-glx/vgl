package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"
)

func newFrameBuffers(ld *vkLogicalDevice, swapChain *vkSwapChain, renderPass vulkan.RenderPass) *vkFrameBuffers {
	fb := &vkFrameBuffers{
		ld:        ld,
		swapChain: swapChain,
	}

	buffers := make([]vulkan.Framebuffer, 0, len(swapChain.imagesView))
	for _, view := range swapChain.imagesView {
		createInfo := &vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      renderPass,
			AttachmentCount: 1,
			PAttachments: []vulkan.ImageView{
				view,
			},
			Width:  swapChain.info.bufferSize.Width,
			Height: swapChain.info.bufferSize.Height,
			Layers: 1,
		}

		var buffer vulkan.Framebuffer
		vkAssert(
			vulkan.CreateFramebuffer(ld.ref, createInfo, nil, &buffer),
			fmt.Errorf("failed create frame buffer"),
		)

		buffers = append(buffers, buffer)
	}

	fb.buffers = buffers
	return fb
}

func (fb *vkFrameBuffers) free() {
	for _, buffer := range fb.buffers {
		vulkan.DestroyFramebuffer(fb.ld.ref, buffer, nil)
	}

	log.Printf("vk: freed: frame buffers\n")
}

func (fb *vkFrameBuffers) renderPassStart(ind int, commandBuffer vulkan.CommandBuffer, renderPass vulkan.RenderPass) {
	renderPassBeginInfo := &vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  renderPass,
		Framebuffer: fb.buffers[ind],
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{
				X: 0,
				Y: 0,
			},
			Extent: vulkan.Extent2D{
				Width:  fb.swapChain.info.bufferSize.Width,
				Height: fb.swapChain.info.bufferSize.Height,
			},
		},
		ClearValueCount: 1,
		PClearValues:    fb.getClearColor(),
	}

	vulkan.CmdBeginRenderPass(commandBuffer, renderPassBeginInfo, vulkan.SubpassContentsInline)
}

func (fb *vkFrameBuffers) renderPassEnd(commandBuffer vulkan.CommandBuffer) {
	vulkan.CmdEndRenderPass(commandBuffer)
}

func (fb *vkFrameBuffers) setClearColor() {
	// todo: color input (color glm.Color)
	fb.clearColor = vulkan.ClearValue{0, 0, 0, 0}
}

func (fb *vkFrameBuffers) getClearColor() []vulkan.ClearValue {
	return []vulkan.ClearValue{fb.clearColor}
}
