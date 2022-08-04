package arch

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type (
	GLFWCustom struct {
		appName    string
		engineName string

		glfwWindow *glfw.Window

		windowResizeCb *glfwResizeCallback
	}

	glfwResizeCallback = func(width int, height int)
)

func NewCustomGLFW(
	appName string,
	engineName string,
	window *glfw.Window,
) *GLFWCustom {
	wm := &GLFWCustom{
		appName:    appName,
		engineName: engineName,
		glfwWindow: window,
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		if wm.windowResizeCb != nil {
			cb := *wm.windowResizeCb
			cb(width, height)
		}
	})

	return wm
}

func (g *GLFWCustom) AppName() string {
	return g.appName
}

func (g *GLFWCustom) EngineName() string {
	return g.engineName
}

func (g *GLFWCustom) OnWindowResized(f func(width int, height int)) {
	g.windowResizeCb = &f
}

func (g *GLFWCustom) CreateSurface(inst vulkan.Instance) (vulkan.Surface, error) {
	surfacePtr, err := g.glfwWindow.CreateWindowSurface(inst, nil)
	if err != nil {
		return nil, fmt.Errorf("failed create glfw surface: %w", err)
	}

	return vulkan.SurfaceFromPointer(surfacePtr), nil
}

func (g *GLFWCustom) GetRequiredInstanceExtensions() []string {
	return g.glfwWindow.GetRequiredInstanceExtensions()
}

func (g *GLFWCustom) GetFramebufferSize() (width, height int) {
	return g.glfwWindow.GetSize()
}

func (g *GLFWCustom) InitVulkanProcAddr() unsafe.Pointer {
	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		panic(fmt.Errorf("failed get vulkan proc address"))
	}

	vulkan.SetGetInstanceProcAddr(procAddr)
	return procAddr
}

func (g *GLFWCustom) Close() error {
	return nil
}
