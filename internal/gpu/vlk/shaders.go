package vlk

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"
)

const (
	buildInShaderTriangle = "triangle"
)

const (
	bindingFormatVec2 = vulkan.FormatR32g32Sfloat
	bindingFormatVec3 = vulkan.FormatR32g32b32Sfloat
	bindingFormatVec4 = vulkan.FormatR32g32b32a32Sfloat
)

var (
	//go:embed shaders/triangle.vert.spv
	triangleVert []byte
	//go:embed shaders/triangle.frag.spv
	triangleFrag []byte
)

type shaderData interface {
	BindingData() []byte
	Indexes() []uint16
}
