package instance

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

func validationLayers(isDebugMode bool) []string {
	if !isDebugMode {
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

	log.Printf("vk: available layers: [%v]\n", found)

	if len(notFound) > 0 {
		log.Printf("vk: debug may not work (turn off it in engine config), because some of extensions not found: %v\n",
			notFound,
		)
	}

	return found
}
