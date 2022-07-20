package physical

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

func (d *Device) score(pd *GPU) int {
	required := map[bool]string{
		// families
		!pd.Families.supportGraphics: "graphics operations not supported",
		!pd.Families.supportPresent:  "window present not supported",

		// extensions
		!pd.isSupportAllRequiredExtensions(): "not all required extensions supported",

		// swap chain
		len(pd.SurfaceProps.formats) <= 0:             "not GPU",
		len(pd.SurfaceProps.presentModes) <= 0:        "not GPU",
		pd.SurfaceProps.richColorSpaceFormat() == nil: "rich colorSpace not supported",
	}

	// filter
	for failed, reason := range required {
		if failed {
			log.Printf(
				"vk: GPU '%s' not pass check: %s\n",
				vkconv.VarcharAsString(pd.Props.DeviceName),
				reason,
			)

			return -1
		}
	}

	// score
	score := 0
	if pd.Props.DeviceType == vulkan.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	return score
}
