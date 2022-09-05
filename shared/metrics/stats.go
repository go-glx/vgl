package metrics

import "time"

// Stats for previous rendered frame
type (
	Stats struct {
		FrameIndex int
		FPS        int

		DrawCalls  int
		DrawGroups int

		SegmentDuration map[string]time.Duration
		Memory          MemoryStats
	}

	MemoryStats struct {
		TotalCapacity  uint32 // total application memory required
		TotalSize      uint32 // total application allocated memory
		IndexBuffers   UsageStats
		VertexBuffers  UsageStats
		UniformBuffers UsageStats
		StorageBuffers UsageStats
	}

	UsageStats struct {
		PagesCount int    // How many pages is allocated (one pages has many areas)
		AreasCount int    // How many areas is allocated (area is logical abstraction from vulkan buffer)
		Capacity   uint32 // total required capacity of this usage type (is not total device capacity)
		Size       uint32 // total size (sum of all allocations of this usage type)
	}
)

func NewStats() Stats {
	return Stats{
		SegmentDuration: map[string]time.Duration{},
	}
}

func (s *Stats) Reset() {
	s.DrawCalls = 0
	s.DrawGroups = 0

	for segment := range s.SegmentDuration {
		s.SegmentDuration[segment] = 0
	}

	s.Memory = MemoryStats{}
}
