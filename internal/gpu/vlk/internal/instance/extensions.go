package instance

import (
	"fmt"
	"log"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type (
	extList map[string]any
)

func fetchAvailableExtensions() extList {
	var extCount uint32
	must.Work(vulkan.EnumerateInstanceExtensionProperties("", &extCount, nil))

	if extCount <= 0 {
		return make(extList)
	}

	extensions := make([]vulkan.ExtensionProperties, extCount)
	must.Work(vulkan.EnumerateInstanceExtensionProperties("", &extCount, extensions))

	availableExt := make(extList)
	for _, extension := range extensions {
		extension.Deref()

		extName := vkconv.VarcharAsString(extension.ExtensionName)
		if extName == "" {
			continue
		}

		log.Printf("vk: available ext: %s (v%d)\n", extName, extension.SpecVersion)
		availableExt[extName] = struct{}{}
	}

	return availableExt
}

func assertRequiredExtensionsIsAvailable(available extList, required []string) {
	notAvailable := make([]string, 0)

	for _, ext := range required {
		if _, exist := available[ext]; exist {
			continue
		}

		notAvailable = append(notAvailable, ext)
	}

	if len(notAvailable) == 0 {
		return
	}

	panic(fmt.Errorf("vk: required extensions [%s] not available",
		strings.Join(notAvailable, ", "),
	))
}
