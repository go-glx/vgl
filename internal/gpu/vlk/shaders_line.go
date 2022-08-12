package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/simple"
)

const (
	shaderLineVertexCount = 2
	shaderLineVertexSize  = glm.SizeOfVec2 + glm.SizeOfVec4
)

func defaultShaderLine() *shader.Meta {
	return shader.NewMeta(
		buildInShaderLine,
		simple.CodeVertex(),
		simple.CodeFragment(),
		vulkan.PrimitiveTopologyLineList, false,
		simple.Bindings(shaderLineVertexSize),
		simple.Attributes(),
		shaderLineVertexCount,
		[]uint16{0, 1},
	)
}

type dataLine struct {
	vertexes [2]glm.Vec2
	colors   [2]glm.Vec4
	width    float32 // todo
}

func (d *dataLine) BindingData() []byte {
	buff := make([]byte, 0, shaderLineVertexCount*shaderLineVertexSize)

	for i := 0; i < shaderLineVertexCount; i++ {
		buff = append(buff, d.vertexes[i].Data()...)
		buff = append(buff, d.colors[i].Data()...)
	}

	return buff
}

func (d *dataLine) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeLine
}
