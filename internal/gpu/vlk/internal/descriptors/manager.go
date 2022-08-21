package descriptors

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

const framesCount = def.OptimalSwapChainBuffersCount

type (
	Manager struct {
		logger config.Logger
		ld     *logical.Device
		pool   *Pool
		heap   *alloc.Heap

		blueprint                *Blueprint
		allocatedSets            allocatedSets
		frameAllocationGlobalUBO [framesCount]alloc.Allocation
	}

	frameID       = uint8
	allocatedSets map[frameID]layoutSet
	layoutSet     map[layoutIndex]vulkan.DescriptorSet
)

func NewManager(
	logger config.Logger,
	ld *logical.Device,
	pool *Pool,
	heap *alloc.Heap,
	blueprint *Blueprint,
) *Manager {
	return &Manager{
		logger: logger,
		ld:     ld,
		pool:   pool,
		heap:   heap,

		blueprint:     blueprint,
		allocatedSets: allocateSets(ld, pool, blueprint),
	}
}

func (m *Manager) UpdateGlobalUBO(frameID uint8, view, projection glm.Mat4, surfaceSize glm.Vec2) Data {
	// clear previously allocated buffer for this frame (if exist)
	if m.frameAllocationGlobalUBO[frameID].Valid {
		m.heap.Free(m.frameAllocationGlobalUBO[frameID])
	}

	// write new data to same buffer
	const sizeUBO = glm.SizeOfMat4 * 2
	const sizeSurface = glm.SizeOfVec2

	staging := make([]byte, 0, sizeUBO+sizeSurface)
	staging = append(staging, view.Data()...)        //  ubo mat4
	staging = append(staging, projection.Data()...)  //  ubo mat4
	staging = append(staging, surfaceSize.Data()...) // surf vec2

	// write data to memory
	allocUBO := m.heap.Write(
		staging,
		alloc.BufferTypeUniform,
		alloc.StorageTargetCoherent,
		alloc.FlagsNone,
	)

	m.frameAllocationGlobalUBO[frameID] = allocUBO

	// take info from alloc
	bufferInfos := []vulkan.DescriptorBufferInfo{
		// 0 = vert UBO (view + projection)
		{
			Buffer: allocUBO.Buffer,
			Offset: allocUBO.Offset,
			Range:  allocUBO.Size,
		},
		// 1 = frag UBO (surface size)
		{
			Buffer: allocUBO.Buffer, // todo: offset for second binding (but it panic)
			Offset: allocUBO.Offset, // todo: offset for second binding (but it panic)
			Range:  allocUBO.Size,   // todo: offset for second binding (but it panic)
		},
	}

	// prepare descriptor write
	descriptorSet := m.allocatedSets[frameID][layoutIndexGlobal]
	blueprint := m.blueprint.LayoutGlobal()
	writeSets := make([]vulkan.WriteDescriptorSet, 0, len(blueprint.bindings))

	for _, binding := range blueprint.bindings {
		writeSets = append(writeSets, vulkan.WriteDescriptorSet{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          descriptorSet,
			DstArrayElement: 0,
			DstBinding:      binding.Binding,
			DescriptorCount: binding.DescriptorCount,
			DescriptorType:  binding.DescriptorType,
			PBufferInfo: []vulkan.DescriptorBufferInfo{
				bufferInfos[binding.Binding],
			},
		})
	}

	// update references
	vulkan.UpdateDescriptorSets(m.ld.Ref(), uint32(len(writeSets)), writeSets, 0, nil)

	return Data{
		DescriptorSet: descriptorSet,
	}
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
