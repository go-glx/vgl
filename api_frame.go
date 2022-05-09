package govgl

// FrameStart should be called before any drawing in current frame
func (r *Render) FrameStart() {
	r.gpuApi.FrameStart()
}

// FrameEnd should be called after any drawing in current frame
// this function will draw all queued objects in GPU
// and swap image buffer from GPU to screen
func (r *Render) FrameEnd() {
	r.gpuApi.Draw()
	r.gpuApi.FrameEnd()
}
