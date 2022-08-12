package vgl

import (
	"github.com/go-glx/vgl/glm"
)

// transform pos in pixel screen space
// to local renderer space [ -1 .. 0 .. +1 ]
// local space is
//  __________________________ -> Width
//  |          Y:-1          |
//  |                        |
//  | X:-1     0, 0     x:+1 |
//  |                        |
//  |          Y:+1          |
//  --------------------------
//                           ^
//                        Height
//
// This is fast function, can be called thousands times
// per frame, surface w/h is cached in renderer
func (r *Render) toLocalSpace2d(pos glm.Local2D) glm.Vec2 {
	w, h := r.api.GetSurfaceSize()

	return glm.Vec2{
		X: 2*(float32(pos[0])/float32(w)) - 1,
		Y: 2*(float32(pos[1])/float32(h)) - 1,
	}
}

func (r *Render) toLocalAspectRation(n int) float32 {
	w, h := r.api.GetSurfaceSize()

	if w <= 0 || h <= 0 {
		return 0
	}

	if w > h {
		return float32(n) / float32(h)
	}

	return float32(n) / float32(w)
}

func (r *Render) cullingPoint(vert glm.Vec2) bool {
	return vert.X >= -1 && vert.X <= 1 && vert.Y >= -1 && vert.Y <= 1
}

func (r *Render) cullingLine(vert [2]glm.Vec2) bool {
	// todo
	return true
}

func (r *Render) cullingTriangle(vert [3]glm.Vec2) bool {
	// todo
	return true
}

func (r *Render) cullingCircle(vert glm.Vec2, radius float32) bool {
	// todo
	return true
}

func (r *Render) cullingRect(vert [4]glm.Vec2) bool {
	// todo
	return true
}
