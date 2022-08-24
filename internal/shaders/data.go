package shaders

import (
	_ "embed"
)

var (
	//go:embed univ2d.vert.spv
	univ2dCodeVert []byte
	//go:embed univ2d.frag.spv
	univ2dCodeFrag []byte

	//go:embed circle2d.vert.spv
	circle2dCodeVert []byte
	//go:embed circle2d.frag.spv
	circle2dCodeFrag []byte
)

func Universal2DVertSpv() []byte {
	return univ2dCodeVert
}

func Universal2DFragSpv() []byte {
	return univ2dCodeFrag
}

func Circle2DVertSpv() []byte {
	return circle2dCodeVert
}

func Circle2DFragSpv() []byte {
	return circle2dCodeFrag
}
