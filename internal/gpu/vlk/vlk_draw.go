package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

// drawQueue used for queue current shader draw with current options to queue
// all drawing actual will be executed in frameEnd
func (vlk *VLK) drawQueue(shader *shader.Shader, instance shaderData) {
	brakeBaking := false

	// brake: shader changed
	if vlk.currentBatch.shader != nil {
		if vlk.currentBatch.shader.Meta().ID() != shader.Meta().ID() {
			brakeBaking = true
		}
	} else {
		// prev shader is nil, so it changed to new
		brakeBaking = true
	}

	// brake: blend mode changed (todo)
	// brake: settings changed (todo)

	if brakeBaking {
		vlk.brakeBaking()
		vlk.currentBatch.shader = shader
	}

	vlk.currentBatch.instances = append(vlk.currentBatch.instances, instance)
}

// brakeBaking should be called when render params or shader is changes
func (vlk *VLK) brakeBaking() {
	if len(vlk.currentBatch.instances) <= 0 {
		return
	}

	vlk.queue = append(vlk.queue, *vlk.currentBatch)
	vlk.currentBatch = &drawCall{}
}

// brakeBaking should be called when render params or shader is changes
func (vlk *VLK) drawAll() {
	if len(vlk.queue) == 0 {
		return
	}

	// todo: memory manager with chunks

	// prepare buffer
	bufferIndex := 0
	bufferOffset := 0

	// write data to vertex buffer
	for _, drawCall := range vlk.queue {
		// todo: move it to memory manager service

		drawCall.bufferIndex = bufferIndex
		drawCall.bufferOffset = bufferOffset

		// todo: vertex buffer
		// todo: buffer flush
		// todo: apply global buffer offset to drawCall
		_ = drawCall
	}

	// render
	countDrawCalls := 0
	vlk.cont.frameManager().FrameApplyCommands(func(_ uint32, cb vulkan.CommandBuffer) {
		for _, drawCall := range vlk.queue {
			// bind pipe with current shader and options
			pipe := vlk.cont.pipelineFactory().NewPipeline(
				pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
					drawCall.shader.ModuleVert().Stage(),
					drawCall.shader.ModuleFrag().Stage(),
				}),
				pipeline.WithTopology(drawCall.shader.Meta().Topology()),
				pipeline.WithVertexInput(
					drawCall.shader.Meta().Bindings(),
					drawCall.shader.Meta().Attributes(),
				),
				pipeline.WithRasterization(vulkan.PolygonModeFill),
				pipeline.WithColorBlend(),
				pipeline.WithMultisampling(),
			)

			vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, pipe)

			// draw instances
			indexCount := uint32(len(drawCall.instances[0].Indexes()))
			instancesCount := uint32(len(drawCall.instances))

			// todo: options for it
			vulkan.CmdBindVertexBuffers(cb, 0, uint32(1), []vulkan.Buffer{vertexBuffer}, []vulkan.DeviceSize{0})
			vulkan.CmdDrawIndexed(cb, indexCount, instancesCount, 0, 0, 0)
			countDrawCalls++
		}
	})

	// reset
	vlk.queue = []drawCall{}
	vlk.currentBatch = &drawCall{}
}
