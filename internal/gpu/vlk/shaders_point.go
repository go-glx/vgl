package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/simple"
)

const (
	shaderPointVertexCount = 1
	shaderPointVertexSize  = glm.SizeOfVec2 + glm.SizeOfVec4
)

func defaultShaderPoint() *shader.Meta {
	return shader.NewMeta(
		buildInShaderPoint,
		simple.CodeVertex(),
		simple.CodeFragment(),
		vulkan.PrimitiveTopologyPointList, false,
		simple.Bindings(shaderPointVertexSize),
		simple.Attributes(),
		shaderPointVertexCount,
		[]uint16{0},
	)
}

type dataPoint struct {
	vertex glm.Vec2
	color  glm.Vec4
}

func (d *dataPoint) BindingData() []byte {
	buff := make([]byte, 0, shaderPointVertexSize*shaderPointVertexCount)
	buff = append(buff, d.vertex.Data()...)
	buff = append(buff, d.color.Data()...)

	return buff
}

func (d *dataPoint) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
