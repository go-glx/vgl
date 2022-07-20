package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

func (vlk *VLK) GPUWait() {
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())
}

func (vlk *VLK) FrameStart() {
	// todo
}

func (vlk *VLK) FrameEnd() {
	// todo
}

func (vlk *VLK) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	// todo
}
