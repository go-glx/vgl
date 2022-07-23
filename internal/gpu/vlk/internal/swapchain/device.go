package swapchain

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

func deviceSharingMode(pd *physical.Device) vulkan.SharingMode {
	if len(pd.PrimaryGPU().Families.UniqueIDs()) > 1 {
		return vulkan.SharingModeConcurrent
	}

	return vulkan.SharingModeExclusive
}
