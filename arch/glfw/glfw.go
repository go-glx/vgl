package glfw

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type (
	GLFW struct {
		custom *Custom
	}
)

func NewGLFW(
	appName string,
	engineName string,
	fullscreen bool,
	width int,
	height int,
) *GLFW {
	// init
	err := glfw.Init()
	if err != nil {
		panic(fmt.Errorf("failed init glfw library: %w", err))
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	// create window
	var monitor *glfw.Monitor
	if fullscreen {
		monitor = glfw.GetPrimaryMonitor()
	}

	//
	window, err := glfw.CreateWindow(width, height, appName, monitor, nil)
	if err != nil {
		panic(fmt.Errorf("failed create glfw window: %w", err))
	}

	return &GLFW{
		custom: NewCustomGLFW(appName, engineName, window),
	}
}

func (g *GLFW) AppName() string {
	return g.custom.AppName()
}

func (g *GLFW) EngineName() string {
	return g.custom.EngineName()
}

func (g *GLFW) OnWindowResized(f func(width int, height int)) {
	g.custom.OnWindowResized(f)
}

func (g *GLFW) CreateSurface(inst vulkan.Instance) (vulkan.Surface, error) {
	return g.custom.CreateSurface(inst)
}

func (g *GLFW) GetRequiredInstanceExtensions() []string {
	return g.custom.GetRequiredInstanceExtensions()
}

func (g *GLFW) GetFramebufferSize() (width, height int) {
	return g.custom.GetFramebufferSize()
}

func (g *GLFW) InitVulkanProcAddr() unsafe.Pointer {
	return g.custom.InitVulkanProcAddr()
}

func (g *GLFW) Close() error {
	glfw.Terminate()

	return g.custom.Close()
}
