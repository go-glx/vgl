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

const appWidth = 960
const appHeight = 720

const (
	demoE1RectDrawOrder = iota
	demoE2Lines
	demoE3Points
)

// todo: switch in demo UI and/or keyboard input
const currentDemo = demoE2Lines

var demos = map[int]func(rnd *vgl.Render){
	demoE1RectDrawOrder: e1RectDrawOrder,
	demoE2Lines:         e2Lines,
	demoE3Points:        e3Points,
}

func main() {
	wm := arch.NewGLFW("examples", "VGL", false, appWidth, appHeight)
	rnd := vgl.NewRender(wm, config.NewConfig(
		config.WithDebug(false),
		config.WithMobileFriendly(true),
	))
	rnd.ListenStats(onFrameEndStats)
	gui := newGUI(rnd)

	appAlive := true
	go listenSignals(&appAlive)

	demo := demos[currentDemo]

	for appAlive {
		rnd.FrameStart()
		gui.frameStart()

		demo(rnd)

		gui.frameEnd()
		rnd.FrameEnd()
	}

	// always should be closed on exit
	// this will clean vulkan resources in GPU/system
	rnd.WaitGPU()
	_ = rnd.Close()
}

func listenSignals(appAlive *bool) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Kill, os.Interrupt)

	select {
	case <-sigCh:
		*appAlive = false
		return
	}
}

var startTime = time.Now()
var combinedStats = glm.Stats{}

func onFrameEndStats(stats glm.Stats) {
	if stats.FrameIndex != 0 {
		// print one time per second
		combinedStats.DrawCalls += stats.DrawCalls
		combinedStats.DrawChunks += stats.DrawChunks
		combinedStats.DrawUniqueShaders += stats.DrawUniqueShaders
		combinedStats.TimeCreatePipeline += stats.TimeCreatePipeline
		combinedStats.TimeFlushVertexBuffer += stats.TimeFlushVertexBuffer
		combinedStats.TimeRenderInstanced += stats.TimeRenderInstanced
		combinedStats.TimeRenderFallback += stats.TimeRenderFallback
		return
	}

	avg := float32(stats.FPS)
	if avg <= 0 {
		avg = 0.0001
	}

	// print AVG stats of each frame
	elapsed := time.Since(startTime)
	fmt.Printf(""+
		"+%3.0fs | FPS=%-2d  drawCalls=%.0f \n"+
		"  us> | Pipl=%-4.0f  WrtVert=%-4.0f  DrwInd=%-4.0f  Drw=%-4.0f\n",
		elapsed.Seconds(),
		stats.FPS,
		float32(combinedStats.DrawCalls)/avg,
		float32(combinedStats.TimeCreatePipeline.Microseconds())/avg,
		float32(combinedStats.TimeFlushVertexBuffer.Microseconds())/avg,
		float32(combinedStats.TimeRenderInstanced.Microseconds())/avg,
		float32(combinedStats.TimeRenderFallback.Microseconds())/avg,
	)

	combinedStats.Reset()
}
