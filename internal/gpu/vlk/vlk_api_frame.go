package vlk

import (
	"github.com/vulkan-go/vulkan"
)

// WarmUp will warm vlk renderer and create all needed
// objects for work, this must be called one time
// before first FrameStart
func (vlk *VLK) WarmUp() {
	// request some managers, this will create it
	// and all dependencies, like swapChain, renderPass, etc..
	_ = vlk.cont.frameManager()
	_ = vlk.cont.shaderManager()
}

func (vlk *VLK) GPUWait() {
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())
}

// GetSurfaceSize returns current surface size in pixels
// for default surface (screen) is window width and height
// this very fast function, will always return cached values
// and can be called thousands times per frame
func (vlk *VLK) GetSurfaceSize() (w uint32, h uint32) {
	return vlk.surfacesSize[vlk.surfaceInd][0],
		vlk.surfacesSize[vlk.surfaceInd][1]
}

func (vlk *VLK) FrameStart() {
	if !vlk.isReady {
		return
	}

	vlk.cont.frameManager().FrameBegin()
}

func (vlk *VLK) FrameEnd() {
	if !vlk.isReady {
		return
	}

	vlk.drawAll()
	vlk.cont.frameManager().FrameEnd()
}
