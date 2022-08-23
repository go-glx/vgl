package main

import (
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

var e1RenderOrderInv = false
var e1SubscribeSwitch = false

func e1RectDrawOrder(rnd *vgl.Render) {
	const rectSize = 50
	const padding = 10

	// this demonstrates how much order of drawing
	// affect draw-calls instancing

	// Drawing top->bottom, left->right
	// O - outline
	// F - filled

	// #0 Single draw call - when [all rect is filled] OR [all rect is outline]

	// #1 Drawing order [> halfWidth] (2 draw calls)
	// 1_O  3_O  5_F  7_F
	// 2_O  4_O  6_F  8_F

	// #2 Drawing order [> halfHeight] (8 draw calls)
	// 1_O  3_O  5_O  7_O
	// 2_F  4_F  6_F  8_F

	if !e1SubscribeSwitch {
		e1SubscribeSwitch = true
		rnd.ListenStats(func(stats glm.Stats) {
			if stats.FrameIndex == 0 {
				e1RenderOrderInv = !e1RenderOrderInv
			}
		})
	}

	for x := int32(padding); x < appWidth-padding; x += rectSize + padding {
		for y := int32(padding); y < appHeight-padding; y += rectSize + padding {
			var filled bool

			if e1RenderOrderInv {
				// 2 draw-calls
				filled = x > appWidth/2
			} else {
				// 8 draw-calls
				filled = y > appHeight/2
			}

			rnd.Draw2dRect(&vgl.Params2dRect{
				Pos: [4]glm.Local2D{
					{x, y},
					{x + rectSize, y},
					{x + rectSize, y + rectSize},
					{x, y + rectSize},
				},
				Color:  colMain,
				Filled: filled,
			})

		}
	}
}
