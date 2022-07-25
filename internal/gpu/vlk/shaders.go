package vlk

import (
	_ "embed"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const (
	buildInShaderTriangle = "triangle"
)

var (
	//go:embed shaders/triangle.vert.spv
	triangleVert []byte
	//go:embed shaders/triangle.frag.spv
	triangleFrag []byte
)

func defaultShaderTriangle() *shader.Meta {
	return shader.NewMeta(
		buildInShaderTriangle,
		triangleVert,
		triangleFrag,
		vulkan.PrimitiveTopologyTriangleList,
		make([]vulkan.VertexInputBindingDescription, 0),
		make([]vulkan.VertexInputAttributeDescription, 0),
	)
}
