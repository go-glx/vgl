package vlk

import (
	"strconv"
	"time"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/glx"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/dscptr"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/pipeline"
	"github.com/go-glx/vgl/shared/metrics"
)

type (
	drawCtxFn       = func(ctx *drawContext)
	drawSurfaceFn   = func(ctx *drawContext, surf *drawSurface)
	drawGroupFn     = func(ctx *drawContext, g *drawGroup)
	drawGroupExecFn = func(cb vulkan.CommandBuffer, ctx *drawContext, surf *drawSurface, g *drawGroup)
	drawCallExecFn  = func(cb vulkan.CommandBuffer, ctx *drawContext, surf *drawSurface, g *drawGroup, c *drawCall)
)

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Drawing order pipeline
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) initDrawingPipeline() {
	vlk.drawContext = newDrawContext()
	vlk.drawExecution = vlk.plEvery(
		vlk.plUpdateGlobalRendererVars,
		vlk.plWhenAvailable(
			vlk.plClearVertexBuffers,
			vlk.plOnEverySurface(
				vlk.plSurfaceUpdateGlobalUniform,
				vlk.plSurfaceOnEveryGroup(
					vlk.plGroupStats,
					vlk.plGroupCreateRenderingPipeline,
					vlk.plGroupFindIndexBuffer,
					vlk.plGroupUpdateVertexBuffer,
				),
			),
			vlk.plOnEverySurface(
				vlk.plSurfaceOnEveryGroupExec(
					vlk.plExecGroupBindPipeline,
					vlk.plExecGroupBindIndexBuffer,
					vlk.plExecGroupOnEveryCall(
						vlk.plExecCallUpdateLocalUniforms,
						vlk.plExecCallBindUniforms,
						vlk.plExecCallBindVertexBuffer,
						vlk.plExecCallInstancedDraw,
					),
				),
			),
			vlk.plClearContext,
		),
	)
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Utility pipe builders
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plEvery(pls ...drawCtxFn) drawCtxFn {
	return func(ctx *drawContext) {
		for ind := range pls {
			pls[ind](ctx)
		}
	}
}

func (vlk *VLK) plWhenAvailable(pls ...drawCtxFn) drawCtxFn {
	return func(ctx *drawContext) {
		if !ctx.available {
			return
		}

		for ind := range pls {
			pls[ind](ctx)
		}
	}
}

func (vlk *VLK) plOnEverySurface(surfaceFns ...drawSurfaceFn) drawCtxFn {
	return func(ctx *drawContext) {
		for _, surface := range ctx.surfaces {
			for _, fn := range surfaceFns {
				fn(ctx, surface)
			}
		}
	}
}

func (vlk *VLK) plSurfaceOnEveryGroup(groupFns ...drawGroupFn) drawSurfaceFn {
	return func(ctx *drawContext, surf *drawSurface) {
		for _, group := range surf.groups {
			if len(group.instances) == 0 {
				continue
			}

			for _, fn := range groupFns {
				fn(ctx, group)
			}
		}
	}
}

func (vlk *VLK) plSurfaceOnEveryGroupExec(groupFns ...drawGroupExecFn) drawSurfaceFn {
	return func(ctx *drawContext, surf *drawSurface) {
		vlk.cont.frameManager().FrameApplyCommands(func(_ uint32, cb vulkan.CommandBuffer) {
			for _, group := range surf.groups {
				if len(group.instances) == 0 {
					continue
				}

				for _, fn := range groupFns {
					fn(cb, ctx, surf, group)
				}
			}
		})
	}
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Functions - Global CTX
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plUpdateGlobalRendererVars(ctx *drawContext) {
	ctx.available = vlk.drawAvailable
	ctx.currentImageID = vlk.drawImageID

	vlk.stats.FrameIndex++
	vlk.stats.DrawCalls = 0
	vlk.stats.DrawGroups = 0
	vlk.stats.Memory = metrics.MemoryStats{}
}

func (vlk *VLK) plClearVertexBuffers(ctx *drawContext) {
	ts := time.Now()

	vlk.cont.allocBuffers().ClearVertexBuffersOwnedBy(ctx.currentImageID)

	vlk.stats.SegmentDuration[metrics.SegmentPlClearBuffers] += time.Since(ts)
}

func (vlk *VLK) plClearContext(ctx *drawContext) {
	ctx.available = false
	ctx.surfaces = make([]*drawSurface, 0, defaultSurfacesCapacity)
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Functions - Surfaces
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plSurfaceUpdateGlobalUniform(ctx *drawContext, surf *drawSurface) {
	ts := time.Now()

	// write global data to uniform buffers that need for every shader
	view := glx.Mat4Identity()       // todo
	projection := glx.Mat4Identity() // todo
	surfaceSize := glx.Vec2{
		X: vlk.surfacesSize[surf.surfaceID][0],
		Y: vlk.surfacesSize[surf.surfaceID][1],
	}

	uboData := make([]byte, 0, glx.SizeOfMat4*2)
	uboData = append(uboData, view.Data()...)
	uboData = append(uboData, projection.Data()...)

	surf.uniform = vlk.cont.descriptorsManager().UpdateSet(
		ctx.currentImageID,
		dscptr.LayoutIndexGlobal,
		map[uint32][]byte{
			0: uboData,            // layout=0, binding=0 (vert shader only)
			1: surfaceSize.Data(), // layout=0, binding=1 (frag shader only)
		})

	vlk.stats.SegmentDuration[metrics.SegmentPlUpdateGlobalUniform] += time.Since(ts)
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Functions - Groups
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plGroupStats(_ *drawContext, _ *drawGroup) {
	vlk.stats.DrawGroups++
}

func (vlk *VLK) plGroupCreateRenderingPipeline(_ *drawContext, g *drawGroup) {
	cacheKey := g.shader.Meta().ID() + strconv.FormatInt(int64(g.polygonMode), 10)

	if pipe, exist := vlk.drawPipelineCache[cacheKey]; exist {
		g.renderPipe = pipe
		return
	}

	ts := time.Now()

	pipe := vlk.cont.pipelineFactory().NewPipeline(
		pipeline.WithStages([]vulkan.PipelineShaderStageCreateInfo{
			g.shader.ModuleVert().Stage(),
			g.shader.ModuleFrag().Stage(),
		}),
		pipeline.WithTopology(
			g.shader.Meta().Topology(),
			g.shader.Meta().TopologyRestartEnable(),
		),
		pipeline.WithVertexInput(
			g.shader.Meta().Bindings(),
			g.shader.Meta().Attributes(),
		),
		pipeline.WithRasterization(g.polygonMode),
		pipeline.WithColorBlend(),
		pipeline.WithMultisampling(),
	)

	g.renderPipe = pipe
	vlk.drawPipelineCache[cacheKey] = pipe

	vlk.stats.SegmentDuration[metrics.SegmentPlCreatePipeline] += time.Since(ts)
}

func (vlk *VLK) plGroupFindIndexBuffer(_ *drawContext, g *drawGroup) {
	ts := time.Now()

	// find global shader indexes of this shader
	indexes := vlk.indexBufferOf(g.shader)

	// bind to current group
	g.indexes = bufferBinding{
		used:   indexes.Valid,
		buffer: indexes.Buffer,
		offset: indexes.Offset,
	}

	vlk.stats.SegmentDuration[metrics.SegmentPlUpdateIndexes] += time.Since(ts)
}

func (vlk *VLK) plGroupUpdateVertexBuffer(ctx *drawContext, g *drawGroup) {
	ts := time.Now()
	chunks := vlk.cont.allocBuffers().WriteVertexBuffersFromInstances(ctx.currentImageID, g.instances)

	firstInst := uint32(0)
	lastInst := uint32(0)

	for _, chunk := range chunks {
		lastInst = firstInst + chunk.InstanceCount

		g.calls = append(g.calls, &drawCall{
			instances: g.instances[firstInst:lastInst],
			vertexes: bufferBinding{
				used:   true,
				buffer: chunk.Buffer,
				offset: chunk.Offset,
			},
		})

		firstInst += chunk.InstanceCount
	}

	vlk.stats.SegmentDuration[metrics.SegmentPlUpdateVertexes] += time.Since(ts)
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Functions - Exec Groups
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plExecGroupBindPipeline(cb vulkan.CommandBuffer, _ *drawContext, _ *drawSurface, g *drawGroup) {
	ts := time.Now()

	vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, g.renderPipe.Pipeline)

	vlk.stats.SegmentDuration[metrics.SegmentPlBindPipeline] += time.Since(ts)
}

func (vlk *VLK) plExecGroupBindIndexBuffer(cb vulkan.CommandBuffer, _ *drawContext, _ *drawSurface, g *drawGroup) {
	if !g.indexes.used {
		return
	}

	ts := time.Now()

	vulkan.CmdBindIndexBuffer(cb, g.indexes.buffer, g.indexes.offset, vulkan.IndexTypeUint16)

	vlk.stats.SegmentDuration[metrics.SegmentPlBindIndexes] += time.Since(ts)
}

func (vlk *VLK) plExecGroupOnEveryCall(callFns ...drawCallExecFn) drawGroupExecFn {
	return func(cb vulkan.CommandBuffer, ctx *drawContext, surf *drawSurface, g *drawGroup) {
		for _, call := range g.calls {
			for _, fn := range callFns {
				fn(cb, ctx, surf, g, call)
			}
		}
	}
}

// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=
// Functions - Exec Calls
// ~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

func (vlk *VLK) plExecCallUpdateLocalUniforms(_ vulkan.CommandBuffer, ctx *drawContext, _ *drawSurface, _ *drawGroup, c *drawCall) {
	ts := time.Now()
	data := make([]byte, 0, 256)

	for _, inst := range c.instances {
		instData := inst.StorageData()
		if len(instData) == 0 {
			continue
		}

		data = append(data, instData...)
	}

	if len(data) == 0 {
		return
	}

	localUniform := vlk.cont.descriptorsManager().UpdateSet(
		ctx.currentImageID,
		dscptr.LayoutIndexObject,
		map[uint32][]byte{
			0: data, // layout=1, binding=0 (all shaders)
		})

	c.uniforms = append(c.uniforms, localUniform)

	vlk.stats.SegmentDuration[metrics.SegmentPlUpdateSSBO] += time.Since(ts)
}

func (vlk *VLK) plExecCallBindUniforms(cb vulkan.CommandBuffer, _ *drawContext, surf *drawSurface, g *drawGroup, c *drawCall) {
	ts := time.Now()

	// layout = 0, global data
	descriptorSets := []vulkan.DescriptorSet{surf.uniform}

	// layout = 1, call data
	if len(c.uniforms) > 0 {
		descriptorSets = append(descriptorSets, c.uniforms...)
	}

	vulkan.CmdBindDescriptorSets(
		cb,
		vulkan.PipelineBindPointGraphics,
		g.renderPipe.Layout,
		0,
		uint32(len(descriptorSets)),
		descriptorSets,
		0,
		nil,
	)

	vlk.stats.SegmentDuration[metrics.SegmentPlBindUniforms] += time.Since(ts)
}

func (vlk *VLK) plExecCallBindVertexBuffer(cb vulkan.CommandBuffer, _ *drawContext, _ *drawSurface, _ *drawGroup, c *drawCall) {
	if !c.vertexes.used {
		return
	}

	ts := time.Now()

	buffers := []vulkan.Buffer{c.vertexes.buffer}
	offsets := []vulkan.DeviceSize{c.vertexes.offset}
	vulkan.CmdBindVertexBuffers(cb, 0, uint32(len(buffers)), buffers, offsets)

	vlk.stats.SegmentDuration[metrics.SegmentPlBindVertex] += time.Since(ts)
}

func (vlk *VLK) plExecCallInstancedDraw(cb vulkan.CommandBuffer, _ *drawContext, _ *drawSurface, g *drawGroup, c *drawCall) {
	ts := time.Now()

	instanceCount := uint32(len(c.instances))
	indexCount := uint32(len(g.shader.Meta().Indexes()))
	vertexCount := g.shader.Meta().VertexCount()

	instPerCall := min(instanceCount, def.BufferIndexMaxInstances)
	for firstInst := uint32(0); firstInst < instanceCount; firstInst += instPerCall {
		// if we try to draw more instances, that fit in warm index cache (>65536)
		// we split it into chunks of def.BufferIndexMaxInstances size each

		// todo: test def.BufferIndexMaxInstances on small values like 1-2
		vulkan.CmdDrawIndexed(
			cb,
			indexCount*instPerCall,
			instPerCall,
			0,
			int32(firstInst*vertexCount),
			firstInst,
		)

		vlk.stats.DrawCalls++
	}

	vlk.stats.SegmentDuration[metrics.SegmentPlDraw] += time.Since(ts)
}
