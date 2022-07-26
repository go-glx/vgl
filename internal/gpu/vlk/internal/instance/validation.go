package instance

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

func validationLayers(opt CreateOptions) []string {
	if !opt.debugMode {
		return []string{}
	}

	layersCount := uint32(0)
	must.Work(vulkan.EnumerateInstanceLayerProperties(&layersCount, nil))

	availableLayers := make([]vulkan.LayerProperties, layersCount)
	must.Work(vulkan.EnumerateInstanceLayerProperties(&layersCount, availableLayers))

	foundLayers := make(map[string]any)
	for _, layer := range availableLayers {
		layer.Deref()
		foundLayers[vkconv.VarcharAsString(layer.LayerName)] = struct{}{}
	}

	notFound := make([]string, 0)
	found := make([]string, 0)

	for _, requiredLayer := range def.RequiredValidationLayers {
		layerName := vkconv.NormalizeString(requiredLayer)

		if _, exist := foundLayers[layerName]; !exist {
			notFound = append(notFound, layerName)
			continue
		}

		found = append(found, layerName)
	}

	opt.logger.Debug(fmt.Sprintf("available layers: [%v]", found))

	if len(notFound) > 0 {
		opt.logger.Error(
			fmt.Sprintf("debug may not work (turn off it in engine config), "+
				"because some of extensions not found: %v",
				notFound,
			),
		)
	}

	return found
}
