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
			VertexCount: 5,
			VertexBinding: []ParamsRegisterShaderInputVertexBinding{
				{
					// x, y
					Location: 0,
					Size:     glm.SizeOfVec2,
					Format:   vulkan.FormatR32g32Sfloat,
				},
				// {
				// 	// r, g, b, a
				// 	Location: 0,
				// 	Size:     glm.SizeOfVec4,
				// 	Format:   vulkan.FormatR32g32b32a32Sfloat,
				// },
				// {
				// 	// w - thickness
				// 	Location: 0,
				// 	Size:     glm.SizeOfVec1,
				// 	Format:   vulkan.FormatR32Sfloat,
				// },
				// {
				// 	// w - smooth
				// 	Location: 0,
				// 	Size:     glm.SizeOfVec1,
				// 	Format:   vulkan.FormatR32Sfloat,
				// },
			},
			Indexes: []uint16{
				// 0         1
				// # ------- #
				// |  \    / |
				// |    #4   |
				// |  /   \  |
				// # ------- #
				// 3         2
				0, 1, 4,
				1, 2, 4,
				2, 3, 4,
				3, 0, 4,
			},
			UseGlobalUniforms: true,
		},
	}
)

type (
	shaderInputCircle2d struct {
		vertexes  []shaderInputCircle2dVertex
		center    glm.Vec2
		thickness glm.Vec1
		smooth    glm.Vec1
	}

	shaderInputCircle2dVertex struct {
		pos glm.Vec2
	}
)

func (d *shaderInputCircle2d) BindingData() []byte {
	// const vertSize = glm.SizeOfVec2 + glm.SizeOfVec4 + glm.SizeOfVec1 + glm.SizeOfVec1 // todo?
	const vertSize = glm.SizeOfVec2
	buff := make([]byte, 0, len(d.vertexes)*vertSize)

	for _, vertex := range d.vertexes {
		buff = append(buff, vertex.pos.Data()...)
		// buff = append(buff, vertex.color.Data()...)
		// buff = append(buff, vertex.thickness.Data()...)
		// buff = append(buff, vertex.smooth.Data()...)
	}

	return buff
}

func (d *shaderInputCircle2d) PolygonMode() vulkan.PolygonMode {
	return vulkan.PolygonModeFill
}
