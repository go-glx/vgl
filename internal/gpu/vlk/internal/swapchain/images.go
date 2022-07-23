package swapchain

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

func createImages(swapChain vulkan.Swapchain, ld *logical.Device) []vulkan.Image {
	imagesCount := uint32(0)
	must.Work(vulkan.GetSwapchainImages(ld.Ref(), swapChain, &imagesCount, nil))

	if imagesCount == 0 {
		panic(fmt.Errorf("swapchain should have at least 1 image buffer"))
	}

	images := make([]vulkan.Image, imagesCount)
	must.Work(vulkan.GetSwapchainImages(ld.Ref(), swapChain, &imagesCount, images))

	return images
}

func createViews(images []vulkan.Image, ld *logical.Device, props ChainProps) []vulkan.ImageView {
	views := make([]vulkan.ImageView, 0, len(images))

	for _, image := range images {
		views = append(views, createView(image, ld, props))
	}

	return views
}

func createView(image vulkan.Image, ld *logical.Device, props ChainProps) vulkan.ImageView {
	info := &vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    image,
		ViewType: vulkan.ImageViewType2d,
		Format:   props.ImageFormat,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var view vulkan.ImageView
	must.Work(vulkan.CreateImageView(ld.Ref(), info, nil, &view))

	return view
}
