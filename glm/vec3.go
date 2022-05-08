package glm

import (
	"fmt"
	"unsafe"
)

// SizeOfVec3 its size for low precision memory data dump (float32)
// dump have size of 12 bytes (x=4 + y=4 + z=4)
const SizeOfVec3 = 12

// Vec3 is common vector data structure
type Vec3 struct {
	R, G, B float32
}

func (v *Vec3) String() string {
	return fmt.Sprintf("Vec3{%.2f, %.2f, %.2f}", v.R, v.G, v.B)
}

func (v *Vec3) Data() []byte {
	return (*(*[SizeOfVec3]byte)(unsafe.Pointer(v)))[:]
}
