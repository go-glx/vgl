package circle

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

const (
	shaderBindingMainID               = 0
	shaderBindingMainIDLocationCenter = 0
	shaderBindingMainIDLocationRadius = 1
)

const (
	// VertexSize =  center(xy) + radius(xy)
	VertexSize = glm.SizeOfVec2 + glm.SizeOfVec2
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
			Location: shaderBindingMainIDLocationCenter,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32g32Sfloat, // x, y
			Offset:   0,
		},
		{
			Location: shaderBindingMainIDLocationRadius,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32g32Sfloat, // x, y
			Offset:   glm.SizeOfVec2,
		},
	}
}
