package pipeline

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/descriptors"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/logical"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/renderpass"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/swapchain"
)

type Factory struct {
	logger         config.Logger
	ld             *logical.Device
	swapChain      *swapchain.Chain
	mainRenderPass *renderpass.Pass
	blueprint      *descriptors.Blueprint
	cache          *Cache

	pipelineLayoutOnlyGlobal vulkan.PipelineLayout
	createdPipelines         []vulkan.Pipeline
}

type Info struct {
	Pipeline vulkan.Pipeline
	Layout   vulkan.PipelineLayout
}

func NewFactory(
	logger config.Logger,
	ld *logical.Device,
	swapChain *swapchain.Chain,
	mainRenderPass *renderpass.Pass,
	blueprint *descriptors.Blueprint,
	cache *Cache,
) *Factory {
	factory := &Factory{
		logger:         logger,
		ld:             ld,
		swapChain:      swapChain,
		mainRenderPass: mainRenderPass,
		blueprint:      blueprint,
		cache:          cache,
	}

	factory.pipelineLayoutOnlyGlobal = factory.newPipelineLayoutOnlyGlobal()

	logger.Debug("pipeline factory created")
	return factory
}

func (f *Factory) Free() {
	vulkan.DestroyPipelineLayout(f.ld.Ref(), f.pipelineLayoutOnlyGlobal, nil)

	for _, pipeline := range f.createdPipelines {
		vulkan.DestroyPipeline(f.ld.Ref(), pipeline, nil)
	}

	f.logger.Debug("freed: pipeline factory")
}

func (f *Factory) NewPipeline(opts ...Initializer) Info {
	info := vulkan.GraphicsPipelineCreateInfo{
		SType: vulkan.StructureTypeGraphicsPipelineCreateInfo,
	}

	// default opts
	opts = append(opts, withDefaultViewport())
	opts = append(opts, withDefaultMainRenderPass())

	// build pipeline info
	for _, applyOpt := range opts {
		applyOpt(&info, f)
	}

	// create pipeline from it
	pipelines := make([]vulkan.Pipeline, 1)
	result := vulkan.CreateGraphicsPipelines(
		f.ld.Ref(),
		f.cache.Ref(),
		1,
		[]vulkan.GraphicsPipelineCreateInfo{info},
		nil,
		pipelines,
	)

	must.Work(result)

	pipeline := pipelines[0]
	f.createdPipelines = append(f.createdPipelines, pipeline)

	return Info{
		Pipeline: pipeline,
		Layout:   info.Layout,
	}
}

func (f *Factory) newPipelineLayoutOnlyGlobal() vulkan.PipelineLayout {
	blueprint := f.blueprint.LayoutGlobal()

	info := &vulkan.PipelineLayoutCreateInfo{
		SType:                  vulkan.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount:         1,
		PSetLayouts:            []vulkan.DescriptorSetLayout{blueprint.Ref()},
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}

	var pipelineLayout vulkan.PipelineLayout
	must.Work(vulkan.CreatePipelineLayout(f.ld.Ref(), info, nil, &pipelineLayout))

	return pipelineLayout
}
