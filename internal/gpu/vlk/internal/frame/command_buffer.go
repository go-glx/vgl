package frame

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func (m *Manager) commandBufferBegin() {
	must.Work(
		vulkan.BeginCommandBuffer(m.commandBuffers[m.frameID], &vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
		}),
	)
}

func (m *Manager) commandBufferEnd() {
	must.Work(
		vulkan.EndCommandBuffer(m.commandBuffers[m.frameID]),
	)
}
