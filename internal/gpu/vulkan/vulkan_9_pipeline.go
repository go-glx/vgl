package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"
)

func newPipelineManager(ld *vkLogicalDevice) *vkPipelineManager {
	return &vkPipelineManager{
		ld:              ld,
		shaderPipelines: map[shaderModuleID]vulkan.Pipeline{},
	}
}

func (pm *vkPipelineManager) free() {
	for _, pipeline := range pm.shaderPipelines {
		vulkan.DestroyPipeline(pm.ld.ref, pipeline, nil)
	}

	log.Printf("vk: freed: pipeline manager\n")
}

func (pm *vkPipelineManager) preloadPipelineFor(sp shaderProgram, pipeline vulkan.Pipeline) {
	if _, exist := pm.shaderPipelines[sp.ID()]; exist {
		panic(fmt.Errorf("failed preload shader pipeline '%s': already exist", sp.ID()))
	}

	log.Printf("vk: pipeline loaded for shader '%s'\n", sp.ID())
	pm.shaderPipelines[sp.ID()] = pipeline
}

func (pm *vkPipelineManager) pipeline(id shaderModuleID) vulkan.Pipeline {
	if pl, exist := pm.shaderPipelines[id]; exist {
		return pl
	}

	panic(fmt.Errorf("failed find pipeline for shader '%s'", id))
}
