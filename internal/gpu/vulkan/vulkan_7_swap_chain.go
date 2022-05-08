package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/config"
)

func newSwapChain(width, height uint32, pd *vkPhysicalDevice, ld *vkLogicalDevice, surface *vkSurface, cfg *config.Config) *vkSwapChain {
	info := vkSwapChainInfo{
		imageFormat:     pd.surfaceProps.richColorSpaceFormat().Format,
		imageColorSpace: pd.surfaceProps.richColorSpaceFormat().ColorSpace,
		bufferSize:      pd.surfaceProps.chooseSwapExtent(width, height),
		presentMode:     pd.surfaceProps.bestPresentMode(cfg.HasGPUVSync()),
	}

	uniqFamilies := pd.families.uniqueIDs()
	sharingMode := vulkan.SharingModeExclusive
	if len(uniqFamilies) > 1 {
		sharingMode = vulkan.SharingModeConcurrent
	}

	// assemble create request
	swapChainCreateInfo := &vulkan.SwapchainCreateInfo{
		SType:                 vulkan.StructureTypeSwapchainCreateInfo,
		Surface:               surface.ref,
		MinImageCount:         pd.surfaceProps.imageBuffersCount(),
		ImageFormat:           info.imageFormat,
		ImageColorSpace:       info.imageColorSpace,
		ImageExtent:           info.bufferSize,
		ImageArrayLayers:      1,
		ImageUsage:            vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		ImageSharingMode:      sharingMode,
		QueueFamilyIndexCount: uint32(len(uniqFamilies)),
		PQueueFamilyIndices:   uniqFamilies,
		PreTransform:          pd.surfaceProps.capabilities.CurrentTransform,
		CompositeAlpha:        vulkan.CompositeAlphaOpaqueBit,
		PresentMode:           info.presentMode,
		Clipped:               vulkan.True,
	}

	// allocate swapChain
	var swapChain vulkan.Swapchain
	vkAssert(
		vulkan.CreateSwapchain(ld.ref, swapChainCreateInfo, nil, &swapChain),
		fmt.Errorf("failed create swapChain"),
	)

	// allocate swap images
	images := createSwapChainImages(swapChain, ld)
	imagesView := createSwapChainImageViews(images, ld, pd)

	log.Printf("vk: swapchain created, images=%d, info=(%s)\n", len(images), info.String())

	return &vkSwapChain{
		ld:           ld,
		ref:          swapChain,
		surfaceProps: pd.surfaceProps,
		images:       images,
		imagesView:   imagesView,
		info:         info,
	}
}

func (sc *vkSwapChain) free() {
	for _, view := range sc.imagesView {
		vulkan.DestroyImageView(sc.ld.ref, view, nil)
	}

	vulkan.DestroySwapchain(sc.ld.ref, sc.ref, nil)
	log.Printf("vk: freed: swapchain\n")
}

func (sc *vkSwapChain) viewport() vulkan.Viewport {
	return vulkan.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(sc.info.bufferSize.Width),
		Height:   float32(sc.info.bufferSize.Height),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}
}

func (sc *vkSwapChain) scissor() vulkan.Rect2D {
	return vulkan.Rect2D{
		Offset: vulkan.Offset2D{
			X: 0,
			Y: 0,
		},
		Extent: sc.info.bufferSize,
	}
}

func createSwapChainImages(swapChain vulkan.Swapchain, ld *vkLogicalDevice) []vulkan.Image {
	imagesCount := uint32(0)
	vkAssert(
		vulkan.GetSwapchainImages(ld.ref, swapChain, &imagesCount, nil),
		fmt.Errorf("failed fetch swapChain images"),
	)
	if imagesCount == 0 {
		panic(fmt.Errorf("swapchain should have at least 1 image buffer"))
	}

	images := make([]vulkan.Image, imagesCount)
	vkAssert(
		vulkan.GetSwapchainImages(ld.ref, swapChain, &imagesCount, images),
		fmt.Errorf("failed fetch swapChain images"),
	)

	return images
}

func createSwapChainImageViews(images []vulkan.Image, ld *vkLogicalDevice, pd *vkPhysicalDevice) []vulkan.ImageView {
	views := make([]vulkan.ImageView, 0, len(images))

	for _, image := range images {
		views = append(views, createSwapChainImageView(image, ld, pd))
	}

	return views
}

func createSwapChainImageView(image vulkan.Image, ld *vkLogicalDevice, pd *vkPhysicalDevice) vulkan.ImageView {
	createInfo := &vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    image,
		ViewType: vulkan.ImageViewType2d,
		Format:   pd.surfaceProps.richColorSpaceFormat().Format,
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

	var imageView vulkan.ImageView
	vkAssert(
		vulkan.CreateImageView(ld.ref, createInfo, nil, &imageView),
		fmt.Errorf("failed create image view"),
	)

	return imageView
}
