package frame

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

type Manager struct {
	chain          *swapchain.Chain
	mainRenderPass *renderpass.Pass

	id      uint8
	count   uint8
	frames  []*frame
	current *frame

	ld *logical.Device
}

func NewManager(ld *logical.Device, pool *command.Pool, chain *swapchain.Chain, renderToScreenPass *renderpass.Pass, onSuboptimal func()) *Manager {
	count := pool.BuffersCount()
	frames := make([]*frame, 0, count)

	for buffID := 0; buffID < count; buffID++ {
		frames = append(frames, newFrame(
			ld,
			chain,
			pool.CommandBuffer(buffID),
			onSuboptimal,
		))
	}

	log.Printf("vk: frame manager created\n")

	const initialID = 0
	return &Manager{
		chain:          chain,
		mainRenderPass: renderToScreenPass,

		id:      initialID,
		count:   uint8(count),
		frames:  frames,
		current: frames[initialID],

		ld: ld,
	}
}

func (m *Manager) Free() {
	for _, frame := range m.frames {
		frame.free()
	}

	log.Printf("vk: freed: frames manager\n")
}

func (m *Manager) FrameBegin() {
	// switch to next frame
	m.id = (m.id + 1) % m.count
	m.current = m.frames[m.id]

	// process frame
	m.current.frameBegin()

	m.FrameApplyCommands(func(cb vulkan.CommandBuffer) {
		m.renderPassMainBegin(cb)
	})
}

func (m *Manager) FrameApplyCommands(apply func(cb vulkan.CommandBuffer)) {
	m.current.frameWrite(apply)
}

func (m *Manager) FrameEnd() {
	m.FrameApplyCommands(func(cb vulkan.CommandBuffer) {
		m.renderPassMainEnd(cb)
	})

	m.current.frameEnd()
}
