package shader

import "github.com/vulkan-go/vulkan"

type InstanceData interface {
	BindingData() []byte
	PolygonMode() vulkan.PolygonMode
}
