package vgl

// SurfaceSize will return size of current surface
// when current surface is 0 (default), it will return window size
// this is fast function, you can call it many times per frame
func (r *Render) SurfaceSize() (width float32, height float32) {
	return r.api.GetSurfaceSize()
}
