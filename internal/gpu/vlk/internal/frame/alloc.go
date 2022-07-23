package frame

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

// Semaphore is GPU <-> GPU sync
// Fence     is CPU <-> GPU sync

func allocateSemaphore(ld *logical.Device) vulkan.Semaphore {
	createInfo := &vulkan.SemaphoreCreateInfo{
		SType: vulkan.StructureTypeSemaphoreCreateInfo,
	}

	var ref vulkan.Semaphore
	must.Work(vulkan.CreateSemaphore(ld.Ref(), createInfo, nil, &ref))

	return ref
}

func allocateFence(ld *logical.Device) vulkan.Fence {
	createInfo := &vulkan.FenceCreateInfo{
		SType: vulkan.StructureTypeFenceCreateInfo,
		Flags: vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit),
	}

	var fence vulkan.Fence
	must.Work(vulkan.CreateFence(ld.Ref(), createInfo, nil, &fence))

	return fence
}
