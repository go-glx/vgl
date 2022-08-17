package alloc

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRelationType string

const (
	testRelationTypeNext testRelationType = "next"
	testRelationTypePrev testRelationType = "prev"
	testRelationTypeFree testRelationType = "free"
)

func Test_h3Area_Claim(t *testing.T) {
	const heapSize = 320
	const alignSize = 32
	const sizeBig = 64
	const sizeSmall = 16

	area := NewArea(heapSize, alignSize)
	assert.Equal(t, uint32(0), area.size)
	assert.Equal(t, uint32(heapSize), area.capacity)

	testVisualizeChanges(t, area,
		[]func(area *h3Area){
			func(curr *h3Area) {},
			func(curr *h3Area) { area.Claim(sizeBig) },
			func(curr *h3Area) { area.Claim(sizeSmall) },
			func(curr *h3Area) { area.Claim(sizeBig) },
			func(curr *h3Area) { area.Claim(150) }, // claim all free space (-10 bytes)
			func(curr *h3Area) {
				// try to allocate more space that exist in area
				allocatedNode, ok := area.Claim(32)
				assert.Nil(t, allocatedNode)
				assert.False(t, ok)
			},
		},
		map[string]string{
			"1": "| ...0 ...0 ...0 ...0 ...0 ...0 ...0 ...0 ...0 ...0 |",
			"2": "| 0000 0000 ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 |",
			"3": "| 0000 0000 11-- ...2 ...2 ...2 ...2 ...2 ...2 ...2 |",
			"4": "| 0000 0000 11-- 2222 2222 ...3 ...3 ...3 ...3 ...3 |",
			"5": "| 0000 0000 11-- 2222 2222 3333 3333 3333 3333 333- |",
			"6": "| 0000 0000 11-- 2222 2222 3333 3333 3333 3333 333- |",
		},
		map[string]string{
			"1": "next | HEAD -> 0",
			"2": "next | HEAD -> 0 1",
			"3": "next | HEAD -> 0 1 2",
			"4": "next | HEAD -> 0 1 2 3",
			"5": "next | HEAD -> 0 1 2 3", // next ptr is nil (no free space)
			"6": "next | HEAD -> 0 1 2 3",
		},
		map[string]string{
			"1": "prev | TAIL -> 0",
			"2": "prev | TAIL -> 1 0",
			"3": "prev | TAIL -> 2 1 0",
			"4": "prev | TAIL -> 3 2 1 0",
			"5": "prev | TAIL -> 3 2 1 0", // 3 is last node
			"6": "prev | TAIL -> 3 2 1 0",
		},
		map[string]string{
			"1": "free | HEAD -> 0",
			"2": "free | HEAD -> 1",
			"3": "free | HEAD -> 2",
			"4": "free | HEAD -> 3",
			"5": "free | HEAD ->", // ptr is nil (no free space and nodes to write)
			"6": "free | HEAD ->",
		},
	)
}

func Test_h3Area_Free(t *testing.T) {
	area := testPrepareTestMemoryLayout(t)

	testVisualizeChanges(t, area,
		[]func(area *h3Area){
			func(curr *h3Area) {},
			func(curr *h3Area) { area.Free(testAreaNodes(curr)[2].offset) },
			func(curr *h3Area) { area.Free(testAreaNodes(curr)[3].offset) },
			func(curr *h3Area) { area.Free(testAreaNodes(curr)[1].offset) },
			func(curr *h3Area) { area.Free(testAreaNodes(curr)[3].offset) },
			func(curr *h3Area) { area.Free(testAreaNodes(curr)[2].offset) },
			func(curr *h3Area) { area.Free(512) }, // outside of memory
			func(curr *h3Area) { area.Free(64) },  // offset of node[1]+32 (invalid)
			func(curr *h3Area) { area.Free(32) },  // offset of node[1] (valid, but free)
		},
		map[string]string{
			"1. initial state":                 "| 000- 1111 2222 2222 3333 3333 44-- 5555 ...6 ...6 |",
			"2. simple free (2)":               "| 000- 1111 ...2 ...2 3333 3333 44-- 5555 ...6 ...6 |",
			"3. merge to left (3)":             "| 000- 1111 ...2 ...2 ...2 ...2 33-- 4444 ...5 ...5 |",
			"4. merge to right (1)":            "| 000- ...1 ...1 ...1 ...1 ...1 22-- 3333 ...4 ...4 |",
			"5. merge to right (to end) (3)":   "| 000- ...1 ...1 ...1 ...1 ...1 22-- ...3 ...3 ...3 |",
			"6. merge in both directions (2)":  "| 000- ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 |",
			"7. invalid offset (out of bound)": "| 000- ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 |",
			"8. invalid offset (inside free)":  "| 000- ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 |",
			"9. invalid offset (already free)": "| 000- ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 ...1 |",
		},
		map[string]string{
			"1": "next | HEAD -> 0 1 2 3 4 5 6",
			"2": "next | HEAD -> 0 1 2 3 4 5 6",
			"3": "next | HEAD -> 0 1 2 3 4 5",
			"4": "next | HEAD -> 0 1 2 3 4",
			"5": "next | HEAD -> 0 1 2 3",
			"6": "next | HEAD -> 0 1",
			"7": "next | HEAD -> 0 1",
			"8": "next | HEAD -> 0 1",
			"9": "next | HEAD -> 0 1",
		},
		map[string]string{
			"1": "prev | TAIL -> 6 5 4 3 2 1 0",
			"2": "prev | TAIL -> 6 5 4 3 2 1 0",
			"3": "prev | TAIL -> 5 4 3 2 1 0",
			"4": "prev | TAIL -> 4 3 2 1 0",
			"5": "prev | TAIL -> 3 2 1 0",
			"6": "prev | TAIL -> 1 0",
			"7": "prev | TAIL -> 1 0",
			"8": "prev | TAIL -> 1 0",
			"9": "prev | TAIL -> 1 0",
		},
		map[string]string{
			"1": "free | HEAD -> 6",
			"2": "free | HEAD -> 2 6",
			"3": "free | HEAD -> 2 5",
			"4": "free | HEAD -> 1 4",
			"5": "free | HEAD -> 1 3",
			"6": "free | HEAD -> 1",
			"7": "free | HEAD -> 1",
			"8": "free | HEAD -> 1",
			"9": "free | HEAD -> 1",
		},
	)
}

// return new area with 320 bytes space (10 blocks by 32 bytes)
// and preallocate some nodes for next testing
//
// this will return area with state
// |---------------------------------------------------|
// | 000- 1111 2222 2222 3333 3333 44-- 5555 ...6 ...6 |
// |---------------------------------------------------|
func testPrepareTestMemoryLayout(t *testing.T) *h3Area {
	area := NewArea(320, 32)

	expectedLayout := "| 000- 1111 2222 2222 3333 3333 44-- 5555 ...6 ...6 |"
	blocks := []uint32{24, 32, 64, 64, 16, 32}

	for _, block := range blocks {
		_, ok := area.Claim(block)
		assert.True(t, ok)
	}

	assert.Equal(t,
		expectedLayout,
		testPrintAreaMemoryLayout(area),
		"unexpected default memory",
	)

	return area
}

func testVisualizeChanges(
	t *testing.T,
	area *h3Area,
	mutate []func(area *h3Area),
	states ...map[string]string,
) {
	assert.NotEmpty(t, states, "expected at least layout changes in states at testVisualizeChanges")
	for _, input := range states {
		assert.Len(t, input, len(mutate), "len of mutate should match len of expectedLayouts")
	}

	// ----------------------------------------------------------
	// 0 - data visualize
	// ----------------------------------------------------------

	var dataMemory [][2]string
	var dataRelationNext [][2]string
	var dataRelationPrev [][2]string
	var dataRelationFree [][2]string

	if len(states) >= 1 {
		dataMemory = testSortInput(states[0])
	}
	if len(states) >= 2 {
		dataRelationNext = testSortInput(states[1])
	}
	if len(states) >= 3 {
		dataRelationPrev = testSortInput(states[2])
	}
	if len(states) >= 4 {
		dataRelationFree = testSortInput(states[3])
	}

	for ind, fn := range mutate {
		var expectedRelNext, expectedRelPref, expectedRelFree string

		if len(dataRelationNext) != 0 {
			expectedRelNext = dataRelationNext[ind][1]
		}
		if len(dataRelationPrev) != 0 {
			expectedRelPref = dataRelationPrev[ind][1]
		}
		if len(dataRelationFree) != 0 {
			expectedRelFree = dataRelationFree[ind][1]
		}

		testVisualizeChange(t, area,
			dataMemory[ind][0], // comment (key)
			dataMemory[ind][1], // expected memory state
			expectedRelNext,    // expected relations "next" state
			expectedRelPref,    // expected relations "prev" state
			expectedRelFree,    // expected relations "free" state
			fn,
		)
	}

	if len(states) == 1 {
		return
	}
}

func testSortInput(in map[string]string) [][2]string {
	out := make([][2]string, 0, len(in))

	for index, value := range in {
		out = append(out, [2]string{index, value})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i][0] <= out[j][0]
	})

	return out
}

func testVisualizeChange(
	t *testing.T,
	area *h3Area,
	comment string,
	expectedMemory string,
	expectedRelationNext string,
	expectedRelationPrev string,
	expectedRelationFree string,
	mutate func(area *h3Area),
) {
	// ------------------

	memBefore := testPrintAreaMemoryLayout(area)
	relNextBefore := testPrintAreaRelations(area, testRelationTypeNext)
	relPrevBefore := testPrintAreaRelations(area, testRelationTypePrev)
	relFreeBefore := testPrintAreaRelations(area, testRelationTypeFree)

	testAreaPropertiesValid(t, area, func(area *h3Area) {
		mutate(area)
	})

	memAfter := testPrintAreaMemoryLayout(area)
	relNextAfter := testPrintAreaRelations(area, testRelationTypeNext)
	relPrevAfter := testPrintAreaRelations(area, testRelationTypePrev)
	relFreeAfter := testPrintAreaRelations(area, testRelationTypeFree)

	// ------------------

	if memAfter != "" && memAfter != expectedMemory {
		t.Log("MEM CHANGES:")
		t.Log("action:" + comment)
		t.Log("BEFORE:" + memBefore)
		t.Log("EXPECT:" + expectedMemory)
		t.Log(" AFTER:" + memAfter)

		assert.Equal(t, expectedMemory, memAfter, "unexpected memory change")
	}

	if expectedRelationNext != "" && relNextAfter != expectedRelationNext {
		t.Log("RELATION 'NEXT' CHANGES:")
		t.Log("action:" + comment)
		t.Log("BEFORE:" + relNextBefore)
		t.Log("EXPECT:" + expectedRelationNext)
		t.Log(" AFTER:" + relNextAfter)

		assert.Equal(t, expectedRelationNext, relNextAfter, "unexpected relation NEXT change")
	}

	if expectedRelationPrev != "" && relPrevAfter != expectedRelationPrev {
		t.Log("RELATION 'PREV' CHANGES:")
		t.Log("action:" + comment)
		t.Log("BEFORE:" + relPrevBefore)
		t.Log("EXPECT:" + expectedRelationPrev)
		t.Log(" AFTER:" + relPrevAfter)

		assert.Equal(t, expectedRelationPrev, relPrevAfter, "unexpected relation PREV change")
	}

	if expectedRelationFree != "" && relFreeAfter != expectedRelationFree {
		t.Log("RELATION 'FREE' CHANGES:")
		t.Log("action:" + comment)
		t.Log("BEFORE:" + relFreeBefore)
		t.Log("EXPECT:" + expectedRelationFree)
		t.Log(" AFTER:" + relFreeAfter)

		assert.Equal(t, expectedRelationFree, relFreeAfter, "unexpected relation FREE change")
	}
}

func testAreaPropertiesValid(t *testing.T, area *h3Area, mutate func(area *h3Area)) {
	prevCapacity := area.capacity
	prevAlign := area.align

	mutate(area)

	// these properties not changed during any mutations:
	// ------------------------
	assert.Equal(t, prevCapacity, area.capacity, "area capacity should not be changed")
	assert.Equal(t, prevAlign, area.align, "area align should not be changed")
	assert.NotNil(t, area.head, "area should have head")

	// test validity of size and capacity between nodes and area
	// ------------------------
	totalSize := uint32(0)
	totalCapacity := uint32(0)

	testAreaWalkForward(area, func(n *h3Node) {
		totalSize += n.size
		totalCapacity += n.capacity
	})

	assert.LessOrEqual(t, totalSize, area.size, "total nodes size should be <= area.Size (not always equal because of align)")
	assert.Equal(t, area.capacity, totalCapacity, "total nodes capacity not equal area capacity")

	// test offsets
	// ------------------------
	nextOffsetAt := uint32(0)

	testAreaWalkForward(area, func(n *h3Node) {
		assert.Equal(t, nextOffsetAt, n.offset, "node offset not in valid position")
		nextOffsetAt += n.capacity
	})
}

// will return string like this:
// | 000- 1111 2222 2222 3333 3333 44-- 5555 ...6 ...6 |
// rules
// - each block has ALIGN bytes
// - nodeIndex start with 0
// - [ 0000 ] - full occupied block #0 (size = capacity)
// - [ ...1 ] - free block #1          (size = 0, capacity = ALIGN)
// - [ 2--- ] - occupied block #2      (size = 25%, capacity = ALIGN)
func testPrintAreaMemoryLayout(area *h3Area) string {
	str := "|"
	align := float64(area.align)
	index := 0

	testAreaWalkForward(area, func(n *h3Node) {
		if index == 10 {
			index = 0
		}

		chunksCount := int(math.Ceil(float64(n.capacity) / align))
		charsCount := float64(n.size) / align
		isFreeNode := n.size == 0

		for chunkID := 0; chunkID < chunksCount; chunkID++ {
			str += " "

			if isFreeNode {
				str += fmt.Sprintf("...%d", index)
				continue
			}

			for ind := 0; ind < 4; ind++ {
				if charsCount > 0 {
					str += fmt.Sprintf("%d", index)
				} else {
					str += "-"
				}

				charsCount -= 0.25
			}
		}

		index++
	})

	str += " |"
	return str
}

// will print area relations in format:
//   <type> | <FROM> -> <id> <id>
// example:
//   next | HEAD -> 0 1 2 3
//
// Where TYPE one of:
//  - next
//  - prev
//  - free
//
// Where FROM one of:
//  - HEAD
//  - TAIL
func testPrintAreaRelations(area *h3Area, relationType testRelationType) string {
	var walkFn func(*h3Area, func(*h3Node))
	dst := fmt.Sprintf("%s | ", relationType)

	switch relationType {
	case testRelationTypeNext:
		dst += "HEAD ->"
		walkFn = testAreaWalkForward
	case testRelationTypeFree:
		dst += "HEAD ->"
		walkFn = testAreaWalkForwardFree
	case testRelationTypePrev:
		dst += "TAIL ->"
		walkFn = testAreaWalkBackward
	}

	nodesID := map[uint32]int{}
	ind := 0
	testAreaWalkForward(area, func(n *h3Node) {
		nodesID[n.offset] = ind
		ind++
	})

	walkFn(area, func(n *h3Node) {
		dst += fmt.Sprintf(" %d", nodesID[n.offset])
	})

	return dst
}

func testAreaWalkForward(area *h3Area, fn func(node *h3Node)) {
	curr := area.head
	for curr != nil {
		fn(curr)
		curr = curr.next
	}
}

func testAreaWalkForwardFree(area *h3Area, fn func(node *h3Node)) {
	testAreaWalkForward(area, func(node *h3Node) {
		if node.size == 0 {
			fn(node)
		}
	})
}

func testAreaWalkBackward(area *h3Area, fn func(n *h3Node)) {
	var last *h3Node

	transferMap := map[uint32]*h3Node{}

	testAreaWalkForward(area, func(node *h3Node) {
		if last == nil || node.offset > last.offset {
			last = node
		}

		if node.next == nil {
			return
		}

		transferMap[node.next.offset] = node
	})

	for last != nil {
		fn(last)

		newLast, exist := transferMap[last.offset]
		if !exist {
			break
		}

		last = newLast
	}
}

func testAreaNodes(area *h3Area) []*h3Node {
	nodes := make([]*h3Node, 0, 32)

	testAreaWalkForward(area, func(n *h3Node) {
		nodes = append(nodes, n)
	})

	return nodes
}
