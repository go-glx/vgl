package physical

import (
	"fmt"
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type (
	SurfaceProps struct {
		capabilities   vulkan.SurfaceCapabilities
		formats        []vulkan.SurfaceFormat
		presentModes   []vulkan.PresentMode
		surfaceSupport bool
	}
)

// ConcurrentBuffersCount defines how much
//  - swap chain images
//  - command buffers in pool
//  - etc...
// we need in pipeline process.
// Usually we need at least two images for [rendering, display].
//
// Each frame, this buffers will be swapped.
// but for good quality we may use triple buffering
func (ds *SurfaceProps) ConcurrentBuffersCount() uint32 {
	return vkconv.ClampUint(
		def.OptimalSwapChainBuffersCount,
		ds.capabilities.MinImageCount,
		ds.capabilities.MaxImageCount,
	)
}

func (ds *SurfaceProps) RichColorSpaceFormat() *vulkan.SurfaceFormat {
	for _, surfaceFormat := range ds.formats {
		if surfaceFormat.Format != def.SurfaceFormat {
			continue
		}

		if surfaceFormat.ColorSpace != def.SurfaceColorSpace {
			continue
		}

		return &surfaceFormat
	}

	return nil
}

func (ds *SurfaceProps) bestPresentMode(mobileFriendly bool) vulkan.PresentMode {
	for _, mode := range ds.presentModes {
		if mobileFriendly && mode == vulkan.PresentModeFifo {
			// 1# [R]       [D.......][R]
			// 2# [D.......][R]       [D.......]
			//    |         |         |
			//   frm1      frm 2     frm 3
			//
			// Fifo is low consumption rendering, friendly for
			// mobile devices. When we render (R) some buffer, it will
			// stay all time before displayed (D) on screen.
			//
			// + vsync
			// + good for mobiles (low power consumption)
			// + always supported
			// - high latency
			return mode
		}

		if mode == vulkan.PresentModeMailbox {
			// 1# [R][R][R] [D] ///// [R]
			// 2# [D] ///// [R][R][R] [D]
			//    |         |         |
			//   frm1      frm 2     frm 3
			//
			// Mailbox is high power consumption mode, that will
			// heat GPU to max, and re-render (R) all frames when
			// GPU has free time (idle). It`s a lot of work, that
			// will be ignored. For example, we render (R) 3 buffers per frame
			// here, and display (D) only last of them.
			//
			// + low latency
			// - not always supported on all GPU's
			// - high power consumption
			return mode
		}

		return mode
	}

	panic(fmt.Errorf("GPU not support any present mode"))
}

func (ds *SurfaceProps) chooseSwapExtent(width, height uint32) vulkan.Extent2D {
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

	actualExtent.Width = vkconv.ClampUint(
		actualExtent.Width,
		ds.capabilities.MinImageExtent.Width,
		ds.capabilities.MaxImageExtent.Width,
	)
	actualExtent.Height = vkconv.ClampUint(
		actualExtent.Height,
		ds.capabilities.MinImageExtent.Height,
		ds.capabilities.MaxImageExtent.Height,
	)

	return actualExtent
}
