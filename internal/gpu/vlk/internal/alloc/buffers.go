package alloc

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type Chunk struct {
	Buffer         vulkan.Buffer
	BufferOffset   vulkan.DeviceSize
	InstancesCount uint32
}

type Buffers struct {
	heap *Heap
}

func NewBuffers(heap *Heap) *Buffers {
	return &Buffers{
		heap: heap,
	}
}

func (b *Buffers) WriteIndexData(data []byte) Allocation {
	return b.heap.Write(
		data,
		BufferTypeIndex,
		StorageTargetImmutable,
		FlagsNone,
	)
}

func (b *Buffers) WriteInstancesVertexes(instances []shader.InstanceData) []Chunk {
	if len(instances) == 0 {
		return nil
	}

	staging := make([]byte, 0, def.BufferVertexSizeBytes)
	spaceLeft := cap(staging)

	chunks := make([]Chunk, 0)

	flush := func() {
		// write staging to GPU
		alloc := b.heap.Write(
			staging,
			BufferTypeVertex,
			StorageTargetCoherent,
			FlagsTemporary,
		)

		// clear staging data
		staging = make([]byte, 0, def.BufferVertexSizeBytes)
		spaceLeft = cap(staging)

		// add chunk with GPU ptr
		chunks = append(chunks, Chunk{
			Buffer:         alloc.Buffer,
			BufferOffset:   alloc.Offset,
			InstancesCount: uint32(len(instances)),
		})
	}

	for _, instance := range instances {
		data := instance.BindingData()
		size := len(data)

		if size > spaceLeft {
			flush()
		}

		staging = append(staging, data...)
		spaceLeft -= size
	}

	if len(staging) > 0 {
		flush()
	}

	return chunks
}
