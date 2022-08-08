package vlk

import (
	"github.com/go-glx/vgl/glm"
)

func (vlk *VLK) DrawPoint(vertex glm.Vec2, color glm.Vec4) {
	if !vlk.isReady {
		return
	}

	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderPoint),
		&dataPoint{
			vertex: vertex,
			color:  color,
		},
	)
}

func (vlk *VLK) DrawTriangle(vertexes [3]glm.Vec2, colors [3]glm.Vec4, filled bool) {
	if !vlk.isReady {
		return
	}

	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderTriangle),
		&dataTriangle{
			vertexes: vertexes,
			colors:   colors,
			filled:   filled,
		},
	)
}
