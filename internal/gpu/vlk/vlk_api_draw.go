package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type (
	DrawOptions struct {
		PolygonMode vulkan.PolygonMode
		// BlendMode vulkan.BlendOp // todo
	}
)

func (vlk *VLK) Draw(name string, opts DrawOptions, data shader.InstanceData) {
	if !vlk.isReady {
		return
	}

	vlk.drawQueue(vlk.cont.shaderManager().ShaderByID(name), opts, data)
}
