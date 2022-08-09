package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/command"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/physical"
)

type (
	internalBufferID uint32

	Allocator struct {
		logger config.Logger
		inst   *instance.Instance
		pd     *physical.Device
		ld     *logical.Device
		pool   *command.Pool

		internalBufferLastID internalBufferID
		allocatedBuffers     map[internalBufferID]internalBuffer
	}
)

func NewAllocator(
	logger config.Logger,
	inst *instance.Instance,
	pd *physical.Device,
	ld *logical.Device,
	pool *command.Pool,
) *Allocator {
	return &Allocator{
		logger: logger,
		inst:   inst,
		pd:     pd,
		ld:     ld,
		pool:   pool,

		internalBufferLastID: 0,
		allocatedBuffers:     make(map[internalBufferID]internalBuffer),
	}
}

func (a *Allocator) Free() {
	for _, buff := range a.allocatedBuffers {
		a.destroyBuffer(buff)
	}

	a.logger.Debug("freed: memory allocator")
}

func (a *Allocator) destroyBuffer(buff internalBuffer) {
	vulkan.DestroyBuffer(a.ld.Ref(), buff.ref, nil)
	vulkan.FreeMemory(a.ld.Ref(), buff.memory, nil)

	delete(a.allocatedBuffers, buff.id)
	a.logger.Debug(fmt.Sprintf("freed: buffer %d", buff.id))
}

func (a *Allocator) createVertexBuffer(size uint32) internalBuffer {
	return a.createBuffer(
		size,
		vulkan.BufferUsageVertexBufferBit,
		memRequirementHostVisible,
	)
}

func (a *Allocator) createIndexBuffer(size uint32) internalBuffer {
	return a.createBuffer(
		size,
		vulkan.BufferUsageIndexBufferBit,
		memRequirementHostVisible,
	)
}
