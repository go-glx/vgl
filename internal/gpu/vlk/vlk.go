package vlk

import (
	"github.com/vulkan-go/vulkan"
)

type VLK struct {
	isReady bool
	cont    *Container

	// drawing
	currentBatch *drawCall
	queue        []drawCall
}

func newVLK(cont *Container) *VLK {
	return &VLK{
		isReady: true,
		cont:    cont,
	}
}

// this will immediately stop render new frames
// wait for all current GPU work is done, then
// run mutate function, that allow change any VLK state
// after this, VLK will be turned on again
func (vlk *VLK) maintenance(mutate func()) {
	// stop render
	vlk.isReady = false

	// wait for GPU end current operations
	vulkan.DeviceWaitIdle(vlk.cont.logicalDevice().Ref())

	// change vulkan state
	// rebuild pipeline, etc..
	mutate()

	// turn on back
	vlk.isReady = true
}
