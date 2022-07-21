package physical

import (
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/surface"
)

type Device struct {
	inst    *instance.Instance
	surface *surface.Surface

	primaryGPU *GPU
}

func NewDevice(inst *instance.Instance, surface *surface.Surface) *Device {
	dev := &Device{inst: inst, surface: surface}
	dev.primaryGPU = dev.pickPrimaryGPU()

	return dev
}

func (d *Device) PrimaryGPU() *GPU {
	return d.primaryGPU
}
