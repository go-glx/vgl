package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/glm"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
)

// WarmUp will warm vlk renderer and create all needed
// objects for work, this must be called one time
// before first FrameStart
func (vlk *VLK) WarmUp() {
	// request some managers, this will create it
	// and all dependencies, like swapChain, renderPass, etc..
	_ = vlk.cont.frameManager()
	_ = vlk.cont.shaderManager()

	// preload shader indexes
	for _, shaderID := range preloadIndexShaders {
		vlk.preloadShaderIndexes(vlk.cont.shaderManager().ShaderByID(shaderID))
	}
}

func (vlk *VLK) GPUWait() {
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())
}

// GetSurfaceSize returns current surface size in pixels
// for default surface (screen) is window width and height
// this very fast function, will always return cached values
// and can be called thousands times per frame
func (vlk *VLK) GetSurfaceSize() (w uint32, h uint32) {
	return vlk.surfacesSize[vlk.surfaceInd][0],
		vlk.surfacesSize[vlk.surfaceInd][1]
}

func (vlk *VLK) FrameStart() {
	if !vlk.isReady {
		return
	}

	if vlk.statsResetQueued {
		vlk.stats.FPS = vlk.stats.FrameIndex
		vlk.stats.FrameIndex = -1
		vlk.statsResetQueued = false
	}

	// clear stats from prev frame
	vlk.stats.Reset()

	// start command buffers
	vlk.cont.frameManager().FrameBegin()
}

func (vlk *VLK) FrameEnd() {
	if !vlk.isReady {
		return
	}

	// draw queued shaders
	vlk.drawAll()

	// submit command buffers
	vlk.cont.frameManager().FrameEnd()

	// collect memory stats and then clean garbage
	vlk.collectMemoryStats()
	vlk.cont.allocHeap().GarbageCollect()

	// send stats
	for _, listener := range vlk.statsListeners {
		listener(vlk.stats)
	}
}

func (vlk *VLK) collectMemoryStats() {
	memStats := vlk.cont.allocHeap().Stats()
	vlk.stats.Memory.TotalCapacity = memStats.TotalCapacity
	vlk.stats.Memory.TotalSize = memStats.TotalSize

	for _, stats := range memStats.Grouped {
		switch stats.BufferType {
		case alloc.BufferTypeIndex:
			vlk.collectMemoryGroupStats(stats, &vlk.stats.Memory.IndexBuffers)
		case alloc.BufferTypeVertex:
			vlk.collectMemoryGroupStats(stats, &vlk.stats.Memory.VertexBuffers)
		case alloc.BufferTypeUniform:
			vlk.collectMemoryGroupStats(stats, &vlk.stats.Memory.UniformBuffers)
		}
	}
}

func (vlk *VLK) collectMemoryGroupStats(in alloc.GroupedStats, out *glm.UsageStats) {
	out.Capacity += in.Capacity
	out.Size += in.Size
	out.PagesCount += int(in.TotalPages)
	out.AreasCount += int(in.TotalAreas)
}
