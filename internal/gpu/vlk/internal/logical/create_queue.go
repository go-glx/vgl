package logical

import "github.com/vulkan-go/vulkan"

func (dev *Device) createQueuesInfo() []vulkan.DeviceQueueCreateInfo {
	infos := make([]vulkan.DeviceQueueCreateInfo, 0)

	for _, familyId := range dev.pd.PrimaryGPU().Families.UniqueIDs() {
		infos = append(infos, vulkan.DeviceQueueCreateInfo{
			SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: familyId,
			QueueCount:       1,
			PQueuePriorities: []float32{1.0},
		})
	}

	return infos
}
