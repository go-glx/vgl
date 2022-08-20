package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type Chunk struct {
	Buffer         vulkan.Buffer
	BufferOffset   vulkan.DeviceSize
	InstancesCount uint32
}

type Buffers struct {
	heap                   *Heap
	frameVertexAllocations [def.OptimalSwapChainBuffersCount][]Allocation
	frameUniformAllocation [def.OptimalSwapChainBuffersCount]Allocation
}

func NewBuffers(heap *Heap) *Buffers {
	return &Buffers{
		heap:                   heap,
		frameVertexAllocations: [def.OptimalSwapChainBuffersCount][]Allocation{},
		frameUniformAllocation: [def.OptimalSwapChainBuffersCount]Allocation{},
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

func (b *Buffers) ClearVertexBuffersOwnedBy(frameID uint32) {
	const defaultAllocationsCapacity = 8

	for _, allocation := range b.frameVertexAllocations[frameID] {
		b.heap.Free(allocation)
	}

	b.frameVertexAllocations[frameID] = make([]Allocation, 0, defaultAllocationsCapacity)
}

func (b *Buffers) WriteVertexBuffersFromInstances(frameID uint32, instances []shader.InstanceData) []Chunk {
	if len(instances) == 0 {
		return nil
	}

	// should be reasonable small for economy golang GC
	// in ideal it should equal avg size of all instances
	const defaultStagingCapacity = 512

	// todo: staging buffer can be reused between all this calls
	//       this will highly reduce GC time
	staging := make([]byte, 0, defaultStagingCapacity)

	// spaceLeft is soft limit, that equal to common buffer size
	// but in we try to write 100MB of vertex buffers in one call
	// this will alloc this 100MB anyway
	spaceLeft := def.BufferVertexSizeBytes

	chunks := make([]Chunk, 0)

	flush := func() {
		if len(staging) <= 0 {
			return
		}

		// write staging to GPU
		alloc := b.heap.Write(
			staging,
			BufferTypeVertex,
			StorageTargetCoherent,
			FlagsNone,
		)
		b.frameVertexAllocations[frameID] = append(b.frameVertexAllocations[frameID], alloc)

		// clear staging data
		staging = make([]byte, 0, defaultStagingCapacity)
		spaceLeft = def.BufferVertexSizeBytes

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

func (b *Buffers) UpdateUniformGlobalData(frameID uint32, model glm.Mat4, view glm.Mat4, projection glm.Mat4) {
	// clear previously allocated buffer for this frame (if exist)
	if b.frameUniformAllocation[frameID].Valid {
		b.heap.Free(b.frameUniformAllocation[frameID])
	}

	// write new data to same buffer
	staging := make([]byte, 0, glm.SizeOfMat4*3)
	staging = append(staging, model.Data()...)
	staging = append(staging, view.Data()...)
	staging = append(staging, projection.Data()...)

	// we always want this UBO data at offset=0
	// because of vulkan layout scheme, so not need to return anything here
	alloc := b.heap.Write(
		staging,
		BufferTypeUniform,
		StorageTargetCoherent,
		FlagsNone,
	)
	b.frameUniformAllocation[frameID] = alloc
	fmt.Printf("ubo offset=%d size=%d buffID=%d allocID=%d\n", alloc.Offset, alloc.Size, alloc.buffID, alloc.allocID)
}
