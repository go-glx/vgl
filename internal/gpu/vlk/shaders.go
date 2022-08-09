package vlk

import (
	_ "embed"
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const (
	buildInShaderPoint    = "point"
	buildInShaderTriangle = "triangle"
)

// shaders that will preload indexes to fast-buffer
// this allows instancing them in one draw call
var preloadIndexShaders = []string{
	buildInShaderPoint,
	buildInShaderTriangle,
}

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

func (vlk *VLK) preloadShaderIndexes(shader *shader.Shader) {
	shaderID := shader.Meta().ID()
	heap := vlk.cont.allocBuffers()

	vlk.cont.logger.Debug(fmt.Sprintf("preload shader '%s' indexes", shaderID))

	// create new index buffer for this shader
	// and pre-generate index data for N instances

	size := uint32(len(shader.Meta().Indexes()))
	indexes := make([]byte, 0, size*2) // uint16 -> uint32

	// index is something like [0,1,2] for one instance (triangle for example)
	// we want to populate index buffer for N instances, for example when N=3
	// buffer will be equal to [0,1,2, 3,4,5, 6,7,8]
	// this allows to draw at least 3 instance of same triangle in one draw call

	for inst := uint32(0); inst < def.BufferIndexMapInstances; inst++ {
		for _, index := range shader.Meta().Indexes() {
			offset := size * inst
			instanceIndex := offset + uint32(index)

			indexes = append(indexes, uint8(instanceIndex&0xff), uint8(instanceIndex>>8))
		}
	}

	// this command will write indexes to GPU fast memory,
	// and later we will reuse this many times, because
	// indexes is not changed later in runtime

	allocationID := heap.AllocateIndexMemory(indexes)
	vlk.shaderIndexPtr[shaderID] = allocationID
}

func (vlk *VLK) indexBufferOf(shader *shader.Shader) alloc.Allocation {
	shaderID := shader.Meta().ID()
	heap := vlk.cont.allocBuffers()

	// return ptr for shader index buffer
	// with pre-generated data for N instances
	if allocationID, exist := vlk.shaderIndexPtr[shaderID]; exist {
		return heap.GetMemoryPointer(allocationID)
	}

	return alloc.Allocation{
		HasData: false,
	}
}
