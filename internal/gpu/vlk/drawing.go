package vlk

import "github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"

type drawCall struct {
	shader    *shader.Shader
	instances []shader.InstanceData

	// todo: blending mode
	// todo: other drawing params
}
