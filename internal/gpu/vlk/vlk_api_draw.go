package vlk

import (
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

func (vlk *VLK) Draw(name string, data shader.InstanceData) {
	if !vlk.isReady {
		return
	}

	vlk.drawQueue(vlk.cont.shaderManager().ShaderByID(name), data)
}
