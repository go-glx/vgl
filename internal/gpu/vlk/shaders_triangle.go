package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const (
	triangleVertexSize = glm.SizeOfVec2 + glm.SizeOfVec4
)

func defaultShaderTriangle() *shader.Meta {
	const bindingData = 0
	const locationVertex = 0
	const locationColor = 1

	return shader.NewMeta(
		buildInShaderTriangle,
		triangleVert,
		triangleFrag,
		vulkan.PrimitiveTopologyTriangleList,
		[]vulkan.VertexInputBindingDescription{
			{
				Binding:   bindingData,
				Stride:    triangleVertexSize,
				InputRate: vulkan.VertexInputRateVertex,
			},
		},
		[]vulkan.VertexInputAttributeDescription{ // [x,y,r,g,b],..
			{
				Location: locationVertex,
				Binding:  bindingData,
				Format:   bindingFormatVec2, // x, y
				Offset:   0,
			},
			{
				Location: locationColor,
				Binding:  bindingData,
				Format:   bindingFormatVec4, // r, g, b, a
				Offset:   glm.SizeOfVec2,
			},
		},
	)
}

type dataTriangle struct {
	vertexes [3]glm.Vec2
	colors   [3]glm.Vec4
	filled   bool
}

func (d *dataTriangle) BindingData() []byte {
	const vertexCount = 3
	buff := make([]byte, 0, vertexCount*triangleVertexSize)

	for i := 0; i < vertexCount; i++ {
		buff = append(buff, d.vertexes[i].Data()...)
		buff = append(buff, d.colors[i].Data()...)
	}

	return buff
}

func (d *dataTriangle) Indexes() []uint16 {
	return []uint16{0, 1, 2}
}

func (d *dataTriangle) PolygonMode() vulkan.PolygonMode {
	if d.filled {
		return vulkan.PolygonModeFill
	}

	return vulkan.PolygonModeLine
}
