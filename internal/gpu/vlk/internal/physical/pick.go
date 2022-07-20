package physical

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/vkconv"
)

func (d *Device) pickPrimaryGPU() *GPU {
	bestScore := -1
	var bestDevice *GPU

	for _, pd := range d.listDevices() {
		score := d.score(pd)
		if score < 0 {
			// device is not suitable at all
			continue
		}

		log.Printf("vk: GPU '%s' is suitable for use, score = %d\n",
			vkconv.VarcharAsString(pd.Props.DeviceName),
			score,
		)

		if score > bestScore {
			bestScore = score
			bestDevice = pd
		}
	}

	if bestDevice == nil {
		panic(fmt.Errorf("not found suitable vulkan GPU for rendering"))
	}

	log.Printf("vk: using GPU: %s\n", vkconv.VarcharAsString(bestDevice.Props.DeviceName))
	return bestDevice
}

func (d *Device) listDevices() []*GPU {
	count := uint32(0)
	must.Work(vulkan.EnumeratePhysicalDevices(d.inst.Ref(), &count, nil))
	if count <= 0 {
		return nil
	}

	physicalDevices := make([]vulkan.PhysicalDevice, count)
	must.Work(vulkan.EnumeratePhysicalDevices(d.inst.Ref(), &count, physicalDevices))

	result := make([]*GPU, 0, len(physicalDevices))
	for _, pd := range physicalDevices {
		result = append(result, d.assembleGPU(pd))
	}

	return result
}
