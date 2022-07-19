package vgl

import "github.com/go-glx/vgl/glm"

// All APIs work with vector vertexes, where:
// x = -1 .. 1
// y = -1 .. 1
// For example:
//  - Top-Left    : {x=-1, y=-1}
//  - Bottom-Right: {x= 1, y= 1}
//  - Center      : {x= 0, y= 0}

// todo: design primitives API

func (r *Render) Draw2DRectExt(
	vertexPos [4]glm.Vec2,
	vertexColor [4]glm.Vec3,
	outline bool,
) {
	// todo: outline
	r.gpuApi.DrawRect(vertexPos, vertexColor)
}
