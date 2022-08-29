package vlk

import (
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/metrics"
)

// todo: multisampling "vulkan.SampleCount4Bit" can be used for test "ErrorDeviceLost"
// todo: set swapChain images count=1 (index out of range [1] with length 1) imageID > 0 can be created
// todo: panic after 10-15 sec in circle demo at (result := vulkan.CreateGraphicsPipelines() in internal/pipeline/factory.go)

type VLK struct {
	isReady bool
	cont    *Container

	// stats
	stats            metrics.Stats
	statsListeners   []func(metrics.Stats)
	statsResetQueued bool

	// surfaces
	surfaceInd   uint8           // 0 - default (Screen, window); 1-255 reserved for user needs
	surfacesSize [255][2]float32 // width, height for each surface

	// drawing
	shaderIndexPtr map[string]alloc.Allocation // shaderID -> allocation (is pointer to index buffer for this shader)
	currentBatch   *drawCall
	queue          []drawCall
}

func newVLK(cont *Container) *VLK {
	vlk := &VLK{
		isReady: true,
		cont:    cont,

		// stats
		stats:            metrics.Stats{},
		statsListeners:   make([]func(metrics.Stats), 0),
		statsResetQueued: false,

		// surface
		surfaceInd:   0, // default - screen
		surfacesSize: [255][2]float32{},

		// drawing
		shaderIndexPtr: make(map[string]alloc.Allocation),
		currentBatch:   &drawCall{},
		queue:          make([]drawCall, 0, queueCapacity),
	}

	// set default screen size
	wWidth, wHeight := cont.wm.GetFramebufferSize()
	vlk.surfacesSize[0] = [2]float32{float32(wWidth), float32(wHeight)}

	go vlk.countFPS()
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
			vlk.statsResetQueued = true
		}
	}
}
