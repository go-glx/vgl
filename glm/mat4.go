package glm

import (
	"fmt"
	"unsafe"
)

const SizeOfMat4 = SizeOfVec4 * 4

type Mat4 struct {
	A Vec4
	B Vec4
	C Vec4
	D Vec4
}

func (v *Mat4) String() string {
	return fmt.Sprintf(`Mat4[
     %0.2f, %0.2f, %0.2f, %0.2f,
     %0.2f, %0.2f, %0.2f, %0.2f,
     %0.2f, %0.2f, %0.2f, %0.2f,
     %0.2f, %0.2f, %0.2f, %0.2f,
]`,

		v.A.R, v.A.G, v.A.B, v.A.A,
		v.B.R, v.B.G, v.B.B, v.B.A,
		v.C.R, v.C.G, v.C.B, v.C.A,
		v.D.R, v.D.G, v.D.B, v.D.A,
	)
}

func (v *Mat4) Data() []byte {
	return (*(*[SizeOfMat4]byte)(unsafe.Pointer(v)))[:]
}

func Mat4Identity() Mat4 {
	return Mat4{
		A: Vec4{1, 0, 0, 0},
		B: Vec4{0, 1, 0, 0},
		C: Vec4{0, 0, 1, 0},
		D: Vec4{0, 0, 0, 1},
	}
}
