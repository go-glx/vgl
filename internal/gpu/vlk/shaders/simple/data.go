package simple

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

const (
	simpleShaderBindingMainID                 = 0
	simpleShaderBindingMainIDLocationPosition = 0
	simpleShaderBindingMainIDLocationColor    = 1
)

var (
	//go:embed vert.spv
	codeVert []byte
	//go:embed frag.spv
	codeFrag []byte
)

func CodeVertex() []byte {
	return codeVert
}

func CodeFragment() []byte {
	return codeFrag
}

func Bindings(stride uint32) []vulkan.VertexInputBindingDescription {
	return []vulkan.VertexInputBindingDescription{
		{
			Binding:   0,
			Stride:    stride,
			InputRate: vulkan.VertexInputRateVertex,
		},
	}
}

func Attributes() []vulkan.VertexInputAttributeDescription {
	return []vulkan.VertexInputAttributeDescription{
		{
			Location: simpleShaderBindingMainIDLocationPosition,
			Binding:  simpleShaderBindingMainID,
			Format:   vulkan.FormatR32g32Sfloat, // x, y
			Offset:   0,
		},
		{
			Location: simpleShaderBindingMainIDLocationColor,
			Binding:  simpleShaderBindingMainID,
			Format:   vulkan.FormatR32g32b32a32Sfloat, // r, g, b, a
			Offset:   glm.SizeOfVec2,
		},
	}
}
