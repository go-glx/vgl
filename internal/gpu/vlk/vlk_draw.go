package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

type drawCallChunk struct {
	shader      *shader.Shader
	chunks      []alloc.Chunk
	indexBuffer alloc.Allocation
	polygonMode vulkan.PolygonMode
}

type drawCall struct {
	shader      *shader.Shader
	instances   []shader.InstanceData
	polygonMode vulkan.PolygonMode

	// todo: blending mode
	// todo: other drawing params
}

// drawQueue used for queue current shader draw with current options to queue
// all drawing actual will be executed in frameEnd
func (vlk *VLK) drawQueue(shader *shader.Shader, instance shader.InstanceData) {
	vlk.autoBake(shader, instance)
	vlk.currentBatch.instances = append(vlk.currentBatch.instances, instance)
}

func (vlk *VLK) autoBake(shader *shader.Shader, instance shader.InstanceData) {
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

	// brake: polygon mode changed
	if vlk.currentBatch.polygonMode != instance.PolygonMode() {
		brakeBaking = true
	}

	// brake: blend mode changed (todo)
	// brake: settings changed (todo)

	if brakeBaking {
		vlk.brakeBaking()
		vlk.currentBatch.shader = shader
		vlk.currentBatch.polygonMode = instance.PolygonMode()
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
			indexBuffer: vlk.indexBufferOf(drawCall.shader),
			polygonMode: drawCall.polygonMode,
		})
	}

	// move data to GPU
	buffs.Flush()

	// render
	countDrawCalls := 0
	vlk.cont.frameManager().FrameApplyCommands(func(_ uint32, cb vulkan.CommandBuffer) {
		for _, drawChunk := range drawChunks {
			sdr := drawChunk.shader
			indexCount := uint32(len(sdr.Meta().Indexes()))

			// bind pipe with current shader and options
			pipe := vlk.cont.pipelineFactory().NewPipeline(
				pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
					sdr.ModuleVert().Stage(),
					sdr.ModuleFrag().Stage(),
				}),
				pipeline.WithTopology(
					sdr.Meta().Topology(),
					sdr.Meta().TopologyRestartEnable(),
				),
				pipeline.WithVertexInput(
					sdr.Meta().Bindings(),
					sdr.Meta().Attributes(),
				),
				pipeline.WithRasterization(drawChunk.polygonMode),
				pipeline.WithColorBlend(),
				pipeline.WithMultisampling(),
			)

			vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, pipe)

			// index
			if drawChunk.indexBuffer.HasData {
				vulkan.CmdBindIndexBuffer(cb, drawChunk.indexBuffer.Buffer, 0, vulkan.IndexTypeUint16)
			}

			// draw instances
			for _, chunk := range drawChunk.chunks {
				// vertex buffer
				buffers := []vulkan.Buffer{chunk.Buffer}
				offsets := []vulkan.DeviceSize{vulkan.DeviceSize(chunk.BufferOffset)}
				vulkan.CmdBindVertexBuffers(cb, 0, uint32(len(buffers)), buffers, offsets)

				// Drawing Type: 1 (optimized indexed draw - instancing)
				if drawChunk.indexBuffer.HasData {
					instPerCall := min(chunk.InstancesCount, def.BufferIndexMapInstances)
					for firstInst := uint32(0); firstInst < chunk.InstancesCount; firstInst += instPerCall {
						// if we try to draw more instances, that fit in warm index map (>65536)
						// we split it into chunks of def.BufferIndexMapInstances size each
						vulkan.CmdDrawIndexed(cb, indexCount*instPerCall, 1, 0, int32(firstInst*sdr.Meta().VertexCount()), 0)
						countDrawCalls++
					}

					continue
				}

				// Drawing Type: 2 (default fallback)
				for i := uint32(0); i < chunk.InstancesCount; i++ {
					vulkan.CmdDraw(cb, indexCount, 1, i*indexCount, 0)
					countDrawCalls++
				}
			}
		}
	})

	// reset
	vlk.queue = make([]drawCall, 0, 32)
	vlk.currentBatch = &drawCall{}
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}

	return b
}
