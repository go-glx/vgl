package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

type (
	codeNames = map[vulkan.Result]string
)

var codeNameReferences = codeNames{
	vulkan.NotReady:                                 "NotReady",
	vulkan.Timeout:                                  "Timeout",
	vulkan.EventSet:                                 "EventSet",
	vulkan.EventReset:                               "EventReset",
	vulkan.Incomplete:                               "Incomplete",
	vulkan.ErrorOutOfHostMemory:                     "ErrorOutOfHostMemory",
	vulkan.ErrorOutOfDeviceMemory:                   "ErrorOutOfDeviceMemory",
	vulkan.ErrorInitializationFailed:                "ErrorInitializationFailed",
	vulkan.ErrorDeviceLost:                          "ErrorDeviceLost",
	vulkan.ErrorMemoryMapFailed:                     "ErrorMemoryMapFailed",
	vulkan.ErrorLayerNotPresent:                     "ErrorLayerNotPresent",
	vulkan.ErrorExtensionNotPresent:                 "ErrorExtensionNotPresent",
	vulkan.ErrorFeatureNotPresent:                   "ErrorFeatureNotPresent",
	vulkan.ErrorIncompatibleDriver:                  "ErrorIncompatibleDriver",
	vulkan.ErrorTooManyObjects:                      "ErrorTooManyObjects",
	vulkan.ErrorFormatNotSupported:                  "ErrorFormatNotSupported",
	vulkan.ErrorFragmentedPool:                      "ErrorFragmentedPool",
	vulkan.ErrorOutOfPoolMemory:                     "ErrorOutOfPoolMemory",
	vulkan.ErrorInvalidExternalHandle:               "ErrorInvalidExternalHandle",
	vulkan.ErrorSurfaceLost:                         "ErrorSurfaceLost",
	vulkan.ErrorNativeWindowInUse:                   "ErrorNativeWindowInUse",
	vulkan.Suboptimal:                               "Suboptimal",
	vulkan.ErrorOutOfDate:                           "ErrorOutOfDate",
	vulkan.ErrorIncompatibleDisplay:                 "ErrorIncompatibleDisplay",
	vulkan.ErrorValidationFailed:                    "ErrorValidationFailed",
	vulkan.ErrorInvalidShaderNv:                     "ErrorInvalidShaderNv",
	vulkan.ErrorInvalidDrmFormatModifierPlaneLayout: "ErrorInvalidDrmFormatModifierPlaneLayout",
	vulkan.ErrorFragmentation:                       "ErrorFragmentation",
	vulkan.ErrorNotPermitted:                        "ErrorNotPermitted",
}

func vkAssert(result vulkan.Result, err error) {
	if result == vulkan.Success {
		return
	}

	possibleVkError := "unknown"
	if ref, ok := codeNameReferences[result]; ok {
		possibleVkError = ref
	}

	panic(fmt.Errorf("%s: (possible: %s)", err, possibleVkError))
}
