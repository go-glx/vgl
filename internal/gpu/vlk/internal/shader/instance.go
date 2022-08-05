package shader

type InstanceData interface {
	BindingData() []byte
	Indexes() []uint16
}
