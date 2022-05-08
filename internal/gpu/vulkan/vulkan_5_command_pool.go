package vulkan

import (
	"fmt"
	"log"

	"github.com/vulkan-go/vulkan"
)

func newCommandPool(pd *vkPhysicalDevice, ld *vkLogicalDevice) *vkCommandPool {
	commandPool := &vkCommandPool{
		ld: ld,
	}

	// pool
	poolInfo := &vulkan.CommandPoolCreateInfo{
		SType:            vulkan.StructureTypeCommandPoolCreateInfo,
		QueueFamilyIndex: pd.families.graphicsFamilyId,
	}

	var pool vulkan.CommandPool
	vkAssert(
		vulkan.CreateCommandPool(ld.ref, poolInfo, nil, &pool),
		fmt.Errorf("failed create command pool"),
	)
	commandPool.ref = pool

	// command buffers
	buffersCount := pd.surfaceProps.imageBuffersCount()
	allocInfo := &vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        pool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: buffersCount,
	}

	commandBuffers := make([]vulkan.CommandBuffer, buffersCount)
	vkAssert(
		vulkan.AllocateCommandBuffers(ld.ref, allocInfo, commandBuffers),
		fmt.Errorf("failed allocate command buffers"),
	)

	commandPool.buffers = commandBuffers
	return commandPool
}

func (pool *vkCommandPool) free() {
	vulkan.FreeCommandBuffers(pool.ld.ref, pool.ref, uint32(len(pool.buffers)), pool.buffers)
	vulkan.DestroyCommandPool(pool.ld.ref, pool.ref, nil)

	log.Printf("vk: freed: command pool\n")
}

func (pool *vkCommandPool) commandBufferStart(ind int) {
	vkAssert(
		vulkan.BeginCommandBuffer(pool.buffers[ind], &vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
		}),
		fmt.Errorf("failed begin command buffer"),
	)
}

func (pool *vkCommandPool) commandBufferEnd(ind int) {
	vkAssert(
		vulkan.EndCommandBuffer(pool.buffers[ind]),
		fmt.Errorf("failed end command buffer"),
	)
}

func (pool *vkCommandPool) commandBuffer(ind int) vulkan.CommandBuffer {
	return pool.buffers[ind]
}
