package vlk

import (
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/shared/metrics"
)

// todo: api for view,projection mat4
// todo: multisampling "vulkan.SampleCount4Bit" can be used for test "ErrorDeviceLost"
// todo: set swapChain images count=1 (index out of range [1] with length 1) imageID > 0 can be created
// todo: panic after 10-15 sec in circle demo at (result := vulkan.CreateGraphicsPipelines() in internal/pipeline/factory.go)
// todo: slow points API (x500)
// todo: validation warn, in e5_blending when switch screen resolution:
//         UNASSIGNED-CoreValidation-DrawState-InvalidCommandBuffer-VkDescriptorSet(ERROR / SPEC):
//         msgNum: -396268558 - Validation Error: [
//         UNASSIGNED-CoreValidation-DrawState-InvalidCommandBuffer-VkDescriptorSet ]
//         Object 0: handle = 0x1ea2518, type = VK_OBJECT_TYPE_COMMAND_BUFFER;
//         Object 1: handle = 0x310000000031, type = VK_OBJECT_TYPE_DESCRIPTOR_SET;
//         | MessageID = 0xe8616bf2 |
//         You are adding vkCmdBindVertexBuffers() to VkCommandBuffer 0x1ea2518[]
//         that is invalid because bound VkDescriptorSet 0x310000000031[] was destroyed or updated.
// todo: new metrics api for timing groups

type (
	surfaceID uint8
)

const surfaceIdMainWindow = 0

type VLK struct {
	isReady bool
	cont    *Container

	// stats
	stats                metrics.Stats
	statsListeners       []func(metrics.Stats)
	statsUpdateFPSQueued bool

	// surfaces
	surfaceInd   surfaceID       // 0 - default (Screen, window); 1-255 reserved for user needs
	surfacesSize [255][2]float32 // width, height for each surface

	// drawing
	drawAvailable        bool
	drawImageID          uint32
	drawContext          *drawContext
	drawExecution        drawCtxFn
	drawShaderIndexesMap map[string]alloc.Allocation // shaderID -> allocation (is pointer to index buffer for this shader)'
	drawPipelineCache    map[string]pipeline.Info
}

func newVLK(cont *Container) *VLK {
	vlk := &VLK{
		isReady: true,
		cont:    cont,

		// stats
		stats:                metrics.NewStats(),
		statsListeners:       make([]func(metrics.Stats), 0),
		statsUpdateFPSQueued: false,

		// surface
		surfaceInd:   surfaceIdMainWindow,
		surfacesSize: [255][2]float32{},

		// drawing
		drawShaderIndexesMap: make(map[string]alloc.Allocation),
		drawPipelineCache:    make(map[string]pipeline.Info), // todo: is useful / fast?
	}

	// set default screen size
	wWidth, wHeight := cont.wm.GetFramebufferSize()
	vlk.surfacesSize[surfaceIdMainWindow] = [2]float32{float32(wWidth), float32(wHeight)}

	// build drawing pipeline
	vlk.initDrawingPipeline()

	// run background workers
	go vlk.countFPS()

	// return
	return vlk
}

// this will immediately stop render new frames
// wait for all current GPU work is done, then
// run mutate function, that allow change any VLK state
// after this, VLK will be turned on again
func (vlk *VLK) maintenance(mutate func()) {
	// stop render
	vlk.isReady = false

	// wait for GPU end current operations
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())

	// change vulkan state
	// rebuild pipeline, etc..
	mutate()

	// turn on back
	vlk.isReady = true
}

// ListenStats will subscribe listener to frame stats
// listener will be executed every frame with last frame Stats
func (vlk *VLK) ListenStats(listener func(stats metrics.Stats)) {
	vlk.statsListeners = append(vlk.statsListeners, listener)
}

func (vlk *VLK) countFPS() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			vlk.statsUpdateFPSQueued = true
		}
	}
}
