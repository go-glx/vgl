package frame

import "github.com/vulkan-go/vulkan"

func (m *Manager) renderPassMainBegin(cb vulkan.CommandBuffer) {
	renderPassBeginInfo := &vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  m.mainRenderPass.Ref(),
		Framebuffer: m.chain.FrameBuffer(int(m.id)),
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{
				X: 0,
				Y: 0,
			},
			Extent: vulkan.Extent2D{
				Width:  m.chain.Props().BufferSize.Width,
				Height: m.chain.Props().BufferSize.Height,
			},
		},
		ClearValueCount: 1,
		PClearValues: []vulkan.ClearValue{
			{0, 0, 0, 0},
		},
	}

	vulkan.CmdBeginRenderPass(cb, renderPassBeginInfo, vulkan.SubpassContentsInline)
}

func (m *Manager) renderPassMainEnd(cb vulkan.CommandBuffer) {
	vulkan.CmdEndRenderPass(cb)
}
