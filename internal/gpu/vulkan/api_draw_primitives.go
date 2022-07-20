package vulkan

import (
	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vulkan/internal/shader/shaderm"
)

func (vk *Vk) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	vk.appendToRenderQueue(&shaderm.Rect{
		Position: vertexPos,
		Color:    vertexColor,
	})
}
