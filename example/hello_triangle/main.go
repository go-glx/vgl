package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/glm"
)

const appWidth = 640
const appHeight = 320

func main() {
	wm := arch.NewGLFW("example", "example", false, appWidth, appHeight)
	api := vgl.NewRender(wm, config.NewConfig())

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

	fps := 0

	go func() {
		tmr := time.NewTicker(time.Second)
		for {
			select {
			case <-tmr.C:
				fmt.Println(fps)
				fps = 0
			}
		}
	}()

	for appAlive {
		api.FrameStart()

		// hello triangle
		const count = 100

		for i := 0; i <= count; i++ {
			offsetY := (float32(i) / float32(count)) * 50

			api.Draw2dTriangle(&vgl.Params2dTriangle{
				Pos: [3]glm.Local2D{ // in clock-wise order
					{appWidth / 2, 100 + int32(offsetY)}, // center-top vertex
					{appWidth - 100, appHeight - 100},    // right-bottom vertex
					{100, appHeight - 100},               // left-bottom vertex
				},
				ColorGradient: [3]glm.Color{
					glm.NewColor(255, 0, 0, 255),
					glm.NewColor(0, 255, 0, 255),
					glm.NewColor(0, 0, 255, 255),
				},
				ColorUseGradient: true,
				Filled:           false,
			})
		}

		api.FrameEnd()
		fps++
	}

	// always should be closed on exit
	// this will clean vulkan resources in GPU/system
	api.WaitGPU()
	_ = api.Close()
}
