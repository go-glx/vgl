package vulkan

import (
	"fmt"
	"math"

	"github.com/vulkan-go/vulkan"
)

func (ds *vkPhysicalDeviceSurfaceProps) richColorSpaceFormat() *vulkan.SurfaceFormat {
	for _, surfaceFormat := range ds.formats {
		if surfaceFormat.Format != vulkan.FormatB8g8r8a8Srgb {
			continue
		}

		if surfaceFormat.ColorSpace != vulkan.ColorSpaceSrgbNonlinear {
			continue
		}

		return &surfaceFormat
	}

	return nil
}

func (ds *vkPhysicalDeviceSurfaceProps) bestPresentMode(vSync bool) vulkan.PresentMode {
	for _, mode := range ds.presentModes {
		if vSync && mode == vulkan.PresentModeFifo {
			return mode
		}

		if mode == vulkan.PresentModeMailbox {
			return mode
		}

		return mode
	}

	panic(fmt.Errorf("GPU not support any present mode"))
}

func (ds *vkPhysicalDeviceSurfaceProps) chooseSwapExtent(width, height uint32) vulkan.Extent2D {
	calculatedMax := uint32((math.MaxInt32 * 2) + 1)

	// Vulkan tells us to match the resolution of the window by setting the width and height in the currentExtent member.
	// However, some window managers do allow us to differ here and this is indicated by setting the width and height
	// in currentExtent to a special value: the maximum value of uint32_t.
	//
	// In that case we'll pick the resolution that best matches the window within the minImageExtent and maxImageExtent bounds.
	// But we must specify the resolution in the correct unit.

	if width == 0 || height == 0 {
		if ds.capabilities.CurrentExtent.Width != calculatedMax {
			curr := ds.capabilities.CurrentExtent
			curr.Deref()
			return curr
		}
	}

	actualExtent := vulkan.Extent2D{
		Width:  width,
		Height: height,
	}

	maxWidth := ds.capabilities.MaxImageExtent.Width
	maxHeight := ds.capabilities.MaxImageExtent.Height

	if maxWidth == 0 || maxHeight == 0 {
		return actualExtent
	}

	actualExtent.Width = vkClampUint(
		actualExtent.Width,
		ds.capabilities.MinImageExtent.Width,
		ds.capabilities.MaxImageExtent.Width,
	)
	actualExtent.Height = vkClampUint(
		actualExtent.Height,
		ds.capabilities.MinImageExtent.Height,
		ds.capabilities.MaxImageExtent.Height,
	)

	return actualExtent
}

func (ds *vkPhysicalDeviceSurfaceProps) imageBuffersCount() uint32 {
	optimalCount := ds.capabilities.MinImageCount + 1

	if ds.capabilities.MaxImageCount == 0 {
		// no maximum
		return optimalCount
	}

	if optimalCount > ds.capabilities.MaxImageCount {
		// clamp to max
		optimalCount = ds.capabilities.MaxImageCount
	}

	return optimalCount
}
