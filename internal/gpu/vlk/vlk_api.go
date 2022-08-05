package vlk

import (
	"fmt"
	"math"
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
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

	// todo: remove test draw
	triangle := vlk.cont.shaderManager().ShaderByID(buildInShaderTriangle)

	const step = 0.1
	const stepHalf = step / 2

	ms := float32(math.Sin(float64(time.Now().UnixMilli()) * 0.005))
	fmt.Println(ms)

	for xx := float32(-1.0); xx < 1; xx += step {
		for yy := float32(-1.0); yy < 1; yy += step {
			xLeft := xx
			xCenter := xx + (stepHalf + (stepHalf * ms))
			xRight := xx + step

			vlk.drawQueue(triangle, &dataTriangle{
				vertexes: [3]glm.Vec2{
					{X: xCenter, Y: yy},
					{X: xRight, Y: yy + step},
					{X: xLeft, Y: yy + step},
				},
				colors: [3]glm.Vec3{
					{R: 1.0, G: 0.0, B: 0.0},
					{R: 0.0, G: 1.0, B: 0.0},
					{R: 0.0, G: 0.0, B: 1.0},
				},
			})
		}
	}

	// todo: ^^^^^^^^^^^^^^

	vlk.drawAll()
	vlk.cont.frameManager().FrameEnd()
}

func (vlk *VLK) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	if !vlk.isReady {
		return
	}

	// todo:
	// vlk.cont.frameManager().frameApplyCommands(func(cb vulkan.CommandBuffer) {
	// 	// todo: add command to draw rect
	// 	// todo: in current command buffer
	// })
}
