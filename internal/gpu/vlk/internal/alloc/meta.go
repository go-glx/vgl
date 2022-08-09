package alloc

import "github.com/vulkan-go/vulkan"

// AllocationID is internal VLK pointer to Allocation
// this allows to fragment/relocate memory buffers inside
// and this ID will always point to valid memory buffer
type AllocationID uint32

type Allocation struct {
	HasData bool          // if true - require index binding before drawing
	Buffer  vulkan.Buffer // vulkan buffer for binding
	Range   Range         // memory range with stored data in this buffer
}

type Range struct {
	PositionFrom uint32
	Size         uint32
}
