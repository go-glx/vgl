package shader

import "github.com/vulkan-go/vulkan"

type InstanceData interface {
	BindingData() []byte
	IndexesCount() int
	PolygonMode() vulkan.PolygonMode
}
