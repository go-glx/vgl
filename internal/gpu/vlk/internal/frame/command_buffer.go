package frame

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func (m *Manager) commandBufferBegin(ctx Context) {
	must.Work(
		vulkan.BeginCommandBuffer(m.commandBuffers[ctx.frameID], &vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
		}),
	)
}

func (m *Manager) commandBufferEnd(ctx Context) {
	must.Work(
		vulkan.EndCommandBuffer(m.commandBuffers[ctx.frameID]),
	)
}
