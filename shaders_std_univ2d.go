package vgl

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl/internal/shaders"
)

var universal2dBindings = []ParamsRegisterShaderInputVertexBinding{
	{
		// position vec2 x,y
		Location: 0,
		Size:     glx.SizeOfVec2,
		Format:   vulkan.FormatR32g32Sfloat,
	},
	{
		// color vec4 r,g,b,a
		Location: 1,
		Size:     glx.SizeOfVec4,
		Format:   vulkan.FormatR32g32b32a32Sfloat,
	},
}

var (
	stdShaderPoint = ParamsRegisterShader{
		ShaderName:       buildInShaderPoint,
		ProgramVert:      shaders.Universal2DVertSpv(),
		ProgramFrag:      shaders.Universal2DFragSpv(),
		Topology:         vulkan.PrimitiveTopologyPointList,
		TopologyRestarts: false,
		InputLayout: ParamsRegisterShaderInputLayout{
			VertexCount:   1,
			VertexBinding: universal2dBindings,
			Indexes:       []uint16{0},
		},
	}

	stdShaderLine = ParamsRegisterShader{
		ShaderName:       buildInShaderLine,
		ProgramVert:      shaders.Universal2DVertSpv(),
		ProgramFrag:      shaders.Universal2DFragSpv(),
		Topology:         vulkan.PrimitiveTopologyLineList,
		TopologyRestarts: false,
		InputLayout: ParamsRegisterShaderInputLayout{
			VertexCount:   2,
			VertexBinding: universal2dBindings,
			Indexes:       []uint16{0, 1},
		},
	}

	stdShaderTriangle = ParamsRegisterShader{
		ShaderName:       buildInShaderTriangle,
		ProgramVert:      shaders.Universal2DVertSpv(),
		ProgramFrag:      shaders.Universal2DFragSpv(),
		Topology:         vulkan.PrimitiveTopologyTriangleList,
		TopologyRestarts: false,
		InputLayout: ParamsRegisterShaderInputLayout{
			VertexCount:   3,
			VertexBinding: universal2dBindings,
			Indexes:       []uint16{0, 1, 2},
		},
	}

	stdShaderRect = ParamsRegisterShader{
		ShaderName:       buildInShaderRect,
		ProgramVert:      shaders.Universal2DVertSpv(),
		ProgramFrag:      shaders.Universal2DFragSpv(),
		Topology:         vulkan.PrimitiveTopologyLineStrip,
		TopologyRestarts: true,
		InputLayout: ParamsRegisterShaderInputLayout{
			VertexCount:   4,
			VertexBinding: universal2dBindings,
			Indexes:       []uint16{0, 1, 2, 3, 0, 0xffff},
		},
	}
)

type (
	shaderInputUniversal2d struct {
		vertexes []shaderInputUniversal2dVertex
	}

	shaderInputUniversal2dVertex struct {
		pos   glx.Vec2
		color glx.Vec4
	}
)

func (d *shaderInputUniversal2d) VertexData() []byte {
	const vertSize = glx.SizeOfVec2 + glx.SizeOfVec4
	buff := make([]byte, 0, len(d.vertexes)*vertSize)

	for _, vertex := range d.vertexes {
		buff = append(buff, vertex.pos.Data()...)
		buff = append(buff, vertex.color.Data()...)
	}

	return buff
}

func (d *shaderInputUniversal2d) StorageData() []byte {
	return nil
}
