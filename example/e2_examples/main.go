package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/arch"
	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/glm"
)

const appWidth = 720
const appHeight = 480

// todo: switch in demo UI and/or keyboard input
const currentDemo = demoE4Circles

const (
	demoE0HelloTriangle = iota
	demoE1RectDrawOrder
	demoE2Lines
	demoE3Points
	demoE4Circles
)

const enablePprof = false

var demos = map[int]func(rnd *vgl.Render){
	demoE0HelloTriangle: e0HelloTriangle,
	demoE1RectDrawOrder: e1RectDrawOrder,
	demoE2Lines:         e2Lines,
	demoE3Points:        e3Points,
	demoE4Circles:       e4Circles,
}

func main() {
	wm := arch.NewGLFW("examples", "VGL", false, appWidth, appHeight)
	rnd := vgl.NewRender(wm, config.NewConfig(
		config.WithDebug(true),
		config.WithMobileFriendly(true),
	))
	rnd.ListenStats(onFrameEndStats)
	gui := newGUI(rnd)

	appAlive := true
	go listenSignals(&appAlive)
	go pprof()

	demo := demos[currentDemo]

	for appAlive {
		rnd.FrameStart()
		gui.frameStart()

		// draw bg
		w, h := rnd.SurfaceSize()
		rnd.Draw2dRect(&vgl.Params2dRect{
			Pos:    [4]glm.Local2D{{0, 0}, {w, 0}, {w, h}, {0, h}},
			Color:  colBackground,
			Filled: true,
		})

		// draw demo scene
		demo(rnd)

		gui.frameEnd()
		rnd.FrameEnd()
	}

	// always should be closed on exit
	// this will clean vulkan resources in GPU/system
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

func pprof() {
	if !enablePprof {
		return
	}

	log.Println(http.ListenAndServe("localhost:47887", nil))
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

	printMem := func(size uint32, capacity uint32) string {
		return fmt.Sprintf("%.1fmb/%.1fmb",
			float32(size)/1024/1024,
			float32(capacity)/1024/1024,
		)
	}

	// print AVG stats of each frame
	elapsed := time.Since(startTime)
	fmt.Printf(""+
		"+%3.0fs | FPS=%-2d  drawCalls=%.0f \n"+
		"  us> | Pipl=%-4.0f  WrtVert=%-4.0f  DrwInd=%-4.0f  Drw=%-4.0f\n"+
		"  mem | total=%s  vert=%s  ind=%s  ubo=%s  ssbo=%s\n",
		elapsed.Seconds(),
		stats.FPS,
		float32(combinedStats.DrawCalls)/avg,
		float32(combinedStats.TimeCreatePipeline.Microseconds())/avg,
		float32(combinedStats.TimeFlushVertexBuffer.Microseconds())/avg,
		float32(combinedStats.TimeRenderInstanced.Microseconds())/avg,
		float32(combinedStats.TimeRenderFallback.Microseconds())/avg,
		printMem(stats.Memory.TotalSize, stats.Memory.TotalCapacity),
		printMem(stats.Memory.VertexBuffers.Size, stats.Memory.VertexBuffers.Capacity),
		printMem(stats.Memory.IndexBuffers.Size, stats.Memory.IndexBuffers.Capacity),
		printMem(stats.Memory.UniformBuffers.Size, stats.Memory.UniformBuffers.Capacity),
		printMem(stats.Memory.StorageBuffers.Size, stats.Memory.StorageBuffers.Capacity),
	)

	combinedStats.Reset()
}
