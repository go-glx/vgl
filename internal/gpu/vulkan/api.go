package vulkan

import (
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
)

func (vk *Vk) appendToRenderQueue(sp shaderProgram) {
	if queue, exist := vk.renderQueue[sp.ID()]; exist {
		vk.renderQueue[sp.ID()] = append(queue, sp)
		return
	}

	vk.renderQueue[sp.ID()] = []shaderProgram{sp}
}

func (vk *Vk) FrameStart() {
	if vk.inResizing || vk.isMinimized {
		vk.currentFrameAvailableForRender = false
		return
	}

	// 1. start frame
	vk.currentFrameImageID, vk.currentFrameAvailableForRender = vk.frameManager.frameStart(vk.swapChain)
	if !vk.currentFrameAvailableForRender {
		return
	}

	// 2. start command buffer
	vk.renderQueue = make(map[string][]shaderProgram)
	vk.commandPool.commandBufferStart(int(vk.currentFrameImageID))

	// 3. start render pass
	commandBuffer := vk.commandPool.commandBuffer(int(vk.currentFrameImageID))
	renderPass := vk.container.defaultRenderPass()
	vk.frameBuffers.renderPassStart(int(vk.currentFrameImageID), commandBuffer, renderPass)
}

func (vk *Vk) FrameEnd() {
	if !vk.currentFrameAvailableForRender {
		return
	}

	commandBuffer := vk.commandPool.commandBuffer(int(vk.currentFrameImageID))

	// 3. end render pass
	vk.frameBuffers.renderPassEnd(commandBuffer)

	// 2. end command buffer
	vk.commandPool.commandBufferEnd(int(vk.currentFrameImageID))

	// 1. end frame
	vk.frameManager.frameEnd(vk.swapChain, commandBuffer)
}

// todo: remove this temp stuff
func mat4Perspective(fovy, aspect, near, far float32) glm.Mat4 {
	nmf, f := 1/(near-far), 1./math.Tan(float64(fovy)/2.0)

	return glm.Mat4{
		A: glm.Vec4{f / float64(aspect), 0, 0, 0},
		B: glm.Vec4{0, f, 0, 0},
		C: glm.Vec4{0, 0, float64((near + far) * nmf), -1},
		D: glm.Vec4{0, 0, float64((2. * far * near) * nmf), 0},
	}
}

func (vk *Vk) Draw() {
	if !vk.currentFrameAvailableForRender {
		return
	}

	// todo: uniform buffers for per instance data (transform, rotation, scale)
	// todo: global uniform buffer for mat4 (projection, view) - single frame buffer for all objects

	// reset
	drawCalls := 0
	drawInstances := 0
	commandBuffer := vk.commandPool.commandBuffer(int(vk.currentFrameImageID))

	// write to global buffers
	// uboProjection := mat4Perspective(
	// 	float32(glm.NewAngle(45).Radians()),
	// 	vk.swapChain.viewport().Width/vk.swapChain.viewport().Height,
	// 	0.1,
	// 	10,
	// )
	// uboProjection := glm.Mat4Identity()
	// uboView := glm.Mat4Identity()
	// uboModel := glm.Mat4Identity()

	// vk.dataBuffersManager.updateGlobalUniformBuffer(vk.currentFrameImageID, uboProjection, uboView, uboModel)

	// draw all shaders
	for pipelineID, instances := range vk.renderQueue {
		// bind pipeline for this group of shaders
		pipeline := vk.pipelineManager.pipeline(pipelineID)
		vulkan.CmdBindPipeline(commandBuffer, vulkan.PipelineBindPointGraphics, pipeline)

		// clean buffers
		vk.dataBuffersManager.resetBuffers()

		// copy data to buffers
		for _, instance := range instances {
			vk.dataBuffersManager.writeToBuffers(instance)
			drawInstances++
		}

		// flush buffers to GPU memory
		result := vk.dataBuffersManager.flushBuffers()

		// bind indexes
		vulkan.CmdBindIndexBuffer(commandBuffer, result.indexBuffer, 0, vulkan.IndexTypeUint16)
		indexesCount := uint32(len(instances[0].Indexes()))

		// draw commands
		for _, batch := range result.vertexChunks {
			vulkan.CmdBindVertexBuffers(commandBuffer, 0, uint32(1), batch.buffers, batch.offsets)
			vulkan.CmdDrawIndexed(commandBuffer, indexesCount*batch.instanceCount, 1, 0, 0, 0)
			drawCalls++
		}
	}

	// log.Printf("draw calls=%d, instances=%d\n", drawCalls, drawInstances)
}

func (vk *Vk) GPUWait() {
	vulkan.DeviceWaitIdle(vk.ld.ref)
}
