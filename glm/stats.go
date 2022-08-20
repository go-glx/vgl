package glm

import "time"

// Stats for previous rendered frame
type Stats struct {
	FrameIndex int
	FPS        int

	DrawCalls         int
	DrawChunks        int
	DrawUniqueShaders int

	TimeCreatePipeline    time.Duration
	TimeFlushVertexBuffer time.Duration
	TimeRenderInstanced   time.Duration
	TimeRenderFallback    time.Duration

	Memory MemoryStats
}

type MemoryStats struct {
	TotalCapacity  uint32 // total application memory required
	TotalSize      uint32 // total application allocated memory
	IndexBuffers   UsageStats
	VertexBuffers  UsageStats
	UniformBuffers UsageStats
}

type UsageStats struct {
	PagesCount int    // How many pages is allocated (one pages has many areas)
	AreasCount int    // How many areas is allocated (area is logical abstraction from vulkan buffer)
	Capacity   uint32 // total required capacity of this usage type (is not total device capacity)
	Size       uint32 // total size (sum of all allocations of this usage type)
}

func (s *Stats) Reset() {
	// frameIndex and FPS should not be reset

	s.DrawCalls = 0
	s.DrawChunks = 0
	s.DrawUniqueShaders = 0

	s.TimeCreatePipeline = 0
	s.TimeFlushVertexBuffer = 0
	s.TimeRenderInstanced = 0
	s.TimeRenderFallback = 0

	s.Memory = MemoryStats{}
}
