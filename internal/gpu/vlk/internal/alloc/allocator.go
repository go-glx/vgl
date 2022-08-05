package alloc

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type (
	Allocator struct {
		logger config.Logger
		inst   *instance.Instance
		pd     *physical.Device
		ld     *logical.Device

		allocatedBuffers []internalBuffer
	}
)

func NewAllocator(logger config.Logger, inst *instance.Instance, pd *physical.Device, ld *logical.Device) *Allocator {
	return &Allocator{
		logger:           logger,
		inst:             inst,
		pd:               pd,
		ld:               ld,
		allocatedBuffers: make([]internalBuffer, 0),
	}
}

func (a *Allocator) Free() {
	for _, buff := range a.allocatedBuffers {
		vulkan.DestroyBuffer(a.ld.Ref(), buff.ref, nil)
		vulkan.FreeMemory(a.ld.Ref(), buff.memory, nil)
	}

	a.logger.Debug("freed: memory allocator")
}

func (a *Allocator) createVertexBuffer() internalBuffer {
	return a.createBuffer(def.BufferVertexSizeBytes, vulkan.BufferUsageFlags(
		vulkan.BufferUsageVertexBufferBit,
	))
}

func (a *Allocator) createIndexBuffer() internalBuffer {
	return a.createBuffer(def.BufferVertexSizeBytes, vulkan.BufferUsageFlags(
		vulkan.BufferUsageIndexBufferBit,
	))
}
