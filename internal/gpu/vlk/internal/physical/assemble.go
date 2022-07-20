package physical

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

func (d *Device) assembleGPU(pd vulkan.PhysicalDevice) *GPU {
	var props vulkan.PhysicalDeviceProperties
	vulkan.GetPhysicalDeviceProperties(pd, &props)
	props.Deref()

	var features vulkan.PhysicalDeviceFeatures
	vulkan.GetPhysicalDeviceFeatures(pd, &features)
	features.Deref()

	vkExtList := make([]string, 0, len(requiredDeviceExtensions))
	for _, extName := range requiredDeviceExtensions {
		vkExtList = append(vkExtList, vkconv.NormalizeString(extName))
	}

	return &GPU{
		Ref:                pd,
		Props:              props,
		Features:           features,
		Families:           d.assembleFamilies(pd),
		Extensions:         d.assembleExtensions(pd),
		SurfaceProps:       d.assembleSurfaceProps(pd),
		RequiredExtensions: vkExtList,
	}
}

func (d *Device) assembleFamilies(device vulkan.PhysicalDevice) Families {
	count := uint32(0)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &count, nil)
	if count == 0 {
		// can return empty families here
		// because supportGraphics and supportPresent will be false inside
		// and these families not pass next validation
		return Families{}
	}

	families := make([]vulkan.QueueFamilyProperties, count)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &count, families)

	result := Families{}

	for familyId, properties := range families {
		properties.Deref()
		if properties.QueueFlags&vulkan.QueueFlags(vulkan.QueueGraphicsBit) != 0 {
			result.GraphicsFamilyId = uint32(familyId)
			result.supportGraphics = true
		}

		var presentSupport vulkan.Bool32
		must.Work(vulkan.GetPhysicalDeviceSurfaceSupport(device, uint32(familyId), d.surface.Ref(), &presentSupport))

		if presentSupport != 0 {
			result.PresentFamilyId = uint32(familyId)
			result.supportPresent = true
		}
	}

	return result
}

func (d *Device) assembleExtensions(pd vulkan.PhysicalDevice) []vulkan.ExtensionProperties {
	count := uint32(0)
	must.Work(vulkan.EnumerateDeviceExtensionProperties(pd, "", &count, nil))
	if count == 0 {
		return nil
	}

	propList := make([]vulkan.ExtensionProperties, count)
	must.Work(vulkan.EnumerateDeviceExtensionProperties(pd, "", &count, propList))

	result := make([]vulkan.ExtensionProperties, 0, count)
	for _, properties := range propList {
		properties.Deref()
		result = append(result, properties)
	}

	return result
}

func (d *Device) assembleSurfaceProps(pd vulkan.PhysicalDevice) SurfaceProps {
	return SurfaceProps{
		capabilities: d.assembleSurfacePropsCapabilities(pd),
		formats:      d.assembleSurfacePropsFormats(pd),
		presentModes: d.assembleSurfacePropsPresentModes(pd),
	}
}

func (d *Device) assembleSurfacePropsCapabilities(device vulkan.PhysicalDevice) vulkan.SurfaceCapabilities {
	var capabilities vulkan.SurfaceCapabilities
	must.Work(vulkan.GetPhysicalDeviceSurfaceCapabilities(device, d.surface.Ref(), &capabilities))

	capabilities.Deref()
	return capabilities
}

func (d *Device) assembleSurfacePropsFormats(device vulkan.PhysicalDevice) []vulkan.SurfaceFormat {
	formatsCount := uint32(0)
	must.Work(vulkan.GetPhysicalDeviceSurfaceFormats(device, d.surface.Ref(), &formatsCount, nil))

	surfaceFormats := make([]vulkan.SurfaceFormat, formatsCount)
	must.Work(vulkan.GetPhysicalDeviceSurfaceFormats(device, d.surface.Ref(), &formatsCount, surfaceFormats))

	resultFormats := make([]vulkan.SurfaceFormat, 0, len(surfaceFormats))
	for _, format := range surfaceFormats {
		format.Deref()
		resultFormats = append(resultFormats, format)
	}

	return resultFormats
}

func (d *Device) assembleSurfacePropsPresentModes(device vulkan.PhysicalDevice) []vulkan.PresentMode {
	modesCount := uint32(0)
	must.Work(vulkan.GetPhysicalDeviceSurfacePresentModes(device, d.surface.Ref(), &modesCount, nil))

	presentModes := make([]vulkan.PresentMode, modesCount)
	must.Work(vulkan.GetPhysicalDeviceSurfacePresentModes(device, d.surface.Ref(), &modesCount, presentModes))

	return presentModes
}
