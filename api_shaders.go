package vgl

import (
	"github.com/vulkan-go/vulkan"
)

type (
	ParamsRegisterShader struct {
		// unique shader name
		ShaderName string
		// raw shader program in GLSL language for vertex shader
		ProgramVert []byte
		// raw shader program in GLSL language for fragment(pixel) shader
		ProgramFrag []byte

		// Topology define how GPU should draw vertexes
		// see: https://registry.khronos.org/vulkan/specs/1.3-extensions/man/html/VkPrimitiveTopology.html
		Topology vulkan.PrimitiveTopology

		// When true you need write special value to input indexes "0xffff"
		// to visually break drawing figures (lines, triangles, etc..)
		// its fully depend on Topology
		TopologyRestarts bool

		// InputLayout is specification (blueprint) of all kind of input information
		// that should be provided into shaders.
		InputLayout ParamsRegisterShaderInputLayout
	}

	ParamsRegisterShaderInputLayout struct {
		// how many vertexes has one instance (2d/3d figure/model)
		// for example:
		//   for drawing square we need 4 vertexes
		VertexCount uint32

		// layout for every location in vertex shader
		VertexBinding []ParamsRegisterShaderInputVertexBinding

		// index order for each vertex in clock-wise order.
		// for example:
		//   for drawing square with 2 triangles, indexes may be: [0,1,2,2,3,0]
		//
		//   0       1
		//   * ----- *
		//   |   \   |
		//   * ----- *
		//   3       2
		Indexes []uint16

		// will allow usage of global UBO layer=0
		// this UBO object contain global application values:
		//    vertex shader:
		//      set=0, binding = 0
		//        - view matrix (mat4)
		//        - projection matrix (mat4)
		//    index shader:
		//      set=0, binding = 1
		//        - surfaceSize (vec2) - current surface size (usually screen) in pixels (width/height)
		UseGlobalUniforms bool

		// additional layout for shader-wide local uniforms
		// provided input values for shader, will be written to
		// special uniform/storage buffer and auto bind to shader set with layout=1
		LocalUniforms []ParamsRegisterShaderInputLocalUniform
	}

	ParamsRegisterShaderInputVertexBinding struct {
		// Location ID in vertex shader for this input
		// example:
		//   layout(location = 0) in vec2 inPosition;
		//                     ^
		//                 this value
		Location uint32

		// Size of data used for one vertex
		//
		// example:
		//   layout(location = 0) in vec2 inPosition;
		//                            ^
		//                        this value
		//
		//   vec2 for example has 2 x float32
		//   each float32 has 4 bytes in size (32 / 8 bits)
		//   so Size of vec2 = 8
		Size uint32

		// Format of data
		// see: https://registry.khronos.org/vulkan/specs/1.3-extensions/man/html/VkFormat.html
		// for example, common formats:
		//   vec1 = vulkan.FormatR32Sfloat
		//   vec2 = vulkan.FormatR32g32Sfloat
		//   vec3 = vulkan.FormatR32g32b32Sfloat
		//   vec4 = vulkan.FormatR32g32b32a32Sfloat
		Format vulkan.Format
	}

	ParamsRegisterShaderInputLocalUniform struct {
		// todo
	}
)

// RegisterShader allow to use custom vert/frag shaders
// with various data/bindings/layout with automatic compilation
func (r *Render) RegisterShader(p *ParamsRegisterShader) {
	attributes := make([]vulkan.VertexInputAttributeDescription, 0)
	bindings := make([]vulkan.VertexInputBindingDescription, 0)

	strideSize := uint32(0)
	currentOffset := uint32(0)

	for _, binding := range p.InputLayout.VertexBinding {
		attributes = append(attributes, vulkan.VertexInputAttributeDescription{
			Location: binding.Location,
			Binding:  0,
			Format:   binding.Format,
			Offset:   currentOffset,
		})

		currentOffset += binding.Size
		strideSize += binding.Size
	}

	bindings = append(bindings, vulkan.VertexInputBindingDescription{
		Binding:   0,
		Stride:    strideSize,
		InputRate: vulkan.VertexInputRateVertex,
	})

	r.api.RegisterShader(
		p.ShaderName,
		p.ProgramVert,
		p.ProgramFrag,
		p.Topology,
		p.TopologyRestarts,
		bindings,
		attributes,
		p.InputLayout.VertexCount,
		p.InputLayout.Indexes,
	)
}
