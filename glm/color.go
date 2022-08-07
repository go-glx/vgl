package glm

import (
	"encoding/binary"
)

type Color uint32

const colorValueMax = float32(255)

func NewColor(r, g, b, a uint8) Color {
	return Color(binary.BigEndian.Uint32([]byte{r, g, b, a}))
}

// VecRGBA converts encoded uint32 color into 4 byte values (r,g,b,a)
// and cast in to Vec4 float32 in range of (0 .. 1)
func (c Color) VecRGBA() Vec4 {
	return Vec4{
		R: float32(c>>24&0xff) / colorValueMax,
		G: float32(c>>16&0xff) / colorValueMax,
		B: float32(c>>8&0xff) / colorValueMax,
		A: float32(c&0xff) / colorValueMax,
	}
}
