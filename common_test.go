package govgl

import (
	"fmt"
	"testing"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/arch"
	"github.com/fe3dback/govgl/config"
)

type virtualWM struct {
}

func (v *virtualWM) CreateSurface(_ vulkan.Instance) (vulkan.Surface, error) {
	return vulkan.NullSurface, nil
}

func (v *virtualWM) InitVulkanProcAddr() {
	err := vulkan.SetDefaultGetInstanceProcAddr()
	if err != nil {
		panic(fmt.Errorf("failed get vulkan proc addr: %w", err))
	}
}

func (v *virtualWM) AppName() string {
	return "govgl_test"
}

func (v *virtualWM) EngineName() string {
	return "govgl_test"
}

func (v *virtualWM) GetRequiredInstanceExtensions() []string {
	return []string{"VK_KHR_surface"}
}

func (v *virtualWM) GetFramebufferSize() (width, height int) {
	return 1280, 720
}

func (v *virtualWM) OnWindowResized(f func(width int, height int)) {
	// no virtual events
	return
}

func (v *virtualWM) Close() error {
	return nil
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDebug(true),
			)

			wm := arch.NewGLFW("test", "test", false, 320, 240)

			renderer := NewRender(wm, cfg)
			_ = renderer.Close()
		})
	}
}
