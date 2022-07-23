package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

// WarmUp will warm vlk renderer and create all needed
// objects for work, this must be called one time
// before first FrameStart
func (vlk *VLK) WarmUp() {
	// request frameManager, this will create it
	// and all dependencies, like swapChain, renderPass, etc..
	_ = vlk.cont.frameManager()
}

func (vlk *VLK) GPUWait() {
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())
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

	vlk.cont.frameManager().FrameEnd()
}

func (vlk *VLK) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	if !vlk.isReady {
		return
	}

	vlk.cont.frameManager().FrameApplyCommands(func(cb vulkan.CommandBuffer) {
		// todo: add command to draw rect
		// todo: in current command buffer
	})
}
