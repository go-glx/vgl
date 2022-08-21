package descriptors

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
)

const (
	layoutIndexGlobal   layoutIndex = 0
	layoutIndexMaterial layoutIndex = 1
	layoutIndexInstance layoutIndex = 2
)

const (
	// # layout 0 (global)
	bindingGlobalUniforms          bindingIndex = 0
	bindingGlobalSurfaceSize       bindingIndex = 1
	bindingGlobalLightDataUniforms bindingIndex = 2

	// # layout 1 (material)
	bindingMaterialDataStorage bindingIndex = 0
	bindingMaterialTextures    bindingIndex = 1

	// # layout 2 (instance)
	bindingInstanceModelMatrix bindingIndex = 0
)

type (
	Blueprint struct {
		logger config.Logger
		ld     *logical.Device

		layouts layouts
	}

	layoutIndex    uint8
	bindingIndex   uint8
	layouts        map[layoutIndex]BlueprintLayout
	layoutBindings map[bindingIndex]vulkan.DescriptorSetLayoutBinding
)

func NewBlueprint(logger config.Logger, ld *logical.Device) *Blueprint {
	bpLayouts := layouts{
		layoutIndexGlobal: layoutGlobal(ld.Ref()),
	}

	logger.Debug("descriptor set blueprint:")
	for _, layout := range bpLayouts {
		logger.Debug(fmt.Sprintf("- %s", layout.String()))
		logger.Debug(fmt.Sprintf("  %s", layout.usage))
	}

	return &Blueprint{
		logger:  logger,
		ld:      ld,
		layouts: bpLayouts,
	}
}

func (bp *Blueprint) Free() {
	for _, layout := range bp.layouts {
		vulkan.DestroyDescriptorSetLayout(bp.ld.Ref(), layout.layout, nil)
		bp.logger.Debug(fmt.Sprintf("freed: layout: %s", layout.title))
	}
}

func (bp *Blueprint) LayoutGlobal() BlueprintLayout {
	return bp.layouts[layoutIndexGlobal]
}

func (bp *Blueprint) LayoutMaterial() BlueprintLayout {
	// todo: required for 3D objects
	panic(fmt.Errorf("not implemented yet"))
}

func (bp *Blueprint) LayoutLocal() BlueprintLayout {
	// todo: required for 3D objects
	panic(fmt.Errorf("not implemented yet"))
}
