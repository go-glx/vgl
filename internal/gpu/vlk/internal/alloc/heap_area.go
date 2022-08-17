package alloc

import "math"

type h3Area struct {
	generation uint64
	align      uint32
	capacity   uint32
	size       uint32
	head       *h3Node
}

func NewArea(capacity uint32, alignSize uint32) *h3Area {
	space := newH3Node(0, capacity)

	return &h3Area{
		generation: 0,
		align:      alignSize,
		capacity:   capacity,
		size:       0,
		head:       space,
	}
}

type h3Node struct {
	generationID uint64
	offset       uint32
	capacity     uint32
	size         uint32
	next         *h3Node
}

func newH3Node(generation uint64, capacity uint32) *h3Node {
	return &h3Node{
		generationID: generation,
		capacity:     capacity,
		size:         0,
	}
}

func (h3 *h3Area) freeSize() uint32 {
	return h3.capacity - h3.size
}

func (curr *h3Node) freeSize() uint32 {
	return curr.capacity - curr.size
}

// Claim will create new virtual memory node with size
// and return offset(ID) of created node
// returns false if node not have enough free space
func (h3 *h3Area) Claim(size uint32) (*h3Node, bool) {
	if h3.freeSize() < size {
		return nil, false
	}

	return h3.walk(func(node *h3Node) bool {
		return h3.splitNodes(node, size)
	})
}

// Free will find occupied memory node at offset
// and mark it as free
func (h3 *h3Area) Free(offset uint32) (*h3Node, bool) {
	return h3.walk(func(node *h3Node) bool {
		if node.offset != offset {
			return false
		}

		return h3.mergeNodes(node)
	})
}

func (h3 *h3Area) walk(act func(node *h3Node) bool) (exitNode *h3Node, found bool) {
	curr := h3.head

	for curr != nil {
		if found := act(curr); found {
			return curr, true
		}

		curr = curr.next
	}

	return nil, false
}

// This will split one node space into 2 nodes:
//   before:  [PREV][ ______________ ][NEXT]
//    after:  [PREV][ XXXXXX ][ ____ ][NEXT]
//                      ^        ^
//              current node     |
//                               new node with free space
//
// returns false if node not have enough free space
func (h3 *h3Area) splitNodes(node *h3Node, realSize uint32) bool {
	size := uint32(math.Ceil(float64(realSize)/float64(h3.align)) * float64(h3.align))

	if node.freeSize() < size {
		// not enough space in current node
		return false
	}

	// is suitable node, now we can split it into two nodes
	unusedSpaces := node.capacity - size

	// allocate space
	node.size = realSize
	node.capacity = size
	h3.size += size

	if unusedSpaces <= 0 {
		// all node space is used, not need to split it
		return true
	}

	// create new node for unused space
	newRight := newH3Node(node.generationID, unusedSpaces)
	newRight.size = 0
	newRight.capacity = unusedSpaces
	newRight.offset = node.offset + node.capacity
	newRight.next = node.next

	// move ptr
	node.next = newRight

	// ok
	return true
}

// before:  [PREV][ ____ ][ XXXXXX ][ ____ ][NEXT]
//  after:  [PREV][ ______________________ ][NEXT]
func (h3 *h3Area) mergeNodes(node *h3Node) bool {
	// free current node
	node.size = 0
	right := node.next

	// if right node is empty
	// we can merge it with current
	if right != nil && right.size == 0 {
		// extend current
		node.capacity += right.capacity
		node.next = right.next

		// remove right node
		right.next = nil
		right.capacity = 0
	}

	// find prev node, and if free, we can extend to left too
	prevNode, exist := h3.walk(func(sample *h3Node) bool {
		return sample.next == node
	})

	// left is empty too, recursive extend to left
	if exist && prevNode != nil && prevNode.size == 0 {
		h3.mergeNodes(prevNode)
	}

	return true
}
