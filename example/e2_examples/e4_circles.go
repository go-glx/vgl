package main

import (
	"time"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
)

type animRadius struct{}
type animTransformX struct{}
type animTransformY struct{}
type animHoleSize struct{}
type animSmooth struct{}

func e4Circles(rnd *vgl.Render) {
	w, h := rnd.SurfaceSize()
	centerX, centerY := w/2, h/2

	radius := anim(animRadius{}, time.Second*5, 150, 300)
	transformX := anim(animTransformX{}, time.Second*10, -200, 200)
	transformY := anim(animTransformY{}, time.Second*15, -50, 50)
	holeSize := anim(animHoleSize{}, time.Second*25, 0, 0.99)
	smooth := anim(animSmooth{}, time.Second*5, 0.005, 0.01)

	rnd.Draw2dCircle(&vgl.Params2dCircle{
		Center: glx.Vec2{X: centerX + transformX, Y: centerY + transformY},
		Radius: radius,
		ColorGradient: [4]glx.Color{
			glx.ColorRed,
			glx.ColorGreen,
			glx.ColorBlue,
			glx.ColorWhite,
		},
		ColorUseGradient: true,
		HoleRadius:       holeSize,
		Smooth:           smooth,
	})

	// todo: remove
	rnd.Draw2dCircle(&vgl.Params2dCircle{
		Center: glx.Vec2{X: centerX - transformX, Y: centerY + transformY},
		Radius: radius,
		ColorGradient: [4]glx.Color{
			glx.ColorRed,
			glx.ColorGreen,
			glx.ColorBlue,
			glx.ColorWhite,
		},
		ColorUseGradient: true,
		HoleRadius:       holeSize,
		Smooth:           smooth,
	})
}
