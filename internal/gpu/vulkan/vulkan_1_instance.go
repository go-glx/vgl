package vulkan

import (
	"fmt"
	"log"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/arch"
)

var requiredValidationLayers = []string{
	"VK_LAYER_KHRONOS_validation",
}

func newVkInstance(wm arch.WindowManager, requiredExt []string, debugMode bool) *vkInstance {
	var inst vulkan.Instance

	vkAssert(
		vulkan.CreateInstance(instanceCreateInfo(wm.AppName(), wm.EngineName(), requiredExt, debugMode), nil, &inst),
		fmt.Errorf("create vulkan instance failed"),
	)

	return &vkInstance{
		ref: inst,
	}
}

func (inst *vkInstance) free() {
	vulkan.DestroyInstance(inst.ref, nil)

	log.Printf("vk: freed: vulkan instance\n")
}

func instanceCreateInfo(windowTitle, engineTitle string, requiredExt []string, debugMode bool) *vulkan.InstanceCreateInfo {
	log.Printf("vk: init '%s', required extensions: [%v]\n", engineTitle, requiredExt)

	instInfo := &vulkan.InstanceCreateInfo{
		SType: vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vulkan.ApplicationInfo{
			SType:              vulkan.StructureTypeApplicationInfo,
			PApplicationName:   windowTitle,
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			PEngineName:        engineTitle,
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			ApiVersion:         vulkan.ApiVersion11,
		},
	}

	// setup extensions
	availableExt := enumerateAvailableExtensions()
	ensureAllRequiredExtensionsIsAvailable(availableExt, requiredExt)
	instInfo.PpEnabledExtensionNames = requiredExt
	instInfo.EnabledExtensionCount = uint32(len(instInfo.PpEnabledExtensionNames))

	// setup validation (debug)
	validationLayers := validationLayers(debugMode)
	instInfo.EnabledLayerCount = uint32(len(validationLayers))
	instInfo.PpEnabledLayerNames = validationLayers

	return instInfo
}

func ensureAllRequiredExtensionsIsAvailable(available map[string]struct{}, required []string) {
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

func enumerateAvailableExtensions() map[string]struct{} {
	var extCount uint32
	vkAssert(
		vulkan.EnumerateInstanceExtensionProperties("", &extCount, nil),
		fmt.Errorf("failed enumerate extensions count"),
	)

	if extCount <= 0 {
		return map[string]struct{}{}
	}

	extensions := make([]vulkan.ExtensionProperties, extCount)
	vkAssert(
		vulkan.EnumerateInstanceExtensionProperties("", &extCount, extensions),
		fmt.Errorf("failed enumerate extensions"),
	)

	availableExt := make(map[string]struct{})
	for _, extension := range extensions {
		extension.Deref()

		extLabel := vkLabelToString(extension.ExtensionName)
		if extLabel == "" {
			continue
		}

		fmt.Printf("vk: available ext: %s (v%d)\n", extLabel, extension.SpecVersion)
		availableExt[extLabel] = struct{}{}
	}

	return availableExt
}

func validationLayers(isDebugMode bool) []string {
	if !isDebugMode {
		return []string{}
	}

	layersCount := uint32(0)
	vkAssert(
		vulkan.EnumerateInstanceLayerProperties(&layersCount, nil),
		fmt.Errorf("failed enumerate layer properties"),
	)

	availableLayers := make([]vulkan.LayerProperties, layersCount)
	vkAssert(
		vulkan.EnumerateInstanceLayerProperties(&layersCount, availableLayers),
		fmt.Errorf("failed enumerate layer properties"),
	)

	foundLayers := make(map[string]struct{})
	for _, layer := range availableLayers {
		layer.Deref()
		foundLayers[vkLabelToString(layer.LayerName)] = struct{}{}
	}

	notFound := make([]string, 0)
	found := make([]string, 0)

	for _, requiredLayer := range requiredValidationLayers {
		layerLabel := vkRepackLabel(requiredLayer)
		if _, exist := foundLayers[layerLabel]; !exist {
			notFound = append(notFound, layerLabel)
			continue
		}

		found = append(found, layerLabel)
	}

	log.Printf("vk: available layers: [%v]\n", found)

	if len(notFound) > 0 {
		log.Printf("vk: debug may not work (turn off it in game config), because some of extensions not found: %v\n",
			notFound,
		)
	}

	return found
}
