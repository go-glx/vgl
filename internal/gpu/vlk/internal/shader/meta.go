package shader

import "github.com/vulkan-go/vulkan"

type Meta struct {
	id   string
	vert []byte
	frag []byte

	topology        vulkan.PrimitiveTopology
	topologyRestart bool
	bindings        []vulkan.VertexInputBindingDescription
	attributes      []vulkan.VertexInputAttributeDescription
	vertexCount     uint32
	useIndexBuffer  bool
	indexes         []uint16
}

func NewMeta(
	id string,
	vert []byte,
	frag []byte,
	topology vulkan.PrimitiveTopology,
	topologyRestart bool,
	bindings []vulkan.VertexInputBindingDescription,
	attributes []vulkan.VertexInputAttributeDescription,
	vertexCount uint32,
	useIndexBuffer bool,
	indexes []uint16,
) *Meta {
	return &Meta{
		id:              id,
		vert:            vert,
		frag:            frag,
		topology:        topology,
		topologyRestart: topologyRestart,
		bindings:        bindings,
		attributes:      attributes,
		vertexCount:     vertexCount,
		useIndexBuffer:  useIndexBuffer,
		indexes:         indexes,
	}
}

func (s *Meta) ID() string {
	return s.id
}

func (s *Meta) Topology() vulkan.PrimitiveTopology {
	return s.topology
}

func (s *Meta) TopologyRestartEnable() bool {
	return s.topologyRestart
}

func (s *Meta) Bindings() []vulkan.VertexInputBindingDescription {
	return s.bindings
}

func (s *Meta) Attributes() []vulkan.VertexInputAttributeDescription {
	return s.attributes
}

func (s *Meta) VertexCount() uint32 {
	return s.vertexCount
}

func (s *Meta) UseIndexBuffer() bool {
	return s.useIndexBuffer
}

func (s *Meta) Indexes() []uint16 {
	return s.indexes
}
