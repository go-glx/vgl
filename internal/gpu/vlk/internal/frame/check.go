package frame

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func (m *Manager) notice(result vulkan.Result) bool {
	return must.NotCare(m.logger, result)
}
