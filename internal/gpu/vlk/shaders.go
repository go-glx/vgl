package vlk

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"
)

const (
	buildInShaderPoint    = "point"
	buildInShaderTriangle = "triangle"
)

const (
	bindingFormatVec2 = vulkan.FormatR32g32Sfloat
	bindingFormatVec3 = vulkan.FormatR32g32b32Sfloat
	bindingFormatVec4 = vulkan.FormatR32g32b32a32Sfloat
)

var (
	//go:embed shaders/simple.vert.spv
	simpleVert []byte
	//go:embed shaders/simple.frag.spv
	simpleFrag []byte
)
