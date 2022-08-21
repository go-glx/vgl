package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/circle"
)

const shaderCircleVertexCount = 5

func defaultShaderCircle() *shader.Meta {
	return shader.NewMeta(
		buildInShaderCircle,
		circle.CodeVertex(),
		circle.CodeFragment(),
		vulkan.PrimitiveTopologyTriangleList, false,
		circle.Bindings(),
		circle.Attributes(),
		shaderCircleVertexCount,
		true,

		// 0         1
		// # ------- #
		// |  \    / |
		// |    #4   |
		// |  /   \  |
		// # ------- #
		// 3         2
		[]uint16{
			0, 1, 4,
			1, 2, 4,
			2, 3, 4,
			3, 0, 4,
		},
	)
}

type dataCircle struct {
	pos [5]glm.Vec2
	// colors    [4]glm.Vec4
	// thickness [4]glm.Vec1
	// smooth    [4]glm.Vec1
}

func (d *dataCircle) BindingData() []byte {
	buff := make([]byte, 0, shaderCircleVertexCount*circle.VertexSize)

	for i := 0; i < shaderCircleVertexCount; i++ {
		buff = append(buff, d.pos[i].Data()...)
	}

	return buff
}

func (d *dataCircle) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
