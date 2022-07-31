package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
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

	vlk.cont.frameManager().FrameApplyCommands(func(_ uint32, cb vulkan.CommandBuffer) {
		triangle := vlk.cont.shaderManager().ShaderByID(buildInShaderTriangle)

		pipe := vlk.cont.pipelineFactory().NewPipeline(
			pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
				*triangle.ModuleVert().Stage(),
				*triangle.ModuleFrag().Stage(),
			}),
			pipeline.WithTopology(triangle.Meta().Topology()),
			pipeline.WithVertexInput(
				triangle.Meta().Bindings(),
				triangle.Meta().Attributes(),
			),
			pipeline.WithRasterization(vulkan.PolygonModeFill),
			pipeline.WithColorBlend(),
			pipeline.WithMultisampling(),
		)

		vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, pipe)

		// todo: 3,1 to shader
		for i := 0; i < 1024; i++ {
			vulkan.CmdDraw(cb, 3, 1, 0, 0)
		}
	})
	// todo: ^^^^^^^^^^^^^^

	vlk.cont.frameManager().FrameEnd()
}

func (vlk *VLK) DrawRect(vertexPos [4]glm.Vec2, vertexColor [4]glm.Vec3) {
	if !vlk.isReady {
		return
	}

	// todo:
	// vlk.cont.frameManager().FrameApplyCommands(func(cb vulkan.CommandBuffer) {
	// 	// todo: add command to draw rect
	// 	// todo: in current command buffer
	// })
}
