package vgl

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/shaders"
)

var (
	stdShaderCircle = ParamsRegisterShader{
		ShaderName:       buildInShaderCircle,
		ProgramVert:      shaders.Circle2DVertSpv(),
		ProgramFrag:      shaders.Circle2DFragSpv(),
		Topology:         vulkan.PrimitiveTopologyTriangleList,
		TopologyRestarts: false,
		InputLayout: ParamsRegisterShaderInputLayout{
			VertexCount: 4,
			VertexBinding: []ParamsRegisterShaderInputVertexBinding{
				{
					// x, y
					Location: 0,
					Size:     glm.SizeOfVec2,
					Format:   vulkan.FormatR32g32Sfloat,
				},
				{
					// r, g, b, a
					Location: 1,
					Size:     glm.SizeOfVec4,
					Format:   vulkan.FormatR32g32b32a32Sfloat,
				},
			},
			Indexes: []uint16{0, 1, 2, 2, 3, 0},
		},
	}
)

type (
	shaderInputCircle2d struct {
		vertexes  []shaderInputCircle2dVertex
		center    glm.Vec2
		radius    glm.Vec1
		thickness glm.Vec1
		smooth    glm.Vec1
	}

	shaderInputCircle2dVertex struct {
		pos   glm.Vec2
		color glm.Vec4
	}
)

func (d *shaderInputCircle2d) VertexData() []byte {
	const vertSize = glm.SizeOfVec2 + glm.SizeOfVec4
	buff := make([]byte, 0, stdShaderCircle.InputLayout.VertexCount*vertSize)

	for _, vertex := range d.vertexes {
		buff = append(buff, vertex.pos.Data()...)
		buff = append(buff, vertex.color.Data()...)
	}

	return buff
}

func (d *shaderInputCircle2d) StorageData() []byte {
	buff := make([]byte, 0, glm.SizeOfVec2+(glm.SizeOfVec1*3))

	buff = append(buff, d.center.Data()...)
	buff = append(buff, d.radius.Data()...)
	buff = append(buff, d.thickness.Data()...)
	buff = append(buff, d.smooth.Data()...)

	return buff
}

func (d *shaderInputCircle2d) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
