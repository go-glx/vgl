package shader

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type Manager struct {
	logger  config.Logger
	shaders map[string]*Shader

	ld *logical.Device
}

func NewManager(logger config.Logger, ld *logical.Device) *Manager {
	return &Manager{
		logger:  logger,
		shaders: make(map[string]*Shader),

		ld: ld,
	}
}

func (m *Manager) Free() {
	for _, shader := range m.shaders {
		vulkan.DestroyShaderModule(m.ld.Ref(), shader.moduleVert.module, nil)
		vulkan.DestroyShaderModule(m.ld.Ref(), shader.moduleFrag.module, nil)
	}

	m.logger.Debug("freed: shaders")
}

func (m *Manager) ShaderByID(id string) *Shader {
	if shader, exist := m.shaders[id]; exist {
		return shader
	}

	panic(fmt.Errorf("shader '%s' not registered in manager and cannot be executed", id))
}

func (m *Manager) RegisterShader(meta *Meta) {
	m.shaders[meta.id] = m.createCompiledShader(meta)
}

func (m *Manager) createCompiledShader(meta *Meta) *Shader {
	return &Shader{
		meta:       meta,
		moduleVert: m.createModule(meta.id, meta.vert, TypeVertexBit),
		moduleFrag: m.createModule(meta.id, meta.frag, TypeFragmentBit),
	}
}

func (m *Manager) createModule(id string, byteCode []byte, shaderType Type) *Module {
	info := &vulkan.ShaderModuleCreateInfo{
		SType:    vulkan.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(byteCode)),
		PCode:    vkconv.TransformByteCode(byteCode),
	}

	var shaderModule vulkan.ShaderModule
	must.Work(vulkan.CreateShaderModule(m.ld.Ref(), info, nil, &shaderModule))

	m.logger.Debug(fmt.Sprintf("created shader '%s' of type '%s', len=%d", id, shaderType, len(byteCode)))
	return &Module{
		module: shaderModule,
		stageInfo: vulkan.PipelineShaderStageCreateInfo{
			SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
			Stage:  shaderType.VulkanShaderStage(),
			Module: shaderModule,
			PName:  fmt.Sprintf("%s\x00", def.ShaderEntryPoint),
		},
	}
}
