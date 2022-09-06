package vgl

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/glx"
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
					Size:     glx.SizeOfVec2,
					Format:   vulkan.FormatR32g32Sfloat,
				},
				{
					// r, g, b, a
					Location: 1,
					Size:     glx.SizeOfVec4,
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
		thickness glx.Vec1
		smooth    glx.Vec1
	}

	shaderInputCircle2dVertex struct {
		pos   glx.Vec2
		color glx.Vec4
	}
)

func (d *shaderInputCircle2d) VertexData() []byte {
	const vertSize = glx.SizeOfVec2 + glx.SizeOfVec4
	buff := make([]byte, 0, stdShaderCircle.InputLayout.VertexCount*vertSize)

	for _, vertex := range d.vertexes {
		buff = append(buff, vertex.pos.Data()...)
		buff = append(buff, vertex.color.Data()...)
	}

	return buff
}

func (d *shaderInputCircle2d) StorageData() []byte {
	buff := make([]byte, 0, glx.SizeOfVec1*2)

	// 8
	buff = append(buff, d.thickness.Data()...)
	buff = append(buff, d.smooth.Data()...)

	// example:
	// buff = append(buff, bytes.Repeat([]byte("0"), 4)...) // align to 8 bytes
	// todo: move align outside of shader definition
	// todo: 4 bytes is "magic" align, use device real align instead

	return buff
}
