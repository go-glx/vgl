package shader

import "github.com/vulkan-go/vulkan"

type InstanceData interface {
	BindingData() []byte
	Indexes() []uint16
	PolygonMode() vulkan.PolygonMode
}
