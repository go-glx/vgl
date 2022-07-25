package vkconv

import (
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func TransformByteCode(data []byte) []uint32 {
	buf := make([]uint32, len(data)/4)
	vulkan.Memcopy(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buf)).Data), data)
	return buf
}
