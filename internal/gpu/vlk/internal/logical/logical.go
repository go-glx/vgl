package logical

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

type Device struct {
	pd *physical.Device

	ref           vulkan.Device
	queueGraphics vulkan.Queue
	queuePresent  vulkan.Queue
}

func NewDevice(pd *physical.Device) *Device {
	dev := &Device{
		pd: pd,
	}
	dev.createLogicalAndEnrich()

	return dev
}

func (dev *Device) Ref() vulkan.Device {
	return dev.ref
}

func (dev *Device) Free() {
	vulkan.DestroyDevice(dev.ref, nil)
	log.Printf("vk: freed: logical device\n")
}

func (dev *Device) createLogicalAndEnrich() {
	// prepare data
	gpu := dev.pd.PrimaryGPU()
	queues := dev.createQueuesInfo()

	createInfo := &vulkan.DeviceCreateInfo{
		SType:                   vulkan.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    uint32(len(queues)),
		PQueueCreateInfos:       queues,
		PEnabledFeatures:        []vulkan.PhysicalDeviceFeatures{gpu.Features},
		EnabledExtensionCount:   uint32(len(gpu.RequiredExtensions)),
		PpEnabledExtensionNames: vkconv.NormalizeStringList(gpu.RequiredExtensions),
	}

	// create device
	var logicalDevice vulkan.Device
	must.Work(vulkan.CreateDevice(gpu.Ref, createInfo, nil, &logicalDevice))

	// create queues
	var queueGraphics vulkan.Queue
	var queuePresent vulkan.Queue
	vulkan.GetDeviceQueue(logicalDevice, gpu.Families.GraphicsFamilyId, 0, &queueGraphics)
	vulkan.GetDeviceQueue(logicalDevice, gpu.Families.PresentFamilyId, 0, &queuePresent)

	// log
	log.Printf("vk: logical device created (graphicsQ: %d, presentQ: %d)\n",
		gpu.Families.GraphicsFamilyId,
		gpu.Families.PresentFamilyId,
	)

	// enrich
	dev.ref = logicalDevice
	dev.queueGraphics = queueGraphics
	dev.queuePresent = queuePresent
}
