package main

import (
	"time"

	"github.com/inkyblackness/imgui-go/v4"

	"github.com/go-glx/vgl"
)

type gui struct {
	app      *vgl.Render
	context  *imgui.Context
	io       imgui.IO
	prevTime time.Time
}

func newGUI(app *vgl.Render) *gui {
	context := imgui.CreateContext(nil)
	io := imgui.CurrentIO()

	// todo: font drawing
	fontImage := io.Fonts().TextureDataAlpha8()
	_ = fontImage // todo: load font image to vulkan
	fontTexID := 1
	io.Fonts().SetTextureID(imgui.TextureID(fontTexID))

	return &gui{
		app:      app,
		context:  context,
		io:       io,
		prevTime: time.Now(),
	}
}

func (g *gui) close() {
	g.context.Destroy()
}

func (g *gui) frameStart() {
	// set win size
	w, h := g.app.SurfaceSize()
	g.io.SetDisplaySize(imgui.Vec2{X: float32(w), Y: float32(h)})

	// calc delta time
	now := time.Now()
	dt := now.Sub(g.prevTime)
	g.prevTime = now
	g.io.SetDeltaTime(float32(dt.Seconds()))

	// start frame
	imgui.NewFrame()
}

func (g *gui) frameEnd() {
	imgui.Render()
	data := imgui.RenderedDrawData()

	// todo: draw data
	_ = data
}
