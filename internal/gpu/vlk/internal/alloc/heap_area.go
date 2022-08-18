package alloc

import "math"

type h3Area struct {
	align    uint32  // memory align, typically it`s something like "32 bytes". It is also minimum node size
	capacity uint32  // total area capacity (equal to capacity sum of all nodes)
	size     uint32  // total area physical size (aligned) (can be not equal to sum of nodes size)
	head     *h3Node // ptr to first node
}

func newArea(capacity uint32, alignSize uint32) *h3Area {
	return &h3Area{
		align:    alignSize,
		capacity: capacity,
		size:     0,
		head:     newH3Node(capacity),
	}
}

type h3Node struct {
	offset   uint32  // offset from area start (also nodeID for apis)
	capacity uint32  // node capacity (logical size)
	size     uint32  // node physical size (if size is aligned, its will be equal to capacity)
	next     *h3Node // ptr to next node
}

func newH3Node(capacity uint32) *h3Node {
	return &h3Node{
		capacity: capacity,
		size:     0,
	}
}

func (h3 *h3Area) freeSize() uint32 {
	return h3.capacity - h3.size
}

func (curr *h3Node) freeSize() uint32 {
	return curr.capacity - curr.size
}

// claim will create new virtual memory node with size
// and return offset(ID) of created node
// returns false if node not have enough free space
func (h3 *h3Area) claim(size uint32) (*h3Node, bool) {
	if h3.freeSize() < size {
		return nil, false
	}

	// hint: can be optimized to walk only by free nodes
	//       but currently not see reason of it
	return h3.walk(func(node *h3Node) bool {
		return h3.splitNodes(node, size)
	})
}

// free will find occupied memory node at offset
// and mark it as free
func (h3 *h3Area) free(offset uint32) (*h3Node, bool) {
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
	// align physical size
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
	newRight := newH3Node(unusedSpaces)
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
	// update area size, if current node non empty
	if node.size > 0 {
		h3.size -= node.capacity // capacity is aligned size
	}

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
