package main

import (
	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
)

func e0HelloTriangle(rnd *vgl.Render) {
	width, height := rnd.SurfaceSize()

	rnd.Draw2dTriangle(&vgl.Params2dTriangle{
		Pos: [3]glx.Vec2{ // in clock-wise order
			{width / 2, 100},            // center-top vertex
			{width - 100, height - 100}, // right-bottom vertex
			{100, height - 100},         // left-bottom vertex
		},
		ColorGradient: [3]glx.Color{
			glx.ColorRed,
			glx.ColorGreen,
			glx.ColorBlue,
		},
		ColorUseGradient: true,
		Filled:           true,
	})
}
