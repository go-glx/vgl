package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
)

type pageFeatures struct {
	bufferType    BufferType
	storageTarget StorageTarget
	flags         Flags
}

func (pf *pageFeatures) defaultBufferAlign() uint32 {
	// currently 256 bytes for any type of buffers
	return 256
}

func (pf *pageFeatures) defaultBufferCapacity() uint32 {
	switch pf.bufferType {
	case BufferTypeIndex:
		return def.BufferIndexSizeBytes
	case BufferTypeVertex:
		return def.BufferVertexSizeBytes
	case BufferTypeUniform:
		return def.BufferUniformSizeBytes
	default:
		return 1 * 1024 * 512 // 512KB
	}
}

func (pf *pageFeatures) vulkanBufferUsage() vulkan.BufferUsageFlagBits {
	// immutable (device only memory) has transferDst because data
	// uploading will work through copy bytes from src buffer

	switch pf.storageTarget {
	case StorageTargetImmutable:
		return vulkan.BufferUsageTransferDstBit | pf.vulkanUsageType()
	case StorageTargetWritable:
		return vulkan.BufferUsageTransferDstBit | pf.vulkanUsageType()
	case StorageTargetCoherent:
		return pf.vulkanUsageType()
	default:
		panic(fmt.Errorf("unknown storage target type %d", pf.storageTarget))
	}
}

func (pf *pageFeatures) vulkanMemoryFlags() vulkan.MemoryPropertyFlagBits {
	switch pf.storageTarget {
	case StorageTargetImmutable:
		return vulkan.MemoryPropertyDeviceLocalBit
	case StorageTargetWritable:
		return vulkan.MemoryPropertyDeviceLocalBit
	case StorageTargetCoherent:
		return vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit
	default:
		panic(fmt.Errorf("unknown storage target type %d", pf.storageTarget))
	}
}

func (pf *pageFeatures) vulkanUsageType() vulkan.BufferUsageFlagBits {
	switch pf.bufferType {
	case BufferTypeIndex:
		return vulkan.BufferUsageIndexBufferBit
	case BufferTypeVertex:
		return vulkan.BufferUsageVertexBufferBit
	case BufferTypeUniform:
		return vulkan.BufferUsageUniformBufferBit
	default:
		panic(fmt.Errorf("unknown buffer type %d", pf.bufferType))
	}
}
