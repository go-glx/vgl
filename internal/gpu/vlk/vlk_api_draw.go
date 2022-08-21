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

func (vlk *VLK) DrawLine(vertexes [2]glm.Vec2, colors [2]glm.Vec4) {
	if !vlk.isReady {
		return
	}

	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderLine),
		&dataLine{
			vertexes: vertexes,
			colors:   colors,
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

func (vlk *VLK) DrawCircle(vertexes [4]glm.Vec2, colors [4]glm.Vec4, thickness glm.Vec1, smooth glm.Vec1) {
	if !vlk.isReady {
		return
	}

	centerX := (vertexes[0].X + vertexes[1].X + vertexes[2].X + vertexes[3].X) / 4
	centerY := (vertexes[0].Y + vertexes[1].Y + vertexes[2].Y + vertexes[3].Y) / 4

	center := glm.Vec2{
		X: centerX,
		Y: centerY,
	}

	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderCircle),
		&dataCircle{
			[5]glm.Vec2{
				vertexes[0],
				vertexes[1],
				vertexes[2],
				vertexes[3],
				center,
			},
		},
	)

	// todo: remove:
	vlk.DrawRect(vertexes, colors, false)
	vlk.DrawPoint(center, colors[0])
}

// DrawRect input vertexes order is [tl,tr,br,bl]
func (vlk *VLK) DrawRect(vertexes [4]glm.Vec2, colors [4]glm.Vec4, filled bool) {
	if !vlk.isReady {
		return
	}

	if !filled {
		vlk.drawQueue(
			vlk.cont.shaderManager().ShaderByID(buildInShaderRect),
			&dataRectOutline{
				vertexes: vertexes,
				colors:   colors,
			},
		)

		return
	}

	// drawing two triangles faster, that outlined rect
	// with custom polygon mode
	// when rect is filled, two triangles visually looks same as rect
	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderTriangle),
		&dataTriangle{
			vertexes: [3]glm.Vec2{
				vertexes[2], // br
				vertexes[3], // bl
				vertexes[0], // tl
			},
			colors: [3]glm.Vec4{
				colors[2],
				colors[3],
				colors[0],
			},
			filled: true,
		},
	)
	vlk.drawQueue(
		vlk.cont.shaderManager().ShaderByID(buildInShaderTriangle),
		&dataTriangle{
			vertexes: [3]glm.Vec2{
				vertexes[2], // br
				vertexes[0], // tl
				vertexes[1], // tr
			},
			colors: [3]glm.Vec4{
				colors[2],
				colors[0],
				colors[1],
			},
			filled: true,
		},
	)
}
