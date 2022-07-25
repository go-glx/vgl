package shader

import "github.com/vulkan-go/vulkan"

type Meta struct {
	id   string
	vert []byte
	frag []byte

	topology   vulkan.PrimitiveTopology
	bindings   []vulkan.VertexInputBindingDescription
	attributes []vulkan.VertexInputAttributeDescription
}

func NewMeta(
	id string,
	vert []byte,
	frag []byte,
	topology vulkan.PrimitiveTopology,
	bindings []vulkan.VertexInputBindingDescription,
	attributes []vulkan.VertexInputAttributeDescription,
) *Meta {
	return &Meta{
		id:         id,
		vert:       vert,
		frag:       frag,
		topology:   topology,
		bindings:   bindings,
		attributes: attributes,
	}
}

func (s *Meta) ID() string {
	return s.id
}

func (s *Meta) Topology() vulkan.PrimitiveTopology {
	return s.topology
}

func (s *Meta) Bindings() []vulkan.VertexInputBindingDescription {
	return s.bindings
}

func (s *Meta) Attributes() []vulkan.VertexInputAttributeDescription {
	return s.attributes
}
