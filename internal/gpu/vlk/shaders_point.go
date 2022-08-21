package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/simple"
)

func defaultShaderPoint() *shader.Meta {
	return shader.NewMeta(
		buildInShaderPoint,
		simple.CodeVertex(),
		simple.CodeFragment(),
		vulkan.PrimitiveTopologyPointList, false,
		simple.Bindings(),
		simple.Attributes(),
		1,
		true,
		[]uint16{0},
	)
}

type dataPoint struct {
	vertex glm.Vec2
	color  glm.Vec4
}

func (d *dataPoint) BindingData() []byte {
	buff := make([]byte, 0, simple.VertexSize)
	buff = append(buff, d.vertex.Data()...)
	buff = append(buff, d.color.Data()...)

	return buff
}

func (d *dataPoint) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
