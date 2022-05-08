package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"
)

func newShaderManager(ld *vkLogicalDevice) *vkShaderManager {
	return &vkShaderManager{
		ld:      ld,
		modules: make(map[shaderModuleID]*vkShaderModule),
	}
}

func (sm *vkShaderManager) free() {
	for _, shaderModule := range sm.modules {
		vulkan.DestroyShaderModule(sm.ld.ref, shaderModule.module, nil)
	}

	log.Printf("vk: freed: shader manager\n")
}

func (sm *vkShaderManager) preloadShader(p shaderProgram) {
	sm.modules[p.ID()+shaderTypeFrag] = sm.allocModule(p.ID()+shaderTypeFrag, p.ProgramFrag(), vulkan.ShaderStageFragmentBit)
	sm.modules[p.ID()+shaderTypeVert] = sm.allocModule(p.ID()+shaderTypeVert, p.ProgramVert(), vulkan.ShaderStageVertexBit)
}

func (sm *vkShaderManager) shaderModule(id shaderModuleID) *vkShaderModule {
	if module, exist := sm.modules[id]; exist {
		return module
	}

	return sm.modules[id]
}

func (sm *vkShaderManager) allocModule(id shaderModuleID, byteCode []byte, stageType vulkan.ShaderStageFlagBits) *vkShaderModule {
	createInfo := &vulkan.ShaderModuleCreateInfo{
		SType:    vulkan.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(byteCode)),
		PCode:    vkTransformBytes(byteCode),
	}

	var shaderModule vulkan.ShaderModule
	vkAssert(
		vulkan.CreateShaderModule(sm.ld.ref, createInfo, nil, &shaderModule),
		fmt.Errorf("failed create shaderPipelines module from '%s' shaderPipelines", id),
	)

	log.Printf("vk: shaderPipelines allocated '%s', len=%d\n", id, len(byteCode))
	return &vkShaderModule{
		module: shaderModule,
		stageInfo: vulkan.PipelineShaderStageCreateInfo{
			SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
			Stage:  stageType,
			Module: shaderModule,
			PName:  fmt.Sprintf("%s\x00", shaderEntryPoint),
		},
	}
}
