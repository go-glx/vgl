package alloc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_h3Area_Claim(t *testing.T) {
	const heapSize = 1024
	const sizeBig = minNodeSize * 2
	const sizeSmall = minNodeSize / 2

	area := NewArea(heapSize)
	assert.Equal(t, uint32(0), area.size)
	assert.Equal(t, uint32(heapSize), area.capacity)

	// 1. claim some big chunk of data
	ptr, ok := area.Claim(sizeBig)
	assert.True(t, ok)

	testAssertNodeSize(t, ptr, sizeBig, sizeBig)
	assert.NotEqual(t, nil, ptr.next)
	assert.NotEqual(t, nil, ptr.nextFree)

	var emptyPref *h3Node
	assert.Equal(t, emptyPref, ptr.prev)
	assert.Equal(t, uint32(0), ptr.offset)

	// 2. check area
	// now should have 2 nodes: [ sizeBig, buffSize-sizeBig ]
	assert.Equal(t, uint32(sizeBig), area.size)
	assert.Equal(t, uint32(heapSize), area.capacity)

	nodes := testAreaNodes(area)
	assert.Len(t, nodes, 2)
	testAssertNodeSize(t, nodes[0], sizeBig, sizeBig)
	testAssertNodeSize(t, nodes[1], 0, heapSize-sizeBig)

	// 3. check relations
	assert.Equal(t, nodes[1], nodes[0].next)
	assert.Equal(t, nodes[1], nodes[0].nextFree)
	assert.Equal(t, nodes[0], nodes[1].prev)

	// 4. claim tiny buffer < min
	ptr, ok = area.Claim(sizeSmall)
	assert.True(t, ok)

	nodes = testAreaNodes(area)
	assert.Len(t, nodes, 3)
	testAssertNodeSize(t, nodes[0], sizeBig, sizeBig)
	testAssertNodeSize(t, nodes[1], sizeSmall, minNodeSize)
	testAssertNodeSize(t, nodes[2], 0, heapSize-sizeBig-minNodeSize)

	assert.Equal(t, uint32(sizeBig+minNodeSize), area.size) // inc +minNodeSize
	assert.Equal(t, uint32(heapSize), area.capacity)        // not changed

	// 5. check offsets
	assert.Equal(t, uint32(0), nodes[0].offset)
	assert.Equal(t, uint32(sizeBig), nodes[1].offset)
	assert.Equal(t, uint32(sizeBig+minNodeSize), nodes[2].offset)

	// 6. check relations
	assert.Equal(t, nodes[1], nodes[0].next)     // not changed
	assert.Equal(t, nodes[2], nodes[0].nextFree) // changed 1 -> 2
	assert.Equal(t, nodes[0], nodes[1].prev)     // not changed
	assert.Equal(t, nodes[2], nodes[1].next)     // new
	assert.Equal(t, nodes[2], nodes[1].nextFree) // new
	assert.Equal(t, nodes[1], nodes[2].prev)     // new
}

func testAreaWalk(area *h3Area, fn func(ind int, n *h3Node)) {
	curr := area.head
	index := 0

	for curr != nil {
		fn(index, curr)

		curr = curr.next
		index++
	}
}

func testAreaNodes(area *h3Area) []*h3Node {
	nodes := make([]*h3Node, 0, 32)

	testAreaWalk(area, func(ind int, n *h3Node) {
		nodes = append(nodes, n)
	})

	return nodes
}

func testAssertNodeSize(t *testing.T, node *h3Node, size int, capacity int) {
	assert.Equal(t, uint32(size), node.size, fmt.Sprintf("node '%d' size invalid", node.offset))
	assert.Equal(t, uint32(capacity), node.capacity, fmt.Sprintf("node '%d' capacity invalid", node.offset))
}
