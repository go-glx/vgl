package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/arch"
)

func newSurfaceFromWindow(inst *vkInstance, wm arch.WindowManager) *vkSurface {
	surface, err := wm.CreateSurface(inst.ref)
	if err != nil {
		panic(fmt.Errorf("failed create vulkan windows surface: %w", err))
	}

	return &vkSurface{
		inst: inst,
		ref:  surface,
	}
}

func (surf *vkSurface) free() {
	vulkan.DestroySurface(surf.inst.ref, surf.ref, nil)

	log.Printf("vk: freed: window surface\n")
}
