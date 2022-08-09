package alloc

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

func (a *Allocator) copyBuffer(src internalBuffer, dst internalBuffer, srcOffset, dstOffset, size uint32) {
	a.pool.TemporaryBuffer(func(cb vulkan.CommandBuffer) {
		copyRegion := vulkan.BufferCopy{
			SrcOffset: vulkan.DeviceSize(srcOffset),
			DstOffset: vulkan.DeviceSize(dstOffset),
			Size:      vulkan.DeviceSize(size),
		}

		vulkan.CmdCopyBuffer(cb, src.ref, dst.ref, 1, []vulkan.BufferCopy{copyRegion})

		a.logger.Debug(fmt.Sprintf("buffer data copied (%d->%d) offsets=[src=%d, dst=%d], size=%.2fKB",
			src.id,
			dst.id,
			srcOffset,
			dstOffset,
			float32(size)/1024,
		))
	})
}
