package glm

import (
	"fmt"
	"unsafe"
)

// SizeOfVec2 its size for low precision memory data dump (float32)
// dump have size of 8 bytes (x=4 + y=4)
const SizeOfVec2 = 8

// Vec2 is common vector data structure
type Vec2 struct {
	X, Y float32
}

func (v *Vec2) String() string {
	return fmt.Sprintf("Vec2{%.2f, %.2f}", v.X, v.Y)
}

func (v *Vec2) Data() []byte {
	return (*(*[SizeOfVec2]byte)(unsafe.Pointer(v)))[:]
}
