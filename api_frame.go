package vgl

import "github.com/go-glx/vgl/glm"

// FrameStart should be called before any drawing in current frame
func (r *Render) FrameStart() {
	r.api.FrameStart()
}

// FrameEnd should be called after any drawing in current frame
// this function will draw all queued objects in GPU
// and swap image buffer from GPU to screen
func (r *Render) FrameEnd() {
	r.api.FrameEnd()
}

// ListenStats allows to subscribe to render frame stats
// this function will execute custom callback function with
// last frame stats
func (r *Render) ListenStats(listener func(glm.Stats)) {
	r.api.ListenStats(listener)
}
