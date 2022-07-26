package instance

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/config"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/def"
	"github.com/go-glx/vgl/internal/gpu/vlk/internal/must"
)

type Instance struct {
	logger config.Logger
	ref    vulkan.Instance
}

func NewInstance(opt CreateOptions) *Instance {
	return &Instance{
		logger: opt.logger,
		ref:    createVk(opt),
	}
}

func (inst *Instance) Free() {
	vulkan.DestroyInstance(inst.ref, nil)

	inst.logger.Debug("freed: vulkan instance")
}

func (inst *Instance) Ref() vulkan.Instance {
	return inst.ref
}

func createVk(opt CreateOptions) vulkan.Instance {
	opt.logger.Info(fmt.Sprintf("init '%s' engine, required extensions: [%v]", opt.engineName, opt.requiredExtensions))

	info := createInfo(opt)

	var inst vulkan.Instance
	must.Work(vulkan.CreateInstance(&info, nil, &inst))

	return inst
}

func createInfo(opt CreateOptions) vulkan.InstanceCreateInfo {
	info := vulkan.InstanceCreateInfo{
		SType: vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vulkan.ApplicationInfo{
			SType:              vulkan.StructureTypeApplicationInfo,
			PApplicationName:   opt.appName,
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			PEngineName:        opt.engineName,
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			ApiVersion:         def.VKApiVersion,
		},
	}

	// setup extensions
	availableExt := fetchAvailableExtensions(opt.logger)
	assertRequiredExtensionsIsAvailable(availableExt, opt.requiredExtensions)
	info.PpEnabledExtensionNames = opt.requiredExtensions
	info.EnabledExtensionCount = uint32(len(info.PpEnabledExtensionNames))

	// setup validation (debug)
	validationLayers := validationLayers(opt)
	info.EnabledLayerCount = uint32(len(validationLayers))
	info.PpEnabledLayerNames = validationLayers

	return info
}
