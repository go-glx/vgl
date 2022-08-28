package main

import (
	"time"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

type animRadius struct{}
type animTransformX struct{}
type animTransformY struct{}
type animHoleSize struct{}
type animSmooth struct{}

func e4Circles(rnd *vgl.Render) {
	w, h := rnd.SurfaceSize()
	centerX, centerY := w/2, h/2

	radius := int32(anim(animRadius{}, time.Second*5, 150, 300))
	transformX := int32(anim(animTransformX{}, time.Second*10, -200, 200))
	transformY := int32(anim(animTransformY{}, time.Second*15, -50, 50))
	holeSize := anim(animHoleSize{}, time.Second*25, 0, 0.99)
	smooth := anim(animSmooth{}, time.Second*5, 0.005, 0.01)

	rnd.Draw2dCircle(&vgl.Params2dCircle{
		Center: [2]int32{centerX + transformX, centerY + transformY},
		Radius: radius,
		ColorGradient: [4]glm.Color{
			glm.NewColor(255, 0, 0, 255),
			glm.NewColor(0, 255, 0, 255),
			glm.NewColor(0, 0, 255, 255),
			glm.NewColor(255, 128, 0, 255),
		},
		ColorUseGradient: true,
		HoleRadius:       holeSize,
		Smooth:           smooth,
	})

	rnd.Draw2dCircle(&vgl.Params2dCircle{
		Center: [2]int32{centerX - transformX, centerY + transformY},
		Radius: radius,
		ColorGradient: [4]glm.Color{
			glm.NewColor(255, 0, 0, 255),
			glm.NewColor(0, 255, 0, 255),
			glm.NewColor(0, 0, 255, 255),
			glm.NewColor(255, 128, 0, 255),
		},
		ColorUseGradient: true,
		HoleRadius:       holeSize,
		Smooth:           smooth,
	})
}
