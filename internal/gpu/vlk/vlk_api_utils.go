package vlk

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/alloc"
	"github.com/go-glx/vgl/internal/gpu/vlk/metrics"
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

// GetSurfaceSize returns current surface size in pixels
// for default surface (screen) is window width and height
// this very fast function, will always return cached values
// and can be called thousands times per frame
func (vlk *VLK) GetSurfaceSize() (w float32, h float32) {
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

	// start command buffers
	vlk.drawImageID, vlk.drawAvailable = vlk.cont.frameManager().FrameBegin()
}

func (vlk *VLK) FrameEnd() {
	if !vlk.isReady {
		return
	}

	// draw queued shaders
	vlk.draw()

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
		case alloc.BufferTypeStorage:
			vlk.collectMemoryGroupStats(stats, &vlk.stats.Memory.StorageBuffers)
		}
	}
}

func (vlk *VLK) collectMemoryGroupStats(in alloc.GroupedStats, out *metrics.UsageStats) {
	out.Capacity += in.Capacity
	out.Size += in.Size
	out.PagesCount += int(in.TotalPages)
	out.AreasCount += int(in.TotalAreas)
}
