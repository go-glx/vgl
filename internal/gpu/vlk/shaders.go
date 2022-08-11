package vlk

import (
	_ "embed"
	"fmt"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const (
	buildInShaderPoint    = "point"
	buildInShaderTriangle = "triangle"
	buildInShaderRect     = "rect"
)

// shaders that will preload indexes to fast-buffer
// this allows instancing them in one draw call
var preloadIndexShaders = []string{
	buildInShaderPoint,
	buildInShaderTriangle,
	buildInShaderRect,
}

func (vlk *VLK) preloadShaderIndexes(shader *shader.Shader) {
	shaderID := shader.Meta().ID()
	heap := vlk.cont.allocBuffers()

	vlk.cont.logger.Debug(fmt.Sprintf("preload shader '%s' indexes", shaderID))

	// create new index buffer for this shader
	// and pre-generate index data for N instances

	indexes := make([]byte, 0, len(shader.Meta().Indexes())*2) // uint16 -> uint32

	// index is something like [0,1,2] for one instance (triangle for example)
	// we want to populate index buffer for N instances, for example when N=3
	// buffer will be equal to [0,1,2, 3,4,5, 6,7,8]
	// this allows to draw at least 3 instance of same triangle in one draw call

	for inst := uint32(0); inst < def.BufferIndexMaxInstances; inst++ {
		for _, index := range shader.Meta().Indexes() {
			var instanceIndex uint32

			if index == 0xffff {
				// do not change offset of breaking sequence
				// this is special index (65535 = 0xffff)
				instanceIndex = uint32(index)
			} else {
				offset := shader.Meta().VertexCount() * inst
				instanceIndex = offset + uint32(index)
			}

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
