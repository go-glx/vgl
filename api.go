package govgl

import (
	"github.com/fe3dback/govgl/arch"
	"github.com/fe3dback/govgl/config"
	"github.com/fe3dback/govgl/internal/gpu/vulkan"
)

type Render struct {
	gpuApi *vulkan.Vk
}

func NewRender(wm arch.WindowManager, cfg *config.Config) *Render {
	return &Render{
		gpuApi: vulkan.NewVulkanApi(wm, cfg),
	}
}

// WaitGPU should be called in graceful engine shutdown
// before application exit. This command will sleep and wait
// current io operation done in GPU.
// SHOULD BE called before Close
func (r *Render) WaitGPU() {
	r.gpuApi.GPUWait()
}

// Close SHOULD BE called on application exit
// this will free all vulkan GPU resources, release
// memory, etc..
// Render.WaitGPU SHOULD BE called right before Close
func (r *Render) Close() error {
	return r.gpuApi.Close()
}
