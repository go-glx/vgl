package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type drawCallChunk struct {
	shader      *shader.Shader
	chunks      []alloc.Chunk
	indexBuffer vulkan.Buffer
}

// drawQueue used for queue current shader draw with current options to queue
// all drawing actual will be executed in frameEnd
func (vlk *VLK) drawQueue(shader *shader.Shader, instance shader.InstanceData) {
	vlk.autoBake(shader)
	vlk.currentBatch.instances = append(vlk.currentBatch.instances, instance)
}

func (vlk *VLK) autoBake(shader *shader.Shader) {
	brakeBaking := false

	if vlk.currentBatch == nil {
		vlk.currentBatch = &drawCall{}
	}

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
	vlk.brakeBaking()
	if len(vlk.queue) == 0 {
		return
	}

	// prepare buffer
	buffs := vlk.cont.allocBuffers()
	buffs.Clear()

	// stage all shader data to buffers
	drawChunks := make([]drawCallChunk, 0)

	for _, drawCall := range vlk.queue {
		if len(drawCall.instances) == 0 {
			continue
		}

		drawChunks = append(drawChunks, drawCallChunk{
			shader:      drawCall.shader,
			chunks:      buffs.Write(drawCall.instances),
			indexBuffer: buffs.IndexBufferOf(drawCall.shader, drawCall.instances[0].Indexes()),
		})
	}

	// move data to GPU
	buffs.Flush()

	// render
	countDrawCalls := 0
	vlk.cont.frameManager().FrameApplyCommands(func(_ uint32, cb vulkan.CommandBuffer) {
		for _, drawChunk := range drawChunks {
			sdr := drawChunk.shader

			// bind pipe with current shader and options
			pipe := vlk.cont.pipelineFactory().NewPipeline(
				pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
					sdr.ModuleVert().Stage(),
					sdr.ModuleFrag().Stage(),
				}),
				pipeline.WithTopology(sdr.Meta().Topology()),
				pipeline.WithVertexInput(
					sdr.Meta().Bindings(),
					sdr.Meta().Attributes(),
				),
				pipeline.WithRasterization(vulkan.PolygonModeFill),
				pipeline.WithColorBlend(),
				pipeline.WithMultisampling(),
			)

			vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, pipe)

			// index
			vulkan.CmdBindIndexBuffer(cb, drawChunk.indexBuffer, 0, vulkan.IndexTypeUint16)

			// draw instances
			for _, chunk := range drawChunk.chunks {
				// vertex buffer
				buffers := []vulkan.Buffer{chunk.Buffer}
				offsets := []vulkan.DeviceSize{vulkan.DeviceSize(chunk.BufferOffset)}
				vulkan.CmdBindVertexBuffers(cb, 0, uint32(len(buffers)), buffers, offsets)

				// draw instances
				for i := uint32(0); i < chunk.InstancesCount; i++ {
					vulkan.CmdDraw(cb, chunk.IndexCount, 1, i*chunk.IndexCount, 0)
					countDrawCalls++
				}

				// todo: optimization:
				// todo: indexed draw all data in one draw call
				// vulkan.CmdDrawIndexed(cb, chunk.IndexCount, chunk.InstancesCount, 0, 0, 0)
			}
		}
	})

	// reset
	vlk.queue = []drawCall{}
	vlk.currentBatch = &drawCall{}
}
