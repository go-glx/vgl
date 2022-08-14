package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/circle"
)

const shaderCircleVertexCount = 4

func defaultShaderCircle() *shader.Meta {
	return shader.NewMeta(
		buildInShaderCircle,
		circle.CodeVertex(),
		circle.CodeFragment(),
		vulkan.PrimitiveTopologyTriangleList, false,
		circle.Bindings(),
		circle.Attributes(),
		shaderCircleVertexCount,
		[]uint16{0, 1, 2, 2, 3, 0}, // todo: can be reused with quad textures (same index for two shaders)
	)
}

type dataCircle struct {
	vertexes  [4]glm.Vec2
	colors    [4]glm.Vec4
	thickness [4]glm.Vec1
	smooth    [4]glm.Vec1
}

func (d *dataCircle) BindingData() []byte {
	buff := make([]byte, 0, shaderCircleVertexCount*circle.VertexSize)

	for i := 0; i < shaderCircleVertexCount; i++ {
		buff = append(buff, d.vertexes[i].Data()...)
		buff = append(buff, d.colors[i].Data()...)
		buff = append(buff, d.thickness[i].Data()...)
		buff = append(buff, d.smooth[i].Data()...)
	}

	return buff
}

func (d *dataCircle) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
