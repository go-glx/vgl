package main

import (
	"time"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

type animRadius struct{}
type animPositionX struct{}
type animPositionY struct{}

func e4Circles(rnd *vgl.Render) {
	radius := anim(animRadius{}, time.Second*3, 100, 600)

	rnd.Draw2dCircle(&vgl.Params2dCircle{
		PosCenter: glm.Local2D{
			int32(anim(animPositionX{}, time.Second*5, 0, appWidth)),
			int32(anim(animPositionY{}, time.Second*25, 0, appHeight)),
		},
		PosCenterRadius: int32(radius),
		ColorGradient: [4]glm.Color{
			glm.NewColor(255, 0, 0, 255),
			glm.NewColor(0, 255, 0, 255),
			glm.NewColor(0, 0, 255, 255),
			glm.NewColor(255, 128, 0, 255),
		},
		ColorUseGradient: true,
	})
}
