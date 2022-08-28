package shader

import "github.com/vulkan-go/vulkan"

type InstanceData interface {
	VertexData() []byte
	StorageData() []byte
	PolygonMode() vulkan.PolygonMode
}
