package vlkext

import (
	"io"
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

type (
	WindowManager interface {
		io.Closer

		AppName() string
		EngineName() string
		OnWindowResized(func(width, height int))
		CreateSurface(inst vulkan.Instance) (vulkan.Surface, error)
		GetRequiredInstanceExtensions() []string
		GetFramebufferSize() (width, height int)
		InitVulkanProcAddr() unsafe.Pointer
	}
)
