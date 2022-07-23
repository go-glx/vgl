package swapchain

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type ChainProps struct {
	ImageFormat     vulkan.Format
	ImageColorSpace vulkan.ColorSpace
	BufferSize      vulkan.Extent2D
	PresentMode     vulkan.PresentMode
	BuffersCount    uint32
}

func newProps(width, height uint32, pd *physical.Device, mobileFriendly bool) ChainProps {
	gpuProps := pd.PrimaryGPU().SurfaceProps
	richColorFormat := gpuProps.RichColorSpaceFormat()

	return ChainProps{
		ImageFormat:     richColorFormat.Format,
		ImageColorSpace: richColorFormat.ColorSpace,
		BufferSize:      gpuProps.ChooseSwapExtent(width, height),
		PresentMode:     gpuProps.BestPresentMode(mobileFriendly),
		BuffersCount:    gpuProps.ConcurrentBuffersCount(),
	}
}

func (p *ChainProps) String() string {
	return fmt.Sprintf("format=%s, colorSpace=%s, buffersCount=%s, bufferSize=%s, presentMode=%s",
		p.formatString(),
		p.colorSpaceString(),
		p.buffersCountString(),
		p.bufferSizeString(),
		p.presentModeString(),
	)
}

func (p *ChainProps) formatString() string {
	if name, ok := propsImageFormats[p.ImageFormat]; ok {
		return name
	}

	return "unknown"
}

func (p *ChainProps) colorSpaceString() string {
	if name, ok := propsColorSpaces[p.ImageColorSpace]; ok {
		return name
	}

	return "unknown"
}

func (p *ChainProps) bufferSizeString() string {
	return fmt.Sprintf("[%dx%d]", p.BufferSize.Width, p.BufferSize.Height)
}

func (p *ChainProps) presentModeString() string {
	if name, ok := propsPresentModes[p.PresentMode]; ok {
		return name
	}

	return "unknown"
}

func (p *ChainProps) buffersCountString() string {
	return fmt.Sprintf("%d", p.BuffersCount)
}
