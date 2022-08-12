package main

import (
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

func e2Lines(rnd *vgl.Render) {
	const offsetY = 150

	points := []glm.Local2D{
		{100, 100},
		{150, 120},
		{240, 130},
		{300, 110},
		{500, 130},
		{650, 100},
		{640, 50},
		{520, 40},
		{340, 60},
		{200, 75},
		{125, 50},
		{100, 100},
	}

	features := []vgl.Params2dLine{
		{
			// simple line
			Color: glm.NewColor(255, 255, 0, 255),
		},
		{
			// same as first, because default width=1
			Color: glm.NewColor(255, 255, 0, 255),
			Width: 1,
		},
		{
			// bold line
			Color: glm.NewColor(255, 255, 0, 255),
			Width: 5,
		},
		{
			// gradient line
			ColorGradient: [2]glm.Color{
				glm.NewColor(255, 0, 0, 255),
				glm.NewColor(0, 0, 255, 255),
			},
			ColorUseGradient: true,
		},
	}

	for featInd, feature := range features {
		for pointInd, curr := range points {
			if pointInd == len(points)-1 {
				break
			}

			next := points[pointInd+1]

			curr[1] += int32(featInd * offsetY)
			next[1] += int32(featInd * offsetY)

			feature.Pos = [2]glm.Local2D{curr, next}

			rnd.Draw2dLine(&feature)
		}
	}

}
