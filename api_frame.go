package vgl

import "github.com/go-glx/vgl/internal/gpu/vlk/metrics"

// FrameStart should be called before any drawing in current frame
func (r *Render) FrameStart() {
	r.api.FrameStart()
}

// FrameEnd should be called after any drawing in current frame
// this function will draw all queued objects in GPU
// and swap image buffer from GPU to screen
func (r *Render) FrameEnd() {
	r.api.FrameEnd()
}

// ListenStats allows to subscribe to render frame stats
// this function will execute custom callback function with
// last frame stats
func (r *Render) ListenStats(listener func(Stats)) {
	// transform internal metrics to public metrics
	r.api.ListenStats(func(stats metrics.Stats) {
		listener(Stats{
			FrameIndex:            stats.FrameIndex,
			FPS:                   stats.FPS,
			DrawCalls:             stats.DrawCalls,
			DrawChunks:            stats.DrawChunks,
			DrawUniqueShaders:     stats.DrawUniqueShaders,
			TimeCreatePipeline:    stats.TimeCreatePipeline,
			TimeFlushVertexBuffer: stats.TimeFlushVertexBuffer,
			TimeRenderInstanced:   stats.TimeRenderInstanced,
			TimeRenderFallback:    stats.TimeRenderFallback,
			Memory: MemoryStats{
				TotalCapacity: stats.Memory.TotalCapacity,
				TotalSize:     stats.Memory.TotalSize,
				IndexBuffers: UsageStats{
					PagesCount: stats.Memory.IndexBuffers.PagesCount,
					AreasCount: stats.Memory.IndexBuffers.AreasCount,
					Capacity:   stats.Memory.IndexBuffers.Capacity,
					Size:       stats.Memory.IndexBuffers.Size,
				},
				VertexBuffers: UsageStats{
					PagesCount: stats.Memory.VertexBuffers.PagesCount,
					AreasCount: stats.Memory.VertexBuffers.AreasCount,
					Capacity:   stats.Memory.VertexBuffers.Capacity,
					Size:       stats.Memory.VertexBuffers.Size,
				},
				UniformBuffers: UsageStats{
					PagesCount: stats.Memory.UniformBuffers.PagesCount,
					AreasCount: stats.Memory.UniformBuffers.AreasCount,
					Capacity:   stats.Memory.UniformBuffers.Capacity,
					Size:       stats.Memory.UniformBuffers.Size,
				},
				StorageBuffers: UsageStats{
					PagesCount: stats.Memory.StorageBuffers.PagesCount,
					AreasCount: stats.Memory.StorageBuffers.AreasCount,
					Capacity:   stats.Memory.StorageBuffers.Capacity,
					Size:       stats.Memory.StorageBuffers.Size,
				},
			},
		})
	})
}
