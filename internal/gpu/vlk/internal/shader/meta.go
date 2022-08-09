package shader

import "github.com/vulkan-go/vulkan"

type Meta struct {
	id   string
	vert []byte
	frag []byte

	topology   vulkan.PrimitiveTopology
	bindings   []vulkan.VertexInputBindingDescription
	attributes []vulkan.VertexInputAttributeDescription
	indexes    []uint16
}

func NewMeta(
	id string,
	vert []byte,
	frag []byte,
	topology vulkan.PrimitiveTopology,
	bindings []vulkan.VertexInputBindingDescription,
	attributes []vulkan.VertexInputAttributeDescription,
	indexes []uint16,
) *Meta {
	return &Meta{
		id:         id,
		vert:       vert,
		frag:       frag,
		topology:   topology,
		bindings:   bindings,
		attributes: attributes,
		indexes:    indexes,
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

func (s *Meta) Indexes() []uint16 {
	return s.indexes
}
