package dscptr

import (
	"bytes"
	"fmt"
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type (
	Manager struct {
		logger config.Logger
		ld     *logical.Device
		heap   *alloc.Heap
		pool   *Pool

		layouts                layoutsMap
		descriptorSets         descriptorSetsMap
		uniformBufferAlignSize uint32
		frameAllocations       frameAllocationsMap
	}

	layoutsMap        map[layoutIndex]vulkan.DescriptorSetLayout
	setsMap           map[layoutIndex]vulkan.DescriptorSet
	descriptorSetsMap map[frameID]setsMap

	layoutAllocationsMap map[layoutIndex][]alloc.Allocation
	frameAllocationsMap  map[frameID]layoutAllocationsMap
	DescriptorUpdates    map[bindingIndex][]byte
)

func NewManager(
	logger config.Logger,
	ld *logical.Device,
	pd *physical.Device,
	heap *alloc.Heap,
	pool *Pool,
) *Manager {
	layouts := make(layoutsMap)
	for index, bpLayout := range blueprint {
		layouts[index] = initializeLayout(ld.Ref(), bpLayout)
	}

	sets := initializeDescriptorSets(ld.Ref(), pool.Pool(), layouts)

	return &Manager{
		logger: logger,
		ld:     ld,
		heap:   heap,
		pool:   pool,

		layouts:                layouts,
		descriptorSets:         sets,
		uniformBufferAlignSize: uint32(pd.PrimaryGPU().Props.Limits.MinUniformBufferOffsetAlignment),
		frameAllocations:       make(frameAllocationsMap),
	}
}

func (m *Manager) Layouts() []vulkan.DescriptorSetLayout {
	layouts := make([]vulkan.DescriptorSetLayout, 0, len(m.layouts))

	for _, layout := range m.layouts {
		layouts = append(layouts, layout)
	}

	return layouts
}

func (m *Manager) UpdateSet(
	frameID frameID,
	index layoutIndex,
	updates DescriptorUpdates,
) vulkan.DescriptorSet {
	// clear before allocated data for this frame,layout
	m.freeMemory(frameID, index)

	// prepare set writes
	descriptorSet := m.descriptorSets[frameID][index]
	writeSets := make([]vulkan.WriteDescriptorSet, 0, len(updates))

	// group all updates by buffer type
	// example input : [ 1: [abc], 2: [ddd], 3: [eee] ]
	// example output: [ uniform: [1,3], storage: [2] ]
	for bufferType, bindingIndexes := range m.uniqueBindingTypes(index, updates) {
		// bindingUpdates is all byte data, that need to write into buffer
		bindingUpdates := make([][]byte, 0, len(bindingIndexes))
		for _, bindingID := range bindingIndexes {
			bindingUpdates = append(bindingUpdates, updates[bindingID])
		}

		// merge all data to single staging buffer
		// here:
		// - staging is united slice of bytes for all bindings
		// - sizes - size in bytes for each update (in same order)
		// - offsets - local offset in bytes for each update (in same order)
		staging, sizes, offsets := m.prepareStaging(bindingUpdates)

		// copy staging data to device
		allocation := m.heap.Write(staging, bufferType, alloc.StorageTargetCoherent, alloc.FlagsNone)
		m.writeToMemory(frameID, index, allocation)

		// add write set
		for writeInd, offset := range offsets {
			// writeInd is ptr to ordered (by bindingUpdates) data in staging slice
			// bindingIndexes ordered by bindingUpdates
			bindingID := bindingIndexes[writeInd]

			writeSets = append(writeSets, vulkan.WriteDescriptorSet{
				SType:           vulkan.StructureTypeWriteDescriptorSet,
				DstSet:          descriptorSet,
				DstBinding:      bindingID,
				DstArrayElement: 0,
				DescriptorCount: 1,
				DescriptorType:  blueprint[index].bindings[bindingID].descriptorType,
				PBufferInfo: []vulkan.DescriptorBufferInfo{
					{
						Buffer: allocation.Buffer,          // vulkan buffer ID
						Offset: allocation.Offset + offset, // global offset in vulkan buffer
						Range:  sizes[writeInd],            // local size for this binding data
					},
				},
			})
		}
	}

	// apply all changes from all write sets
	vulkan.UpdateDescriptorSets(m.ld.Ref(), uint32(len(writeSets)), writeSets, 0, nil)

	// return current descriptorSet for next binding it to pipeline
	return descriptorSet
}

func (m *Manager) uniqueBindingTypes(index layoutIndex, updates DescriptorUpdates) map[alloc.BufferType][]bindingIndex {
	results := make(map[alloc.BufferType][]bindingIndex)

	for bindingId := range updates {
		buffType := m.bufferTypeOfDescriptor(blueprint[index].bindings[bindingId].descriptorType)

		if _, exist := results[buffType]; !exist {
			results[buffType] = make([]bindingIndex, 0, len(updates))
		}

		results[buffType] = append(results[buffType], bindingId)
	}

	return results
}

func (m *Manager) bufferTypeOfDescriptor(dType vulkan.DescriptorType) alloc.BufferType {
	switch dType {
	case vulkan.DescriptorTypeUniformBuffer:
		return alloc.BufferTypeUniform
	case vulkan.DescriptorTypeStorageBuffer:
		return alloc.BufferTypeStorage
	default:
		panic(fmt.Errorf("unexpected descriptor type %d (%s). Possible need add new buffer type for it",
			dType,
			nameOfDescriptorType(dType),
		))
	}
}

// prepareStaging will write all block bytes data to single slice
// and align bytes in each block of zeroed space if needed
// function return result bytes slice and offset of each block start
func (m *Manager) prepareStaging(updates [][]byte) ([]byte, []vulkan.DeviceSize, []vulkan.DeviceSize) {
	totalSize := 0
	for _, update := range updates {
		totalSize += m.alignSize(len(update))
	}

	staging := make([]byte, 0, totalSize)
	sizes := make([]vulkan.DeviceSize, 0, len(updates))
	offsets := make([]vulkan.DeviceSize, 0, len(updates))

	for _, data := range updates {
		size := len(data)
		alignedSize := m.alignSize(size)
		uselessSize := alignedSize - size

		// write actual data
		sizes = append(sizes, vulkan.DeviceSize(size))
		offsets = append(offsets, vulkan.DeviceSize(len(staging))) // current staging len = offset of next block
		staging = append(staging, data...)

		// write trash data if needed, we need match device aligned size exactly
		if uselessSize > 0 {
			staging = append(staging, bytes.Repeat([]byte{0}, uselessSize)...)
		}
	}

	return staging, sizes, offsets
}

func (m *Manager) alignSize(realSize int) int {
	return int(math.Ceil(float64(realSize)/float64(m.uniformBufferAlignSize)) * float64(m.uniformBufferAlignSize))
}

func (m *Manager) freeMemory(frameID frameID, index layoutIndex) {
	const defaultAllocsCapacity = 16

	// get layouts map, or create if not exist
	layoutAllocsMap, exist := m.frameAllocations[frameID]
	if !exist {
		m.frameAllocations[frameID] = make(layoutAllocationsMap)
		layoutAllocsMap = m.frameAllocations[frameID]
	}

	// get allocs map, or create if not exist
	allocs, exist := layoutAllocsMap[index]
	if !exist {
		m.frameAllocations[frameID][index] = make([]alloc.Allocation, 0, defaultAllocsCapacity)
		return
	}

	// free all exist allocs
	for _, allocation := range allocs {
		m.heap.Free(allocation)
	}

	// clear allocs buffer
	m.frameAllocations[frameID][index] = make([]alloc.Allocation, 0, defaultAllocsCapacity)
}

func (m *Manager) writeToMemory(frameID frameID, index layoutIndex, alloc alloc.Allocation) {
	m.frameAllocations[frameID][index] = append(m.frameAllocations[frameID][index], alloc)
}
