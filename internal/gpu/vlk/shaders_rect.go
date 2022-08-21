package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
	"github.com/go-glx/vgl/internal/gpu/vlk/shaders/simple"
)

const (
	shaderRectVertexCount = 4
)

func defaultShaderRect() *shader.Meta {
	return shader.NewMeta(
		buildInShaderRect,
		simple.CodeVertex(),
		simple.CodeFragment(),
		vulkan.PrimitiveTopologyLineStrip, true,
		simple.Bindings(),
		simple.Attributes(),
		shaderRectVertexCount,
		true,
		[]uint16{0, 1, 2, 3, 0, 0xffff},
	)
}

type dataRectOutline struct {
	vertexes [4]glm.Vec2
	colors   [4]glm.Vec4
}

func (d *dataRectOutline) BindingData() []byte {
	buff := make([]byte, 0, shaderRectVertexCount*simple.VertexSize)

	for i := 0; i < shaderRectVertexCount; i++ {
		buff = append(buff, d.vertexes[i].Data()...)
		buff = append(buff, d.colors[i].Data()...)
	}

	return buff
}

func (d *dataRectOutline) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeLine
}
