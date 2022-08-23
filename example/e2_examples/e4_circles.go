package main

import (
	"time"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/glm"
)

type animRadiusX struct{}
type animRadiusY struct{}
type animTransformX struct{}
type animTransformY struct{}

func e4Circles(rnd *vgl.Render) {
	w, h := rnd.SurfaceSize()
	centerX, centerY := w/2, h/2

	radiusW := int32(anim(animRadiusX{}, time.Second*5, 50, 150))
	radiusH := int32(anim(animRadiusY{}, time.Second*30, 100, 200))

	transformX := int32(anim(animTransformX{}, time.Second*10, -100, 100))
	transformY := int32(anim(animTransformY{}, time.Second*15, -50, 50))

	rnd.Draw2dCircle(&vgl.Params2dCircle{
		PosArea: [4]glm.Local2D{
			{transformX + (centerX - radiusW), transformY + (centerY - radiusH)},
			{transformX + (centerX + (radiusW * 2)), transformY + (centerY - radiusH)},
			{transformX + (centerX + radiusW), transformY + (centerY + radiusH)},
			{transformX + (centerX - (radiusW * 2)), transformY + (centerY + radiusH)},
		},
		PosUseArea: true,
		ColorGradient: [4]glm.Color{
			glm.NewColor(255, 0, 0, 255),
			glm.NewColor(0, 255, 0, 255),
			glm.NewColor(0, 0, 255, 255),
			glm.NewColor(255, 128, 0, 255),
		},
		ColorUseGradient: true,
	})
}
