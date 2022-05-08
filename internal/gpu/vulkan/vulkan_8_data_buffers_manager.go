package vulkan

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/fe3dback/govgl/glm"
)

type (
	dataInstance interface {
		Data() []byte
		Indexes() []uint16
		VertexCount() uint32
	}

	bindStats struct {
		instanceCount uint32
		buffers       []vulkan.Buffer
		offsets       []vulkan.DeviceSize
	}

	flush struct {
		vertexChunks []bindStats
		indexBuffer  vulkan.Buffer
	}
)

func newDataBuffersManager(ld *vkLogicalDevice, pd *vkPhysicalDevice) *vkDataBuffersManager {
	return &vkDataBuffersManager{
		vertex: vkBufferTable{
			totalCapacity: 0,
			buffers:       []vkBuffer{},
		},
		index: vkBufferUnion{},
		ld:    ld,
		pd:    pd,
	}
}

func (vb *vkDataBuffersManager) free() {
	buffers := make([]vkBuffer, 0)
	buffers = append(buffers, vb.vertex.buffers...)
	buffers = append(buffers, vb.index.buffer)
	for _, union := range vb.uniform {
		buffers = append(buffers, union.buffer)
	}

	for _, vkBuffer := range buffers {
		vulkan.DestroyBuffer(vb.ld.ref, vkBuffer.handle, nil)
		vulkan.FreeMemory(vb.ld.ref, vkBuffer.memory, nil)
	}

	log.Printf("vk: freed: vertex buffers\n")
}

func (vb *vkDataBuffersManager) resetBuffers() {
	vb.resetIndexBuffer()
	vb.resetVertexBuffer()
}

func (vb *vkDataBuffersManager) resetIndexBuffer() {
	vb.index.staging = []byte{}
	vb.index.offset = 0
}

func (vb *vkDataBuffersManager) resetVertexBuffer() {
	vb.vertex.framePageID = -1
	vb.vertex.framePageCapacity = 0
	vb.vertex.frameInstanceCounts = []uint32{}
	vb.vertex.frameStagedData = [][]byte{}

	for range vb.vertex.buffers {
		vb.vertex.frameStagedData = append(vb.vertex.frameStagedData, []byte{})
		vb.vertex.frameInstanceCounts = append(vb.vertex.frameInstanceCounts, 0)
	}
}

func (vb *vkDataBuffersManager) flushBuffers() flush {
	return flush{
		vertexChunks: vb.flushVertexBuffer(),
		indexBuffer:  vb.flushIndexBuffer(),
	}
}

func (vb *vkDataBuffersManager) flushIndexBuffer() vulkan.Buffer {
	if vb.index.buffer.handle == nil {
		vb.allocateNewIndexBuffer()
	}

	vulkan.Memcopy(vb.index.buffer.dataPtr, vb.index.staging)
	return vb.index.buffer.handle
}

func (vb *vkDataBuffersManager) flushVertexBuffer() []bindStats {
	stats := make([]bindStats, 0)

	for pageID, stagedData := range vb.vertex.frameStagedData {
		buff := vb.vertex.buffers[pageID]
		vulkan.Memcopy(buff.dataPtr, stagedData)

		stats = append(stats, bindStats{
			instanceCount: vb.vertex.frameInstanceCounts[pageID],
			buffers:       []vulkan.Buffer{buff.handle},
			offsets:       []vulkan.DeviceSize{0},
		})
	}

	return stats
}

func (vb *vkDataBuffersManager) writeToBuffers(instance dataInstance) {
	vb.writeToIndexBuffer(instance)
	vb.writeToVertexBuffer(instance)
}

func (vb *vkDataBuffersManager) writeToIndexBuffer(instance dataInstance) {
	for _, index := range instance.Indexes() {
		index += vb.index.offset
		vb.index.staging = append(vb.index.staging, uint8(index&0xff), uint8(index>>8))
	}

	vb.index.offset += uint16(instance.VertexCount())
}

func (vb *vkDataBuffersManager) writeToVertexBuffer(instance dataInstance) {
	data := instance.Data()
	size := uint64(len(data))

	if vb.vertex.framePageCapacity < size {
		vb.vertex.framePageID++

		if int16(len(vb.vertex.buffers)-1) < vb.vertex.framePageID {
			vb.allocateNewVertexBuffer()
		}

		buff := vb.vertex.buffers[vb.vertex.framePageID]
		vb.vertex.framePageCapacity = uint64(buff.capacity)
	}

	// write to buffer
	vb.vertex.frameInstanceCounts[vb.vertex.framePageID]++
	stageData := &(vb.vertex.frameStagedData[vb.vertex.framePageID])
	*stageData = append(*stageData, data...)
	vb.vertex.framePageCapacity -= size
}

func (vb *vkDataBuffersManager) updateGlobalUniformBuffer(frameID uint32, projection, view, model glm.Mat4) {
	if int(frameID) > len(vb.uniform)-1 {
		vb.allocateNewUniformBuffer()
	}

	buf := make([]byte, 0, uniformBufferSize)

	// max memory supported here is 4*mat4
	buf = append(buf, projection.Data()...)
	buf = append(buf, view.Data()...)
	buf = append(buf, model.Data()...)

	vulkan.Memcopy(vb.uniform[frameID].buffer.dataPtr, buf)
}

func (vb *vkDataBuffersManager) allocateNewVertexBuffer() {
	buff := allocatePersistBuffer(vb.ld, vb.pd, vertexBufferSize, vulkan.BufferUsageVertexBufferBit)
	vb.vertex.buffers = append(vb.vertex.buffers, buff)
	vb.vertex.frameStagedData = append(vb.vertex.frameStagedData, []byte{})
	vb.vertex.frameInstanceCounts = append(vb.vertex.frameInstanceCounts, 0)
	vb.vertex.totalCapacity += uint64(buff.capacity)
}

func (vb *vkDataBuffersManager) allocateNewIndexBuffer() {
	vb.index.buffer = allocatePersistBuffer(vb.ld, vb.pd, indexBufferSize, vulkan.BufferUsageIndexBufferBit)
}

func (vb *vkDataBuffersManager) allocateNewUniformBuffer() {
	vb.uniform = append(vb.uniform, vkBufferUnion{
		offset:  0,
		buffer:  allocatePersistBuffer(vb.ld, vb.pd, uniformBufferSize, vulkan.BufferUsageUniformBufferBit),
		staging: []byte{},
	})
}

func allocatePersistBuffer(ld *vkLogicalDevice, pd *vkPhysicalDevice, size int, buffType vulkan.BufferUsageFlagBits) vkBuffer {
	// create new buffer page
	info := &vulkan.BufferCreateInfo{
		SType:       vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(size),
		Usage:       vulkan.BufferUsageFlags(buffType),
		SharingMode: vulkan.SharingModeExclusive,
	}

	var buffer vulkan.Buffer
	vkAssert(
		vulkan.CreateBuffer(ld.ref, info, nil, &buffer),
		fmt.Errorf("failed create vertex buffer"),
	)

	// get device memory requirements for it
	var memoryReq vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(ld.ref, buffer, &memoryReq)
	memoryReq.Deref()

	memoryTypeIndex := findVertexBufferMemoryType(
		pd,
		memoryReq,
		vulkan.MemoryPropertyFlags(
			vulkan.MemoryPropertyHostVisibleBit|
				vulkan.MemoryPropertyHostCoherentBit,
		),
	)

	memAllocInfo := &vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memoryReq.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	var bufferMemory vulkan.DeviceMemory
	vkAssert(
		vulkan.AllocateMemory(ld.ref, memAllocInfo, nil, &bufferMemory),
		fmt.Errorf("failed allocate GPU memory for vertex buffer"),
	)

	vulkan.BindBufferMemory(ld.ref, buffer, bufferMemory, 0)

	var data unsafe.Pointer
	vulkan.MapMemory(ld.ref, bufferMemory, 0, info.Size, 0, &data)

	log.Printf("Buffer %.3fMB capacity - allocated", float64(info.Size/1024))

	return vkBuffer{
		capacity: info.Size,
		dataPtr:  data,
		handle:   buffer,
		memory:   bufferMemory,
	}
}

func findVertexBufferMemoryType(pd *vkPhysicalDevice, memoryReq vulkan.MemoryRequirements, memFlags vulkan.MemoryPropertyFlags) uint32 {
	typeFilter := memoryReq.MemoryTypeBits

	var memProperties vulkan.PhysicalDeviceMemoryProperties
	vulkan.GetPhysicalDeviceMemoryProperties(pd.ref, &memProperties)
	memProperties.Deref()

	for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
		memType := memProperties.MemoryTypes[i]
		memType.Deref()

		if (typeFilter&(1<<i) != 0) && ((memType.PropertyFlags & memFlags) == memFlags) {
			return i
		}
	}

	panic(fmt.Errorf("failed find suitable GPU memory for vertex buffer"))
}
