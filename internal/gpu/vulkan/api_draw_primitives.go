package vulkan

import (
	"github.com/fe3dback/govgl/glm"
	"github.com/fe3dback/govgl/internal/gpu/vulkan/internal/shader/shaderm"
)

func (vk *Vk) DrawTmpTriangle() {
	for i := float32(-1); i < 1; i += 0.1 {
		vk.appendToRenderQueue(&shaderm.Triangle{
			Position: [3]glm.Vec2{
				{X: i, Y: -0.5},
				{X: 0.5, Y: 0.5},
				{X: -0.5, Y: 0.5},
			},
			Color: [3]glm.Vec3{
				{R: (i + 1) / 2, G: 0, B: 0},
				{R: 0, G: 1, B: 0},
				{R: 0, G: 0, B: 1},
			},
		})
	}
}

func (vk *Vk) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	vk.appendToRenderQueue(&shaderm.Rect{
		Position: vertexPos,
		Color:    vertexColor,
	})
}
