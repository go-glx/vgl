package physical

import (
	"fmt"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type (
	GPU struct {
		logger             config.Logger
		Ref                vulkan.PhysicalDevice
		Props              vulkan.PhysicalDeviceProperties
		Features           vulkan.PhysicalDeviceFeatures
		MemProperties      vulkan.PhysicalDeviceMemoryProperties
		Extensions         []vulkan.ExtensionProperties
		Families           Families
		SurfaceProps       SurfaceProps
		RequiredExtensions []string
	}
)

func (pd *GPU) isSupportAllRequiredExtensions() bool {
	supportedExt := make(map[string]any)

	for _, extension := range pd.Extensions {
		vkExtName := vkconv.VarcharAsString(extension.ExtensionName)
		supportedExt[vkExtName] = struct{}{}
	}

	notSupported := make([]string, 0)
	for _, extension := range def.RequiredDeviceExtensions {
		vkExtName := vkconv.NormalizeString(extension)

		if _, supported := supportedExt[vkExtName]; supported {
			continue
		}

		notSupported = append(notSupported, extension)
	}

	if len(notSupported) > 0 {
		pd.logger.Notice(fmt.Sprintf(
			"vk: GPU '%s' not support all required extensions: [%s]\n",
			vkconv.VarcharAsString(pd.Props.DeviceName),
			strings.Join(notSupported, ", "),
		))

		return false
	}

	return true
}
