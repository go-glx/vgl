package vgl

import (
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk"
)

type Render struct {
	closer *Closer
	api    *vlk.VLK
}

func NewRender(wm arch.WindowManager, cfg *config.Config) *Render {
	closer := newCloser()
	container := vlk.NewContainer(closer, wm, cfg)
	renderer := container.VulkanRenderer()
	renderer.GPUWait()

	return &Render{
		closer: closer,
		api:    renderer,
	}
}

// WaitGPU should be called in graceful engine shutdown
// before application exit. This command will sleep and wait
// current io operation done in GPU.
// SHOULD BE called before Close
func (r *Render) WaitGPU() {
	r.api.GPUWait()
}

// Close SHOULD BE called on application exit
// this will free all vulkan GPU resources, release
// memory, etc..
// Render.WaitGPU SHOULD BE called right before Close
func (r *Render) Close() error {
	r.closer.close()
	return nil
}
