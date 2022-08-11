package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/simple"
)

const (
	shaderTriangleVertexCount = 3
	shaderTriangleVertexSize  = glm.SizeOfVec2 + glm.SizeOfVec4
)

func defaultShaderTriangle() *shader.Meta {
	return shader.NewMeta(
		buildInShaderTriangle,
		simple.CodeVertex(),
		simple.CodeFragment(),
		vulkan.PrimitiveTopologyTriangleList, false,
		simple.Bindings(shaderTriangleVertexSize),
		simple.Attributes(),
		shaderTriangleVertexCount,
		[]uint16{0, 1, 2},
	)
}

type dataTriangle struct {
	vertexes [3]glm.Vec2
	colors   [3]glm.Vec4
	filled   bool
}

func (d *dataTriangle) BindingData() []byte {
	buff := make([]byte, 0, shaderTriangleVertexCount*shaderTriangleVertexSize)

	for i := 0; i < shaderTriangleVertexCount; i++ {
		buff = append(buff, d.vertexes[i].Data()...)
		buff = append(buff, d.colors[i].Data()...)
	}

	return buff
}

func (d *dataTriangle) PolygonMode() vulkan.PolygonMode {
	if d.filled {
		return vulkan.PolygonModeFill
	}

	return vulkan.PolygonModeLine
}
