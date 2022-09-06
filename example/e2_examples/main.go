package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl"
	"github.com/go-glx/vgl/arch/glfw"
	"github.com/go-glx/vgl/shared/config"
	"github.com/go-glx/vgl/shared/metrics"
)

// --------------------------
// Settings
// --------------------------

const appWidth = 720
const appHeight = 480
const currentDemo = demoE4Circles
const enablePprof = true

// --------------------------
// Internal setup
// --------------------------

const (
	demoE0HelloTriangle = iota
	demoE1RectDrawOrder
	demoE2Lines
	demoE3Points
	demoE4Circles
	demoE5DefaultBlending
)

var demos = map[int]func(rnd *vgl.Render){
	demoE0HelloTriangle:   e0HelloTriangle,
	demoE1RectDrawOrder:   e1RectDrawOrder,
	demoE2Lines:           e2Lines,
	demoE3Points:          e3Points,
	demoE4Circles:         e4Circles,
	demoE5DefaultBlending: e5DefaultBlending,
}

func main() {
	wm := glfw.NewGLFW("examples", "VGL", false, appWidth, appHeight)
	rnd := vgl.NewRender(wm, config.NewConfig(
		config.WithDebug(true),
		config.WithMobileFriendly(true),
	))
	rnd.ListenStats(onFrameEndStats)
	gui := newGUI(rnd)

	appAlive := true
	go listenSignals(&appAlive)
	go pprof()

	drawDemo := demos[currentDemo]

	for appAlive {
		rnd.FrameStart()
		gui.frameStart()

		clearScreen(rnd)
		drawDemo(rnd)

		gui.frameEnd()
		rnd.FrameEnd()
	}

	// always should be closed on exit
	// this will clean vulkan resources in GPU/system
	err := rnd.Close()
	if err != nil {
		fmt.Printf("VLK close error: %s\n", err.Error())
	}
}

func clearScreen(rnd *vgl.Render) {
	w, h := rnd.SurfaceSize()
	rnd.Draw2dRect(&vgl.Params2dRect{
		Pos:    [4]glx.Vec2{{0, 0}, {w, 0}, {w, h}, {0, h}},
		Color:  colBackground,
		Filled: true,
	})
}

// --------------------------
// Internal demo functions
// --------------------------

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
var combinedStats = metrics.NewStats()

func onFrameEndStats(stats metrics.Stats) {
	if stats.FrameIndex != 0 {
		// print one time per second
		combinedStats.DrawCalls += stats.DrawCalls
		combinedStats.DrawGroups += stats.DrawGroups

		for segment, duration := range stats.SegmentDuration {
			if _, exist := combinedStats.SegmentDuration[segment]; !exist {
				combinedStats.SegmentDuration[segment] = duration
			}

			combinedStats.SegmentDuration[segment] += duration
		}
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
		"  us> | %s\n"+
		"  mem | total=%s  vert=%s  ind=%s  ubo=%s  ssbo=%s\n",
		elapsed.Seconds(),
		stats.FPS,
		float32(combinedStats.DrawCalls)/avg,
		formatMetricDurations(combinedStats, avg),
		printMem(stats.Memory.TotalSize, stats.Memory.TotalCapacity),
		printMem(stats.Memory.VertexBuffers.Size, stats.Memory.VertexBuffers.Capacity),
		printMem(stats.Memory.IndexBuffers.Size, stats.Memory.IndexBuffers.Capacity),
		printMem(stats.Memory.UniformBuffers.Size, stats.Memory.UniformBuffers.Capacity),
		printMem(stats.Memory.StorageBuffers.Size, stats.Memory.StorageBuffers.Capacity),
	)

	combinedStats.Reset()
}

func formatMetricDurations(s metrics.Stats, avg float32) string {
	const maxSlowSegments = 4

	type seg struct {
		name     string
		duration time.Duration
	}

	// map -> slice
	segments := make([]seg, 0, len(s.SegmentDuration))

	for name, duration := range s.SegmentDuration {
		segments = append(segments, seg{
			name:     name,
			duration: duration,
		})
	}

	// sort by duration (slow first)
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].duration > segments[j].duration
	})

	// limit first N
	if len(segments) > maxSlowSegments {
		segments = segments[0:maxSlowSegments]
	}

	// format duration
	formatted := make([]string, 0, len(segments))
	for _, segment := range segments {
		formatted = append(formatted, fmt.Sprintf("%s=%s (%.0f%%)",
			segment.name,
			segment.duration.String(),
			segment.duration.Seconds()*100,
		))
	}

	// join in one line
	return strings.Join(formatted, ", ")
}
