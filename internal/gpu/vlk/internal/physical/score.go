package physical

import (
	"fmt"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type suitableCheck struct {
	isSuitable bool
	reason     string
}

func (d *Device) score(pd *GPU) int {
	checks := []suitableCheck{
		{
			isSuitable: pd.Families.supportGraphics,
			reason:     "graphics operations required for any kind of rendering",
		},
		{
			isSuitable: pd.Families.supportPresent,
			reason:     "window present required for drawing on screen",
		},
		{
			isSuitable: pd.isSupportAllRequiredExtensions(),
			reason:     "gpu should support all required extensions",
		},
		{
			isSuitable: len(pd.SurfaceProps.formats) >= 1,
			reason:     "at least some render formats should exist",
		},
		{
			isSuitable: len(pd.SurfaceProps.presentModes) >= 1,
			reason:     "at least some present modes should exist",
		},
		{
			isSuitable: pd.SurfaceProps.RichColorSpaceFormat() != nil,
			reason:     "gpu should support rich color space rendering",
		},
	}

	// filter
	for _, check := range checks {
		if !check.isSuitable {
			d.logger.Notice(fmt.Sprintf(
				"GPU '%s' not pass check: %s",
				vkconv.VarcharAsString(pd.Props.DeviceName),
				check.reason,
			))

			return -1
		}
	}

	// score
	score := 0
	if pd.Props.DeviceType == vulkan.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	devName := strings.ToLower(vkconv.VarcharAsString(pd.Props.DeviceName))
	if strings.Contains(devName, "llvm") {
		// usually is some nvidia gpu in prime setup (intel+nvidia)
		score += 5
	}

	if strings.Contains(devName, "intel") {
		// possible some intel integrated gpu
		score += 1
	}

	return score
}
