package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

func (f *vkPhysicalDeviceFinder) physicalDeviceExtensions(pd vulkan.PhysicalDevice) []vulkan.ExtensionProperties {
	count := uint32(0)
	vkAssert(
		vulkan.EnumerateDeviceExtensionProperties(pd, "", &count, nil),
		fmt.Errorf("failed enumerate device extensions"),
	)
	if count == 0 {
		return nil
	}

	ext := make([]vulkan.ExtensionProperties, count)
	vkAssert(
		vulkan.EnumerateDeviceExtensionProperties(pd, "", &count, ext),
		fmt.Errorf("failed enumerate device extensions"),
	)

	result := make([]vulkan.ExtensionProperties, 0, count)
	for _, properties := range ext {
		properties.Deref()
		result = append(result, properties)
	}

	return result
}

func (f *vkPhysicalDeviceFinder) physicalDeviceSurfaceProps(pd vulkan.PhysicalDevice) vkPhysicalDeviceSurfaceProps {
	return vkPhysicalDeviceSurfaceProps{
		capabilities: f.physicalDeviceGetSurfaceCapabilities(pd),
		formats:      f.physicalDeviceGetSurfaceFormats(pd),
		presentModes: f.physicalDevicePresentModes(pd),
	}
}

func (f *vkPhysicalDeviceFinder) physicalDeviceGetSurfaceCapabilities(device vulkan.PhysicalDevice) vulkan.SurfaceCapabilities {
	var capabilities vulkan.SurfaceCapabilities
	vkAssert(
		vulkan.GetPhysicalDeviceSurfaceCapabilities(device, f.surface.ref, &capabilities),
		fmt.Errorf("failed get physical device capabilities"),
	)

	capabilities.Deref()
	return capabilities
}

func (f *vkPhysicalDeviceFinder) physicalDeviceGetSurfaceFormats(device vulkan.PhysicalDevice) []vulkan.SurfaceFormat {
	formatsCount := uint32(0)
	vkAssert(
		vulkan.GetPhysicalDeviceSurfaceFormats(device, f.surface.ref, &formatsCount, nil),
		fmt.Errorf("failed get physical device surface formats"),
	)

	surfaceFormats := make([]vulkan.SurfaceFormat, formatsCount)
	vkAssert(
		vulkan.GetPhysicalDeviceSurfaceFormats(device, f.surface.ref, &formatsCount, surfaceFormats),
		fmt.Errorf("failed get physical device surface formats"),
	)

	resultFormats := make([]vulkan.SurfaceFormat, 0, len(surfaceFormats))
	for _, format := range surfaceFormats {
		format.Deref()
		resultFormats = append(resultFormats, format)
	}

	return resultFormats
}

func (f *vkPhysicalDeviceFinder) physicalDevicePresentModes(device vulkan.PhysicalDevice) []vulkan.PresentMode {
	modesCount := uint32(0)
	vkAssert(
		vulkan.GetPhysicalDeviceSurfacePresentModes(device, f.surface.ref, &modesCount, nil),
		fmt.Errorf("failed get physical device present modes"),
	)

	presentModes := make([]vulkan.PresentMode, modesCount)
	vkAssert(
		vulkan.GetPhysicalDeviceSurfacePresentModes(device, f.surface.ref, &modesCount, presentModes),
		fmt.Errorf("failed get physical device present modes"),
	)

	return presentModes
}

func (f *vkPhysicalDeviceFinder) physicalDeviceFamilies(device vulkan.PhysicalDevice) vkPhysicalDeviceFamilies {
	count := uint32(0)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &count, nil)
	if count == 0 {
		return vkPhysicalDeviceFamilies{}
	}

	families := make([]vulkan.QueueFamilyProperties, count)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &count, families)

	result := vkPhysicalDeviceFamilies{}

	for familyId, properties := range families {
		properties.Deref()
		if properties.QueueFlags&vulkan.QueueFlags(vulkan.QueueGraphicsBit) != 0 {
			result.graphicsFamilyId = uint32(familyId)
			result.supportGraphics = true
		}

		var presentSupport vulkan.Bool32
		vkAssert(
			vulkan.GetPhysicalDeviceSurfaceSupport(device, uint32(familyId), f.surface.ref, &presentSupport),
			fmt.Errorf("failed check device surface support"),
		)
		if presentSupport != 0 {
			result.presentFamilyId = uint32(familyId)
			result.supportPresent = true
		}
	}

	return result
}
