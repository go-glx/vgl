package glm

import (
	"fmt"
	"math"
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

func (v *Vec2) AngleTo(to Vec2) (rad float64) {
	return math.Atan2(float64(to.Y-v.Y), float64(v.X-to.X)) + math.Pi
}

func (v *Vec2) PolarOffset(distance float32, rad float64) Vec2 {
	return Vec2{
		X: v.X + distance*float32(math.Cos(rad)),
		Y: v.Y - distance*float32(math.Sin(rad)),
	}
}

func (v Vec2) Add(x, y float32) Vec2 {
	return Vec2{
		X: v.X + x,
		Y: v.Y + y,
	}
}
