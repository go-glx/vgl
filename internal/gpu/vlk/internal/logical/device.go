package logical

import (
	"fmt"
	"strings"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
	"github.com/go-glx/vgl/shared/vlkext"
)

type Device struct {
	logger vlkext.Logger
	pd     *physical.Device

	ref           vulkan.Device
	queueGraphics vulkan.Queue
	queuePresent  vulkan.Queue
}

func NewDevice(logger vlkext.Logger, pd *physical.Device) *Device {
	dev := &Device{
		logger: logger,
		pd:     pd,
	}
	dev.createLogicalAndEnrich()

	return dev
}

func (dev *Device) Ref() vulkan.Device {
	return dev.ref
}

func (dev *Device) QueueGraphics() vulkan.Queue {
	return dev.queueGraphics
}

func (dev *Device) QueuePresent() vulkan.Queue {
	return dev.queuePresent
}

func (dev *Device) Free() {
	vulkan.DestroyDevice(dev.ref, nil)
	dev.logger.Debug("freed: logical device")
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

	dev.logger.Debug(fmt.Sprintf("gpu require ext: [%s]", strings.Join(gpu.RequiredExtensions, ", ")))

	// create device
	var logicalDevice vulkan.Device
	must.Work(vulkan.CreateDevice(gpu.Ref, createInfo, nil, &logicalDevice))

	// create queues
	var queueGraphics vulkan.Queue
	var queuePresent vulkan.Queue
	vulkan.GetDeviceQueue(logicalDevice, gpu.Families.GraphicsFamilyId, 0, &queueGraphics)
	vulkan.GetDeviceQueue(logicalDevice, gpu.Families.PresentFamilyId, 0, &queuePresent)

	// log
	dev.logger.Debug(fmt.Sprintf("logical device created (graphicsQ: %d, presentQ: %d)",
		gpu.Families.GraphicsFamilyId,
		gpu.Families.PresentFamilyId,
	))

	// enrich
	dev.ref = logicalDevice
	dev.queueGraphics = queueGraphics
	dev.queuePresent = queuePresent
}
