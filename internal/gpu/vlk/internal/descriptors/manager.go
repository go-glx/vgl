package descriptors

import (
	"bytes"
	"fmt"
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

const framesCount = def.OptimalSwapChainBuffersCount

type (
	Manager struct {
		logger config.Logger
		ld     *logical.Device
		pool   *Pool
		heap   *alloc.Heap

		uniformBufferAlignSize   uint32
		blueprint                *Blueprint
		allocatedSets            allocatedSets
		frameAllocationGlobalUBO [framesCount]alloc.Allocation
	}

	frameID       = uint8
	allocatedSets map[frameID]layoutSet
	layoutSet     map[layoutIndex]vulkan.DescriptorSet

	uniformData interface {
		Data() []byte
	}

	uniformBlock struct {
		uniforms     []uniformData
		expectedSize int
	}
)

func NewManager(
	logger config.Logger,
	ld *logical.Device,
	pd *physical.Device,
	pool *Pool,
	heap *alloc.Heap,
	blueprint *Blueprint,
) *Manager {
	return &Manager{
		logger: logger,
		ld:     ld,
		pool:   pool,
		heap:   heap,

		uniformBufferAlignSize: uint32(pd.PrimaryGPU().Props.Limits.MinUniformBufferOffsetAlignment),
		blueprint:              blueprint,
		allocatedSets:          allocateSets(ld, pool, blueprint),
	}
}

func (m *Manager) UpdateGlobalUBO(frameID uint8, view, projection glm.Mat4, surfaceSize glm.Vec2) Data {
	// clear previously allocated buffer for this frame (if exist)
	if m.frameAllocationGlobalUBO[frameID].Valid {
		m.heap.Free(m.frameAllocationGlobalUBO[frameID])
	}

	uniforms := []uniformBlock{
		{
			uniforms:     []uniformData{&view, &projection},
			expectedSize: glm.SizeOfMat4 * 2,
		},
		{
			uniforms:     []uniformData{&surfaceSize},
			expectedSize: glm.SizeOfVec2,
		},
	}

	staging, offsets := m.prepareUniformBuffer(uniforms)

	// write data to memory
	allocUBO := m.heap.Write(
		staging,
		alloc.BufferTypeUniform,
		alloc.StorageTargetCoherent,
		alloc.FlagsNone,
	)

	m.frameAllocationGlobalUBO[frameID] = allocUBO

	bufferBindingInfo := make([]vulkan.DescriptorBufferInfo, 0, len(uniforms))
	for index, uniform := range uniforms {
		bufferBindingInfo = append(bufferBindingInfo, vulkan.DescriptorBufferInfo{
			Buffer: allocUBO.Buffer,
			Offset: allocUBO.Offset + offsets[index],
			Range:  vulkan.DeviceSize(uniform.expectedSize),
		})
	}

	// prepare descriptor write
	descriptorSet := m.allocatedSets[frameID][layoutIndexGlobal]
	blueprint := m.blueprint.LayoutGlobal()
	writeSets := make([]vulkan.WriteDescriptorSet, 0, len(blueprint.bindings))

	for _, binding := range blueprint.bindings {
		writeSets = append(writeSets, vulkan.WriteDescriptorSet{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          descriptorSet,
			DstBinding:      binding.Binding,
			DstArrayElement: 0,
			DescriptorCount: binding.DescriptorCount,
			DescriptorType:  binding.DescriptorType,
			PBufferInfo: []vulkan.DescriptorBufferInfo{
				bufferBindingInfo[binding.Binding],
			},
		})
	}

	// update references
	vulkan.UpdateDescriptorSets(m.ld.Ref(), uint32(len(writeSets)), writeSets, 0, nil)

	return Data{
		DescriptorSet: descriptorSet,
	}
}

// prepareUniformBuffer will write all block bytes data to single slice
// and align bytes in each block of zeroed space if needed
// function return result bytes slice and offset of each block start
func (m *Manager) prepareUniformBuffer(blocks []uniformBlock) ([]byte, []vulkan.DeviceSize) {
	totalSize := 0
	for _, block := range blocks {
		totalSize += m.alignUniformSize(block.expectedSize)
	}

	buffer := make([]byte, 0, totalSize)
	offsets := make([]vulkan.DeviceSize, 0, len(blocks))

	for _, block := range blocks {
		data := make([]byte, 0, block.expectedSize)
		for _, uniform := range block.uniforms {
			data = append(data, uniform.Data()...)
		}

		if len(data) != block.expectedSize {
			panic(fmt.Errorf("uniform block size miss calculated. expected=%dB, actual=%dB",
				block.expectedSize,
				len(data),
			))
		}

		alignedSize := m.alignUniformSize(block.expectedSize)
		uselessSize := alignedSize - block.expectedSize

		// write actual data
		offsets = append(offsets, vulkan.DeviceSize(len(buffer))) // current buffer len = offset of next block
		buffer = append(buffer, data...)

		// write trash data if needed, we need match device aligned size exactly
		if uselessSize > 0 {
			buffer = append(buffer, bytes.Repeat([]byte{0}, uselessSize)...)
		}
	}

	return buffer, offsets
}

func (m *Manager) alignUniformSize(realSize int) int {
	return int(math.Ceil(float64(realSize)/float64(m.uniformBufferAlignSize)) * float64(m.uniformBufferAlignSize))
}

func allocateSets(ld *logical.Device, pool *Pool, blueprint *Blueprint) allocatedSets {
	sets := make(allocatedSets)

	for frameID := frameID(0); frameID < framesCount; frameID++ {
		set := make(layoutSet)
		set[layoutIndexGlobal] = allocateSetGlobal(ld, pool, blueprint.LayoutGlobal())

		sets[frameID] = set
	}

	return sets
}

func allocateSetGlobal(ld *logical.Device, pool *Pool, bpLayout BlueprintLayout) vulkan.DescriptorSet {
	setAllocateInfo := vulkan.DescriptorSetAllocateInfo{
		SType:              vulkan.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     pool.Pool(),
		DescriptorSetCount: 1,
		PSetLayouts: []vulkan.DescriptorSetLayout{
			bpLayout.layout,
		},
	}

	var set vulkan.DescriptorSet
	must.Work(vulkan.AllocateDescriptorSets(ld.Ref(), &setAllocateInfo, &set))

	return set
}
