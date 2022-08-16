package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/circle"
)

const shaderCircleVertexCount = 1

func defaultShaderCircle() *shader.Meta {
	return shader.NewMeta(
		buildInShaderCircle,
		circle.CodeVertex(),
		circle.CodeFragment(),
		vulkan.PrimitiveTopologyTriangleList, false,
		circle.Bindings(),
		circle.Attributes(),
		shaderCircleVertexCount,
		[]uint16{0, 0, 0, 0, 0, 0}, // only len is matter here
	)
}

type dataCircle struct {
	center glm.Vec2
	radius glm.Vec2
	// colors    [4]glm.Vec4
	// thickness [4]glm.Vec1
	// smooth    [4]glm.Vec1
}

func (d *dataCircle) BindingData() []byte {
	buff := make([]byte, 0, shaderCircleVertexCount*circle.VertexSize)
	buff = append(buff, d.center.Data()...)
	buff = append(buff, d.radius.Data()...)

	return buff
}

func (d *dataCircle) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
