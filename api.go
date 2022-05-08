package govgl

import (
	"github.com/fe3dback/govgl/arch"
	"github.com/fe3dback/govgl/config"
	"github.com/fe3dback/govgl/internal/gpu/vulkan"
)

type Render struct {
	renderer *vulkan.Vk
}

func NewRender(wm arch.WindowManager, cfg *config.Config) *Render {
	return &Render{
		renderer: vulkan.NewVulkanApi(wm, cfg),
	}
}

func (r *Render) Close() error {
	return r.renderer.Close()
}
