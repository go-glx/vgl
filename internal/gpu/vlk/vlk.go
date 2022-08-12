package vlk

import (
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
)

// todo: fullScreen switch panics in examples
// todo: vertex buffers is override in swapChain (can have multiple frames in-fly)
// todo: panic in rect scene after one second

type VLK struct {
	isReady bool
	cont    *Container

	// stats
	stats          glm.Stats
	statsListeners []func(glm.Stats)

	// surfaces
	surfaceInd   uint8          // 0 - default (Screen, window); 1-255 reserved for user needs
	surfacesSize [255][2]uint32 // width, height for each surface

	// drawing
	shaderIndexPtr map[string]alloc.AllocationID // shaderID -> allocationID (is pointer to index buffer for this shader)
	currentBatch   *drawCall
	queue          []drawCall
}

func newVLK(cont *Container) *VLK {
	vlk := &VLK{
		isReady: true,
		cont:    cont,

		// stats
		stats:          glm.Stats{},
		statsListeners: make([]func(glm.Stats), 0),

		// surface
		surfaceInd:   0, // default - screen
		surfacesSize: [255][2]uint32{},

		// drawing
		shaderIndexPtr: make(map[string]alloc.AllocationID),
		currentBatch:   &drawCall{},
		queue:          make([]drawCall, 0, queueCapacity),
	}

	// set default screen size
	wWidth, wHeight := cont.wm.GetFramebufferSize()
	vlk.surfacesSize[0] = [2]uint32{uint32(wWidth), uint32(wHeight)}

	go countFPS(&vlk.stats)

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
func (vlk *VLK) ListenStats(listener func(stats glm.Stats)) {
	vlk.statsListeners = append(vlk.statsListeners, listener)
}

func countFPS(stats *glm.Stats) {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			stats.FPS = stats.FrameIndex
			stats.FrameIndex = 0
		}
	}
}
