package glm

import (
	"fmt"
	"unsafe"
)

// SizeOfVec4 its size for low precision memory data dump (float32)
// dump have size of 16 bytes (x=4 + y=4 + z=4 + a=4)
const SizeOfVec4 = 16

// Vec4 is common vector data structure
type Vec4 struct {
	R, G, B, A float64
}

func (v *Vec4) String() string {
	return fmt.Sprintf("Vec4{%.2f, %.2f, %.2f, %.2f}", v.R, v.G, v.B, v.A)
}

func (v *Vec4) Data() []byte {
	return (*(*[SizeOfVec4]byte)(unsafe.Pointer(v)))[:]
}
