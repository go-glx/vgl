package vgl

import (
	"github.com/go-glx/glx"
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
func (r *Render) toLocalSpace2d(pos glx.Vec2) glx.Vec2 {
	w, h := r.api.GetSurfaceSize()

	return glx.Vec2{
		X: 2*(pos.X/w) - 1,
		Y: 2*(pos.Y/h) - 1,
	}
}

func (r *Render) toLocalAspectRation(n float32) float32 {
	w, h := r.api.GetSurfaceSize()

	if w <= 0 || h <= 0 {
		return 0
	}

	if w > h {
		return n / w
	}

	return n / h
}

func (r *Render) toLocalAspectRationX(n float32) float32 {
	w, _ := r.api.GetSurfaceSize()
	return n / w
}

func (r *Render) toLocalAspectRationY(n float32) float32 {
	_, h := r.api.GetSurfaceSize()
	return n / h
}

func (r *Render) cullingPoint(vert glx.Vec2) bool {
	return vert.X >= -1 && vert.X <= 1 && vert.Y >= -1 && vert.Y <= 1
}

func (r *Render) cullingLine(vert [2]glx.Vec2) bool {
	// todo
	return true
}

func (r *Render) cullingTriangle(vert [3]glx.Vec2) bool {
	// todo
	return true
}

func (r *Render) cullingRect(vert [4]glx.Vec2) bool {
	// todo
	return true
}
