package pipeline

import (
	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
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
	cache          *Cache

	defaultPipelineLayout vulkan.PipelineLayout
	createdPipelines      []vulkan.Pipeline
}

func NewFactory(
	logger config.Logger,
	ld *logical.Device,
	swapChain *swapchain.Chain,
	mainRenderPass *renderpass.Pass,
	cache *Cache,
) *Factory {
	factory := &Factory{
		logger:         logger,
		ld:             ld,
		swapChain:      swapChain,
		mainRenderPass: mainRenderPass,
		cache:          cache,
	}

	factory.defaultPipelineLayout = factory.newDefaultPipelineLayout()
	logger.Debug("pipeline factory created")
	return factory
}

func (f *Factory) Free() {
	vulkan.DestroyPipelineLayout(f.ld.Ref(), f.defaultPipelineLayout, nil)

	for _, pipeline := range f.createdPipelines {
		vulkan.DestroyPipeline(f.ld.Ref(), pipeline, nil)
	}

	f.logger.Debug("freed: pipeline factory")
}

func (f *Factory) NewPipeline(opts ...Initializer) vulkan.Pipeline {
	info := vulkan.GraphicsPipelineCreateInfo{
		SType: vulkan.StructureTypeGraphicsPipelineCreateInfo,
	}

	// default opts
	opts = append(opts, f.withDefaultViewport())
	opts = append(opts, f.withDefaultMainRenderPass())
	opts = append(opts, f.withDefaultLayout())

	// build pipeline info
	for _, applyOpt := range opts {
		applyOpt(&info)
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
	return pipeline
}
