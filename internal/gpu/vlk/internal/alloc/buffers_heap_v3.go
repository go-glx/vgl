package alloc

type (
	heapV3 struct {
		areas []*h3Area
	}
)

const minNodeSize = 128 // (todo: move to def const)

// Buffer Types:
// - index
// - vertex
// - uniform
//
// Layout:
// - Coherent         (must be in host memory)
// - LocalWritable    (must be in device memory, can be rewritten (has mapped memory and coherent space))
// - LocalImmutable   (must be in device memory, write-only, cannot be changed (less space used))
//
// Flags
// - OneTime  (can be overridden with generationID > allocID)
//
// | heap coherent                                          |
//
//             Real memory chunk (some vk buffer)
// | ------------------------------------------------------ |
//

// Api example:
//
// alloc(Uniform, LocalImmutable, size, OneTime & Flag) Allocation
// free(Allocation)
// write(Allocation, []byte(data))

type h3Area struct {
	generation  uint64
	minNodeSize uint32
	capacity    uint32
	size        uint32
	head        *h3Node
}

func NewArea(capacity uint32) *h3Area {
	return &h3Area{
		generation:  0,
		minNodeSize: minNodeSize,
		capacity:    capacity,
		size:        0,
		head:        newH3Node(0, capacity),
	}
}

type h3Node struct {
	generationID uint64
	offset       uint32
	capacity     uint32
	size         uint32
	prev         *h3Node
	next         *h3Node
	nextFree     *h3Node
}

func newH3Node(generation uint64, capacity uint32) *h3Node {
	return &h3Node{
		generationID: generation,
		capacity:     capacity,
		size:         0,
	}
}

func (curr *h3Node) freeSize() uint32 {
	return curr.capacity - curr.size
}

// Claim will create new virtual memory node with size
// and return offset(ID) of created node
// returns false if node not have enough free space
func (h3 *h3Area) Claim(size uint32) (*h3Node, bool) {
	current := h3.head
	for {
		if current == nil {
			return nil, false
		}

		if h3.splitNodeSpace(current, size) {
			return current, true
		}

		current = current.nextFree
	}
}

// This will split one node space into 2 nodes:
// before:  [ ______________ ]
//  after:  [ XXXXXX ][ ____ ]
//              ^        ^
//      current node     |
//              new node with free space
//
// returns false if node not have enough free space
func (h3 *h3Area) splitNodeSpace(curr *h3Node, realSize uint32) bool {
	size := realSize
	if size < h3.minNodeSize {
		size = h3.minNodeSize
	}

	if curr.freeSize() < size {
		return false
	}

	unclaimedSpace := curr.freeSize() - size
	next := curr.next

	// mark this node as claimed
	curr.size = realSize
	curr.capacity = size
	
	// inc area size usage
	h3.size += size

	// create new next node with unclaimed space
	var freeNode *h3Node

	if unclaimedSpace > 0 {
		freeNode = newH3Node(curr.generationID, unclaimedSpace)
		freeNode.next = next
		freeNode.nextFree = curr.nextFree
		freeNode.prev = curr
		freeNode.offset = curr.offset + size

		curr.next = freeNode
		curr.nextFree = freeNode

		if next != nil {
			next.prev = freeNode
		}
	}

	// change pointers
	if curr.prev != nil {
		curr.prev.nextFree = freeNode
	}

	return true
}
