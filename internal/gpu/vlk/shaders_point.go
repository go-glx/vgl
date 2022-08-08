package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const (
	pointVertexSize = glm.SizeOfVec2 + glm.SizeOfVec4
)

func defaultShaderPoint() *shader.Meta {
	const bindingData = 0
	const locationVertex = 0
	const locationColor = 1

	return shader.NewMeta(
		buildInShaderPoint,
		simpleVert,
		simpleFrag,
		vulkan.PrimitiveTopologyPointList,
		[]vulkan.VertexInputBindingDescription{
			{
				Binding:   bindingData,
				Stride:    pointVertexSize,
				InputRate: vulkan.VertexInputRateVertex,
			},
		},
		[]vulkan.VertexInputAttributeDescription{ // [x,y,r,g,b,a],..
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

type dataPoint struct {
	vertex glm.Vec2
	color  glm.Vec4
}

func (d *dataPoint) BindingData() []byte {
	buff := make([]byte, 0, pointVertexSize)
	buff = append(buff, d.vertex.Data()...)
	buff = append(buff, d.color.Data()...)

	return buff
}

func (d *dataPoint) Indexes() []uint16 {
	return []uint16{0}
}

func (d *dataPoint) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
