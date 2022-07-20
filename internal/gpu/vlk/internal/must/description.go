package must

import "github.com/vulkan-go/vulkan"

type (
	codeNames = map[vulkan.Result][2]string
)

// see descriptions/codes here:
const errorCodesUrl = "https://registry.khronos.org/vulkan/specs/1.3-extensions/man/html/VkResult.html"

var codeNameReferences = codeNames{
	vulkan.ErrorOutOfHostMemory: [2]string{
		"ErrorOutOfHostMemory",
		"A host memory allocation has failed.",
	},
	vulkan.ErrorOutOfDeviceMemory: [2]string{
		"ErrorOutOfDeviceMemory",
		"A device memory allocation has failed.",
	},
	vulkan.ErrorInitializationFailed: [2]string{
		"ErrorInitializationFailed",
		"Initialization of an object could not be completed for implementation-specific reasons.",
	},
	vulkan.ErrorDeviceLost: [2]string{
		"ErrorDeviceLost",
		"The logical or physical device has been lost. See Lost Device",
	},
	vulkan.ErrorMemoryMapFailed: [2]string{
		"ErrorMemoryMapFailed",
		"Mapping of a memory object has failed.",
	},
	vulkan.ErrorLayerNotPresent: [2]string{
		"ErrorLayerNotPresent",
		"A requested layer is not present or could not be loaded.",
	},
	vulkan.ErrorExtensionNotPresent: [2]string{
		"ErrorExtensionNotPresent",
		"A requested extension is not supported.",
	},
	vulkan.ErrorFeatureNotPresent: [2]string{
		"ErrorFeatureNotPresent",
		"A requested feature is not supported.",
	},
	vulkan.ErrorIncompatibleDriver: [2]string{
		"ErrorIncompatibleDriver",
		"The requested version of Vulkan is not supported by the driver or is otherwise incompatible for implementation-specific reasons.",
	},
	vulkan.ErrorTooManyObjects: [2]string{
		"ErrorTooManyObjects",
		"Too many objects of the type have already been created.",
	},
	vulkan.ErrorFormatNotSupported: [2]string{
		"ErrorFormatNotSupported",
		"A requested format is not supported on this device.",
	},
	vulkan.ErrorFragmentedPool: [2]string{
		"ErrorFragmentedPool",
		"A pool allocation has failed due to fragmentation of the pool’s memory. This must only be returned if no attempt to allocate host or device memory was made to accommodate the new allocation. This should be returned in preference to VKERROROUTOFPOOLMEMORY,\" but only if the implementation is certain that the pool allocation failure was due to fragmentation.\"",
	},
	vulkan.ErrorOutOfPoolMemory: [2]string{
		"ErrorOutOfPoolMemory",
		"A pool memory allocation has failed. This must only be returned if no attempt to allocate host or device memory was made to accommodate the new allocation. If the failure was definitely due to fragmentation of the pool, VKERRORFRAGMENTEDPOOL \"should be returned instead.\"",
	},
	vulkan.ErrorInvalidExternalHandle: [2]string{
		"ErrorInvalidExternalHandle",
		"An external handle is not a valid handle of the specified type.",
	},
	vulkan.ErrorSurfaceLost: [2]string{
		"ErrorSurfaceLost",
		"A surface is no longer available.",
	},
	vulkan.ErrorNativeWindowInUse: [2]string{
		"ErrorNativeWindowInUse",
		"The requested window is already in use by Vulkan or another API in a manner which prevents it from being used again.",
	},
	vulkan.ErrorOutOfDate: [2]string{
		"ErrorOutOfDate",
		"A surface has changed in such a way that it is no longer compatible with the swapchain, and further presentation requests using the swapchain will fail. Applications must query the new surface properties and recreate their swapchain if they wish to continue presenting to the surface.",
	},
	vulkan.ErrorIncompatibleDisplay: [2]string{
		"ErrorIncompatibleDisplay",
		"The display used by a swapchain does not use the same presentable image layout, or is incompatible in a way that prevents sharing an image.",
	},
	vulkan.ErrorValidationFailed: [2]string{
		"ErrorValidationFailed",
		"",
	},
	vulkan.ErrorInvalidShaderNv: [2]string{
		"ErrorInvalidShaderNv",
		"One or more shaders failed to compile or link. More details are reported back to the application via VKEXTDebugReport \"if enabled.\"",
	},
	vulkan.ErrorInvalidDrmFormatModifierPlaneLayout: [2]string{
		"ErrorInvalidDrmFormatModifierPlaneLayout",
		"",
	},
	vulkan.ErrorFragmentation: [2]string{
		"ErrorFragmentation",
		"A descriptor pool creation has failed due to fragmentation.",
	},
	vulkan.ErrorNotPermitted: [2]string{
		"ErrorNotPermitted",
		"",
	},
}
