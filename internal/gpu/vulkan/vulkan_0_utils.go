package vulkan

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

func vkMapToList(src map[string]struct{}) []string {
	r := make([]string, 0, len(src))

	for s := range src {
		r = append(r, s)
	}

	return r
}

// string(varchar) -> [256]byte -> string(256)
func vkRepackLabel(src string) string {
	return vkLabelToString(vkStringToLabel(src))
}

func vkStringsToStringLabels(src []string) []string {
	labels := make([]string, 0, len(src))
	for _, s := range src {
		labels = append(labels, vkLabelToString(vkStringToLabel(s)))
	}
	return labels
}

func vkLabelToString(src [256]byte) string {
	r := strings.Builder{}

	for _, b := range src {
		if b == byte(0) { // skip null bytes
			continue
		}

		r.WriteByte(b)
	}

	return strings.TrimSpace(r.String())
}

func vkStringToLabel(src string) [256]byte {
	if len(src) > 256 {
		panic(fmt.Errorf("failed convert string '%s' to vulkan label: len is greater to 256", src))
	}

	label := [256]byte{}

	for i, b := range src {
		label[i] = byte(b)
	}

	return label
}

func vkClampUint(n, min, max uint32) uint32 {
	if n <= min {
		return min
	}

	if n >= max {
		return max
	}

	return n
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func vkTransformBytes(data []byte) []uint32 {
	buf := make([]uint32, len(data)/4)
	vulkan.Memcopy(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buf)).Data), data)
	return buf
}
