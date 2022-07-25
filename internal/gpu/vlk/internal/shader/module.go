package shader

import "github.com/vulkan-go/vulkan"

type Module struct {
	module    vulkan.ShaderModule
	stageInfo vulkan.PipelineShaderStageCreateInfo
}

func (m *Module) Module() vulkan.ShaderModule {
	return m.module
}

func (m *Module) Stage() *vulkan.PipelineShaderStageCreateInfo {
	return &m.stageInfo
}
