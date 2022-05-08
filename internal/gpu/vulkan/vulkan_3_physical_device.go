package vulkan

import (
	"fmt"
	"log"
	"strings"

	"github.com/vulkan-go/vulkan"
)

var requiredDeviceExtensions = []string{
	"VK_KHR_swapchain", // require for display buffer to screen
}

func newPhysicalDeviceFinder(inst *vkInstance, surface *vkSurface) *vkPhysicalDeviceFinder {
	return &vkPhysicalDeviceFinder{
		inst:    inst,
		surface: surface,
	}
}

func (f *vkPhysicalDeviceFinder) physicalDevicePick() *vkPhysicalDevice {
	bestScore := -1
	var bestDevice *vkPhysicalDevice

	for _, pd := range f.physicalDevices() {
		score := f.physicalDeviceScore(pd)
		if score < 0 {
			// device is not suitable at all
			continue
		}

		log.Printf("vk: GPU '%s' is suitable for use, score = %d\n", vkLabelToString(pd.props.DeviceName), score)
		if score > bestScore {
			bestScore = score
			bestDevice = pd
		}
	}

	if bestDevice == nil {
		panic(fmt.Errorf("not found suitable vulkan GPU for rendering"))
	}

	log.Printf("vk: using GPU: %s\n", vkLabelToString(bestDevice.props.DeviceName))
	return bestDevice
}

func (f *vkPhysicalDeviceFinder) physicalDevices() []*vkPhysicalDevice {
	count := uint32(0)
	vkAssert(
		vulkan.EnumeratePhysicalDevices(f.inst.ref, &count, nil),
		fmt.Errorf("failed EnumeratePhysicalDevices"),
	)
	if count <= 0 {
		panic(fmt.Errorf("not found any GPU with vulkan support"))
	}

	physicalDevices := make([]vulkan.PhysicalDevice, count)
	vkAssert(
		vulkan.EnumeratePhysicalDevices(f.inst.ref, &count, physicalDevices),
		fmt.Errorf("failed EnumeratePhysicalDevices"),
	)

	result := make([]*vkPhysicalDevice, 0, len(physicalDevices))
	for _, pd := range physicalDevices {
		result = append(result, f.physicalDeviceAssemble(pd))
	}

	return result
}

func (f *vkPhysicalDeviceFinder) physicalDeviceAssemble(pd vulkan.PhysicalDevice) *vkPhysicalDevice {
	var props vulkan.PhysicalDeviceProperties
	vulkan.GetPhysicalDeviceProperties(pd, &props)
	props.Deref()

	var features vulkan.PhysicalDeviceFeatures
	vulkan.GetPhysicalDeviceFeatures(pd, &features)
	features.Deref()

	return &vkPhysicalDevice{
		ref:          pd,
		props:        props,
		features:     features,
		families:     f.physicalDeviceFamilies(pd),
		extensions:   f.physicalDeviceExtensions(pd),
		surfaceProps: f.physicalDeviceSurfaceProps(pd),
	}
}

// supported values:
// any value < 0 	- device is not suitable for use
// any value >= 0 	- device score
func (f *vkPhysicalDeviceFinder) physicalDeviceScore(pd *vkPhysicalDevice) int {
	required := map[bool]string{
		// families
		!pd.families.supportGraphics: "graphics operations not supported",
		!pd.families.supportPresent:  "window present not supported",

		// extensions
		!pd.isSupportAllRequiredExtensions(): "not all required extensions supported",

		// swap chain
		len(pd.surfaceProps.formats) <= 0:             "not GPU",
		len(pd.surfaceProps.presentModes) <= 0:        "not GPU",
		pd.surfaceProps.richColorSpaceFormat() == nil: "rich colorspace not supported",
	}

	// filter
	for failed, reason := range required {
		if failed {
			log.Printf("vk: GPU '%s' not pass check: %s\n", vkLabelToString(pd.props.DeviceName), reason)
			return -1
		}
	}

	// score
	score := 0
	if pd.props.DeviceType == vulkan.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	return score
}

func (pd *vkPhysicalDevice) isSupportAllRequiredExtensions() bool {
	supportedExt := make(map[string]struct{})
	for _, extension := range pd.extensions {
		supportedExt[vkLabelToString(extension.ExtensionName)] = struct{}{}
	}

	notSupported := make([]string, 0)
	for _, extension := range requiredDeviceExtensions {
		if _, supported := supportedExt[vkRepackLabel(extension)]; supported {
			continue
		}

		notSupported = append(notSupported, extension)
	}

	if len(notSupported) > 0 {
		log.Printf("vk: GPU '%s' not support all required extensions: [%s]\n", vkLabelToString(pd.props.DeviceName), strings.Join(notSupported, ", "))
		return false
	}

	return true
}
