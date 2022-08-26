package dscptr

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func initializeLayout(ld vulkan.Device, bpLayout blueprintLayout) vulkan.DescriptorSetLayout {
	bindings := make([]vulkan.DescriptorSetLayoutBinding, 0, len(bpLayout.bindings))
	for bindingIndex, bpBinding := range bpLayout.bindings {
		bindings = append(bindings, vulkan.DescriptorSetLayoutBinding{
			Binding:         bindingIndex,
			DescriptorType:  bpBinding.descriptorType,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageFlags(bpBinding.flags),
		})
	}

	info := &vulkan.DescriptorSetLayoutCreateInfo{
		SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bpLayout.bindings)),
		PBindings:    bindings,
	}

	var layout vulkan.DescriptorSetLayout
	must.Work(vulkan.CreateDescriptorSetLayout(ld, info, nil, &layout))

	return layout
}

func initializeDescriptorSets(ld vulkan.Device, pool vulkan.DescriptorPool, layouts layoutsMap) descriptorSetsMap {
	descriptorSets := make(descriptorSetsMap)

	for frameID := frameID(0); frameID < framesCount; frameID++ {
		set := make(setsMap)

		for layoutIndex, layout := range layouts {
			set[layoutIndex] = allocateSet(ld, pool, layout)
		}

		descriptorSets[frameID] = set
	}

	return descriptorSets
}

func allocateSet(ld vulkan.Device, pool vulkan.DescriptorPool, layout vulkan.DescriptorSetLayout) vulkan.DescriptorSet {
	setAllocateInfo := vulkan.DescriptorSetAllocateInfo{
		SType:              vulkan.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     pool,
		DescriptorSetCount: 1,
		PSetLayouts:        []vulkan.DescriptorSetLayout{layout},
	}

	var set vulkan.DescriptorSet
	must.Work(vulkan.AllocateDescriptorSets(ld, &setAllocateInfo, &set))

	return set
}
