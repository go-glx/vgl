package main

import (
	"os"
	"os/signal"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/arch/glfw"
	"github.com/go-glx/vgl/shared/config"
)

func main() {
	wm := glfw.NewGLFW("triangle", "triangle", false, 640, 480)
	app := vgl.NewRender(wm, config.NewConfig())

	appAlive := true

	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, os.Kill, os.Interrupt)

		select {
		case <-sigCh:
			appAlive = false
			return
		}
	}()

	for appAlive {
		app.FrameStart()

		width, height := app.SurfaceSize()

		app.Draw2dTriangle(&vgl.Params2dTriangle{
			Pos: [3]glx.Vec2{ // in clock-wise order
				{X: width / 2, Y: 100},            // center-top vertex
				{X: width - 100, Y: height - 100}, // right-bottom vertex
				{X: 100, Y: height - 100},         // left-bottom vertex
			},
			ColorGradient: [3]glx.Color{
				glx.ColorRed,
				glx.ColorGreen,
				glx.ColorBlue,
			},
			ColorUseGradient: true,
			Filled:           true,
		})

		app.FrameEnd()
	}

	// always should be closed on exit
	// this will clean vulkan resources in GPU/system
	_ = app.Close()
}
