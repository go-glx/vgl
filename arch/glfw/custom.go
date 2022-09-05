package glfw

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type (
	Custom struct {
		appName    string
		engineName string
		window     *glfw.Window

		windowResizeCb *glfwResizeCallback
	}

	glfwResizeCallback = func(width int, height int)
)

func NewCustomGLFW(
	appName string,
	engineName string,
	window *glfw.Window,
) *Custom {
	wm := &Custom{
		appName:    appName,
		engineName: engineName,
		window:     window,
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		if wm.windowResizeCb != nil {
			cb := *wm.windowResizeCb
			cb(width, height)
		}
	})

	return wm
}

func (g *Custom) AppName() string {
	return g.appName
}

func (g *Custom) EngineName() string {
	return g.engineName
}

func (g *Custom) OnWindowResized(f func(width int, height int)) {
	g.windowResizeCb = &f
}

func (g *Custom) CreateSurface(inst vulkan.Instance) (vulkan.Surface, error) {
	surfacePtr, err := g.window.CreateWindowSurface(inst, nil)
	if err != nil {
		return nil, fmt.Errorf("failed create glfw surface: %w", err)
	}

	return vulkan.SurfaceFromPointer(surfacePtr), nil
}

func (g *Custom) GetRequiredInstanceExtensions() []string {
	return g.window.GetRequiredInstanceExtensions()
}

func (g *Custom) GetFramebufferSize() (width, height int) {
	return g.window.GetSize()
}

func (g *Custom) InitVulkanProcAddr() unsafe.Pointer {
	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		panic(fmt.Errorf("failed get vulkan proc address"))
	}

	vulkan.SetGetInstanceProcAddr(procAddr)
	return procAddr
}

func (g *Custom) Close() error {
	return nil
}
