package circle

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

const (
	shaderBindingMainID                     = 0
	shaderBindingMainIDLocationPosition     = 0
	shaderBindingMainIDLocationColor        = 1
	shaderBindingMainIDLocationThickness    = 2
	shaderBindingMainIDLocationBorderSmooth = 3
)

const (
	// VertexSize =  pos + color + Thickness + BorderSmooth
	VertexSize = glm.SizeOfVec2 + glm.SizeOfVec4 + glm.SizeOfVec1 + glm.SizeOfVec1
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
		{
			Location: shaderBindingMainIDLocationThickness,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32Sfloat,          // thickness
			Offset:   glm.SizeOfVec2 + glm.SizeOfVec4, // todo: ??
		},
		{
			Location: shaderBindingMainIDLocationBorderSmooth,
			Binding:  shaderBindingMainID,
			Format:   vulkan.FormatR32Sfloat,                           // smooth
			Offset:   glm.SizeOfVec2 + glm.SizeOfVec4 + glm.SizeOfVec1, // todo: ??
		},
	}
}
