package vlk

import (
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

const maxSurfaces = math.MaxUint8
const maxGroups = math.MaxUint16
const defaultSurfacesCapacity = 4
const defaultGroupCapacity = 64
const defaultCallsCapacity = 32
const defaultInstancesCapacity = 512

type (
	drawContext struct {
		available      bool           // is render available?
		currentImageID uint32         // current swapChain imageID
		surfaces       []*drawSurface // surfaces render queue
	}

	drawSurface struct {
		// baking
		surfaceID surfaceID    // current surface ID
		groups    []*drawGroup // shader groups for drawing

		// dynamic
		uniform vulkan.DescriptorSet // global UBO set (view, projection)
	}

	drawGroup struct {
		// baking
		shader      *shader.Shader        // ref to group shader
		instances   []shader.InstanceData // raw instances data that should be used for drawing
		polygonMode vulkan.PolygonMode    // render polygon mode
		// blendingMode vulkan.BlendOp // todo: blending

		// dynamic
		renderPipe pipeline.Info // created vk pipeline object for group params
		indexes    bufferBinding // index buffer info
		calls      []*drawCall   // instances transformed to calls
	}

	drawCall struct {
		instances []shader.InstanceData
		vertexes  bufferBinding
		uniforms  []vulkan.DescriptorSet
	}

	bufferBinding struct {
		used   bool
		buffer vulkan.Buffer
		offset vulkan.DeviceSize
	}
)

func newDrawContext() *drawContext {
	return &drawContext{
		surfaces: make([]*drawSurface, 0, defaultSurfacesCapacity),
	}
}

func newDrawSurface(ID surfaceID) *drawSurface {
	return &drawSurface{
		surfaceID: ID,
		groups:    make([]*drawGroup, 0, defaultGroupCapacity),
	}
}

func newDrawGroup(sdr *shader.Shader, polygonMode vulkan.PolygonMode) *drawGroup {
	return &drawGroup{
		shader:      sdr,
		instances:   make([]shader.InstanceData, 0, defaultInstancesCapacity),
		polygonMode: polygonMode,
		calls:       make([]*drawCall, 0, defaultCallsCapacity),
	}
}
