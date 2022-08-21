package vlk

import (
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/descriptors"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

// default capacity for batch queue
// will be reset every frame
const queueCapacity = 32

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
	vlk.stats.FrameIndex++

	vlk.brakeBaking()
	if len(vlk.queue) == 0 {
		return
	}

	// render
	vlk.stats.DrawCalls = 0
	vlk.cont.frameManager().FrameApplyCommands(func(imageID uint32, cb vulkan.CommandBuffer) {
		// write global data to uniform buffers that need for every shader
		view := glm.Mat4Identity()       // todo
		projection := glm.Mat4Identity() // todo
		surfaceSize := glm.Vec2{
			X: float32(vlk.surfacesSize[vlk.surfaceInd][0]),
			Y: float32(vlk.surfacesSize[vlk.surfaceInd][1]),
		}

		globalUbo := vlk.cont.descriptorsManager().UpdateGlobalUBO(uint8(imageID), view, projection, surfaceSize)

		// split drawing queue to chunks
		// and upload all required data for rendering into buffers
		ts := time.Now()
		drawChunks := vlk.prepareDrawingChunks(imageID)
		vlk.stats.TimeFlushVertexBuffer = time.Since(ts)
		vlk.stats.DrawChunks = len(drawChunks)

		for _, drawChunk := range drawChunks {
			vlk.stats.DrawCalls += vlk.drawChunk(cb, globalUbo, drawChunk)
		}
	})

	// reset
	vlk.queue = make([]drawCall, 0, queueCapacity)
	vlk.currentBatch = &drawCall{}
}

func (vlk *VLK) prepareDrawingChunks(imageID uint32) []drawCallChunk {
	buff := vlk.cont.allocBuffers()

	// clear all previously allocated vertex buffers owned by this frame
	buff.ClearVertexBuffersOwnedBy(imageID)

	// stage all shader data to buffers
	drawChunks := make([]drawCallChunk, 0, len(vlk.queue))
	uniqShaders := make(map[string]struct{}, 16)

	for _, drawCall := range vlk.queue {
		if len(drawCall.instances) == 0 {
			continue
		}

		uniqShaders[drawCall.shader.Meta().ID()] = struct{}{}
		drawChunks = append(drawChunks, drawCallChunk{
			shader:      drawCall.shader,
			chunks:      buff.WriteVertexBuffersFromInstances(imageID, drawCall.instances),
			indexBuffer: vlk.indexBufferOf(drawCall.shader),
			polygonMode: drawCall.polygonMode,
		})
	}

	vlk.stats.DrawUniqueShaders = len(uniqShaders)

	return drawChunks
}

func (vlk *VLK) drawChunk(cb vulkan.CommandBuffer, globalUBO descriptors.Data, drawChunk drawCallChunk) int {
	// bind pipe with current shader and options
	ts := time.Now()
	pipeInfo := vlk.createPipeline(drawChunk)
	vlk.stats.TimeCreatePipeline = time.Since(ts)
	vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, pipeInfo.Pipeline)

	descriptorSets := []vulkan.DescriptorSet{
		globalUBO.DescriptorSet,
	}

	// todo: not rebind sets, if [pipelineLayout and descriptorSets] not changed from previous call
	vulkan.CmdBindDescriptorSets(
		cb, vulkan.PipelineBindPointGraphics,
		pipeInfo.Layout, 0,
		uint32(len(descriptorSets)), descriptorSets,
		0, nil,
	)

	// bind index buffer
	if drawChunk.indexBuffer.Valid {
		vulkan.CmdBindIndexBuffer(cb, drawChunk.indexBuffer.Buffer, drawChunk.indexBuffer.Offset, vulkan.IndexTypeUint16)
	}

	// params
	indexCount := uint32(len(drawChunk.shader.Meta().Indexes()))
	vertexCount := drawChunk.shader.Meta().VertexCount()
	countDrawCalls := 0

	// draw instances
	for _, chunk := range drawChunk.chunks {
		// bind vertex buffer
		buffers := []vulkan.Buffer{chunk.Buffer}
		offsets := []vulkan.DeviceSize{chunk.BufferOffset}
		vulkan.CmdBindVertexBuffers(cb, 0, uint32(len(buffers)), buffers, offsets)

		// Drawing Type: 1 (optimized indexed draw - instancing)
		if drawChunk.indexBuffer.Valid {
			instPerCall := min(chunk.InstancesCount, def.BufferIndexMaxInstances)
			for firstInst := uint32(0); firstInst < chunk.InstancesCount; firstInst += instPerCall {
				// if we try to draw more instances, that fit in warm index cache (>65536)
				// we split it into chunks of def.BufferIndexMaxInstances size each

				ts := time.Now()
				vulkan.CmdDrawIndexed(cb, indexCount*instPerCall, 1, 0, int32(firstInst*vertexCount), 0)
				vlk.stats.TimeRenderInstanced += time.Since(ts)

				countDrawCalls++
			}

			continue
		}

		// Drawing Type: 2 (default fallback)
		for i := uint32(0); i < chunk.InstancesCount; i++ {
			ts := time.Now()
			vulkan.CmdDraw(cb, indexCount, 1, i*indexCount, 0)
			vlk.stats.TimeRenderFallback += time.Since(ts)

			countDrawCalls++
		}
	}

	return countDrawCalls
}

func (vlk *VLK) createPipeline(call drawCallChunk) pipeline.Info {
	return vlk.cont.pipelineFactory().NewPipeline(
		pipeline.WithLayout(pipeline.LayoutTypeOnlyGlobal),
		pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
			call.shader.ModuleVert().Stage(),
			call.shader.ModuleFrag().Stage(),
		}),
		pipeline.WithTopology(
			call.shader.Meta().Topology(),
			call.shader.Meta().TopologyRestartEnable(),
		),
		pipeline.WithVertexInput(
			call.shader.Meta().Bindings(),
			call.shader.Meta().Attributes(),
		),
		pipeline.WithRasterization(call.polygonMode),
		pipeline.WithColorBlend(),
		pipeline.WithMultisampling(),
	)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}

	return b
}
