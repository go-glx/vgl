package metrics

import "time"

// Stats for previous rendered frame
type Stats struct {
	FrameIndex int
	FPS        int

	DrawCalls         int
	DrawChunks        int
	DrawUniqueShaders int

	TimeCreatePipeline    time.Duration // deprecated, todo: remove
	TimeFlushVertexBuffer time.Duration // deprecated, todo: remove
	TimeRenderInstanced   time.Duration // deprecated, todo: remove
	TimeRenderFallback    time.Duration // deprecated, todo: remove

	Memory MemoryStats
}

type MemoryStats struct {
	TotalCapacity  uint32 // total application memory required
	TotalSize      uint32 // total application allocated memory
	IndexBuffers   UsageStats
	VertexBuffers  UsageStats
	UniformBuffers UsageStats
	StorageBuffers UsageStats
}

type UsageStats struct {
	PagesCount int    // How many pages is allocated (one pages has many areas)
	AreasCount int    // How many areas is allocated (area is logical abstraction from vulkan buffer)
	Capacity   uint32 // total required capacity of this usage type (is not total device capacity)
	Size       uint32 // total size (sum of all allocations of this usage type)
}
