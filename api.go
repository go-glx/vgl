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

	// init renderer resources and prepare GPU to work
	renderer.WarmUp()
	renderer.GPUWait()

	api := &Render{
		closer: closer,
		api:    renderer,
	}

	registerStdShaders(api)
	return api
}

// Close SHOULD BE called on application exit
// this will free all vulkan GPU resources, release
// memory, etc..
func (r *Render) Close() error {
	r.api.GPUWait()
	r.closer.close()
	return nil
}

func registerStdShaders(api *Render) {
	for _, buildInShader := range stdShaders {
		api.RegisterShader(&buildInShader)
	}
}
