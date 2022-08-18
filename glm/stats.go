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

	// todo: memory size & capacity (grouped by buffer type (vert,ind,ubo..)
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
}
