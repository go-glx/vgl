package simple

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

const (
	shaderBindingMainID                 = 0
	shaderBindingMainIDLocationPosition = 0
	shaderBindingMainIDLocationColor    = 1
)

const (
	// VertexSize = pos + color
	VertexSize = glm.SizeOfVec2 + glm.SizeOfVec4
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

func Bindings() []vulkan.VertexInputBindingDescription {
	return []vulkan.VertexInputBindingDescription{
		{
			Binding:   0,
			Stride:    VertexSize,
			InputRate: vulkan.VertexInputRateVertex,
		},
	}
}

func Attributes() []vulkan.VertexInputAttributeDescription {
	return []vulkan.VertexInputAttributeDescription{
		{
			Location: shaderBindingMainIDLocationPosition,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32g32Sfloat, // x, y
			Offset:   0,
		},
		{
			Location: shaderBindingMainIDLocationColor,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32g32b32a32Sfloat, // r, g, b, a
			Offset:   glm.SizeOfVec2,
		},
	}
}
