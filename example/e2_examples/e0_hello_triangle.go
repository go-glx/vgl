package main

import (
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

func e0HelloTriangle(rnd *vgl.Render) {
	width, height := rnd.SurfaceSize()

	rnd.Draw2dTriangle(&vgl.Params2dTriangle{
		Pos: [3]glm.Local2D{ // in clock-wise order
			{width / 2, 100},            // center-top vertex
			{width - 100, height - 100}, // right-bottom vertex
			{100, height - 100},         // left-bottom vertex
		},
		ColorGradient: [3]glm.Color{
			glm.NewColor(255, 0, 0, 255),
			glm.NewColor(0, 255, 0, 255),
			glm.NewColor(0, 0, 255, 255),
		},
		ColorUseGradient: true,
		Filled:           true,
	})
}
