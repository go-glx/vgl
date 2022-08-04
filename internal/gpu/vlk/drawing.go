package vlk

import "github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"

type drawCall struct {
	shader       *shader.Shader
	instances    []shaderData
	bufferIndex  uint32
	bufferOffset uint32

	// todo: blending mode
	// todo: other drawing params
}
