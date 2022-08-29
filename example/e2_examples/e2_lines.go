package main

import (
	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
)

func e2Lines(rnd *vgl.Render) {
	const unitSize = 18          // each character size
	const offsetY = unitSize * 4 // offset between features

	// Drawing "VGL" word with lines
	chars := [][]glx.Vec2{
		{
			// V
			{2, 1},
			{3, 4},
			{4, 1},
		},
		{
			// G
			{7, 1},
			{5, 1},
			{5, 4},
			{7, 4},
			{7, 2},
			{6, 2},
		},
		{
			// L
			{8, 1},
			{8, 4},
			{10, 4},
		},
	}

	features := []vgl.Params2dLine{
		{
			// simple line
			Color: colGrayDark,
		},
		{
			// same as first, because default width=1
			Color: colGrayDark,
			Width: 1,
		},
		{
			// gradient line
			ColorGradient: [2]glx.Color{
				glx.ColorRed,
				glx.ColorGreen,
			},
			ColorUseGradient: true,
		},
		{
			// bold line
			Color: colGrayDark,
			Width: 10,
		},
		{
			// bold + gradient line
			ColorGradient: [2]glx.Color{
				glx.ColorRed,
				glx.ColorGreen,
			},
			ColorUseGradient: true,
			Width:            10,
		},
	}

	for featInd, feature := range features {
		for _, points := range chars {
			for pointInd, curr := range points {
				if pointInd == len(points)-1 {
					break
				}

				next := points[pointInd+1]

				// abstractPos -> realPos
				curr = curr.Scale(unitSize)
				next = next.Scale(unitSize)

				// add feature Y offset
				curr.Y += float32(featInd * offsetY)
				next.Y += float32(featInd * offsetY)

				feature.Pos = [2]glx.Vec2{curr, next}

				rnd.Draw2dLine(&feature)
			}
		}
	}
}
