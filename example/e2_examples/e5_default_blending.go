package main

import (
	"math"
	"time"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
)

type animBlendingAngle struct{}
type animBlendingSize struct{}
type animBlendingMerge struct{}

func e5DefaultBlending(rnd *vgl.Render) {
	type figureRenderFn = func(rnd *vgl.Render, center glx.Vec2, angle glx.Angle, size float32)
	figures := []figureRenderFn{
		e5DefaultBlendingDrawBox,
		e5DefaultBlendingDrawCircle,
		e5DefaultBlendingDrawTriangle,
	}

	w, h := rnd.SurfaceSize()
	centerX := w / 2
	gridX := []float32{centerX - (centerX / 2), centerX + (centerX / 2)}
	gridY := h / float32(len(figures)+1)

	angle := anim(animBlendingAngle{}, time.Second*10, 0, math.Pi*2)
	size := anim(animBlendingSize{}, time.Second*30, gridY, gridY/1.5)
	mergeX := anim(animBlendingMerge{}, time.Second*10, centerX/3, centerX/2)

	// todo: change order of for loops
	//       this will show bug with overridden ssbo
	for yInd, figure := range figures {
		for xInd, xx := range gridX {
			origin := glx.Vec2{
				X: xx,
				Y: gridY * float32(yInd+1),
			}

			if xInd == 0 {
				origin.X += mergeX
			} else {
				origin.X -= mergeX
			}

			figure(rnd, origin, glx.Radians(angle+glx.Angle180*float32(xInd)), size)
		}
	}
}

func e5DefaultBlendingDrawBox(rnd *vgl.Render, origin glx.Vec2, angle glx.Angle, size float32) {
	rnd.Draw2dRect(&vgl.Params2dRect{
		Pos: [4]glx.Vec2{
			origin.PolarOffset(size, angle+glx.Angle90+glx.Angle45),
			origin.PolarOffset(size, angle+glx.Angle90-glx.Angle45),
			origin.PolarOffset(size, angle+glx.Angle270+glx.Angle45),
			origin.PolarOffset(size, angle+glx.Angle270-glx.Angle45),
		},
		ColorGradient: [4]glx.Color{
			glx.ColorRed,
			glx.ColorGreen,
			glx.ColorBlue,
			glx.ColorWhite,
		},
		ColorUseGradient: true,
		Filled:           true,
	})
}

func e5DefaultBlendingDrawTriangle(rnd *vgl.Render, origin glx.Vec2, angle glx.Angle, size float32) {
	rnd.Draw2dTriangle(&vgl.Params2dTriangle{
		Pos: [3]glx.Vec2{
			origin.PolarOffset(size, angle+glx.Angle90),
			origin.PolarOffset(size, angle+glx.Angle270+glx.Angle45),
			origin.PolarOffset(size, angle+glx.Angle270-glx.Angle45),
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

func e5DefaultBlendingDrawCircle(rnd *vgl.Render, origin glx.Vec2, _ glx.Angle, size float32) {
	rnd.Draw2dCircle(&vgl.Params2dCircle{
		PosCenter: origin,
		PosRadius: size,
		ColorGradient: [4]glx.Color{
			glx.ColorRed,
			glx.ColorGreen,
			glx.ColorBlue,
			glx.ColorWhite,
		},
		HoleRadius:       0.25,
		Smooth:           0.005,
		ColorUseGradient: true,
	})
}
