package shader

type InstanceData interface {
	VertexData() []byte
	StorageData() []byte
}
