package main

import (
	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/shared/metrics"
)

var e1RenderOrderInv = false
var e1SubscribeSwitch = false

func e1RectDrawOrder(rnd *vgl.Render) {
	const slotsCount = 10
	const slotPaddingPercent = 1

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
		rnd.ListenStats(func(stats metrics.Stats) {
			if stats.FrameIndex == 0 {
				e1RenderOrderInv = !e1RenderOrderInv
			}
		})
	}

	width, height := rnd.SurfaceSize()
	slotWidth := glx.Floor(width / slotsCount)
	slotHeight := glx.Floor(height / slotsCount)
	paddingX := width / 100 * slotPaddingPercent
	paddingY := height / 100 * slotPaddingPercent

	for x := float32(0); x < width; x += slotWidth {
		for y := float32(0); y < height; y += slotHeight {

			var filled bool

			if e1RenderOrderInv {
				// 2 draw-calls
				filled = x >= width/2
			} else {
				// 8 draw-calls
				filled = y >= height/2
			}

			rnd.Draw2dRect(&vgl.Params2dRect{
				Pos: [4]glx.Vec2{
					{x + paddingX, y + paddingY},
					{(x + slotWidth) - paddingX, y + paddingY},
					{(x + slotWidth) - paddingX, (y + slotHeight) - paddingY},
					{x + paddingX, (y + slotHeight) - paddingY},
				},
				Color:  colMain,
				Filled: filled,
			})

		}
	}

}
