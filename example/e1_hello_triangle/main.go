package main

import (
	"os"
	"os/signal"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/glm"
)

const appWidth = 640
const appHeight = 480

func main() {
	wm := arch.NewGLFW("triangle", "triangle", false, appWidth, appHeight)
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

		app.Draw2dTriangle(&vgl.Params2dTriangle{
			Pos: [3]glm.Local2D{ // in clock-wise order
				{appWidth / 2, 100},               // center-top vertex
				{appWidth - 100, appHeight - 100}, // right-bottom vertex
				{100, appHeight - 100},            // left-bottom vertex
			},
			ColorGradient: [3]glm.Color{
				glm.NewColor(255, 0, 0, 255),
				glm.NewColor(0, 255, 0, 255),
				glm.NewColor(0, 0, 255, 255),
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
