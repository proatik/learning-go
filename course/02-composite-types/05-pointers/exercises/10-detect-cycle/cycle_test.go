package cycle

import (
	"fmt"
	"testing"
)

// build makes a list of n nodes holding 1..n. If loopTo >= 0, the last
// node's Next points back to node loopTo (0-based), creating a cycle.
// It returns the head and all nodes in order.
func build(n, loopTo int) (*Node, []*Node) {
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = &Node{Val: i + 1}
	}
	for i := 0; i+1 < n; i++ {
		nodes[i].Next = nodes[i+1]
	}
	if n > 0 && loopTo >= 0 {
		nodes[n-1].Next = nodes[loopTo]
	}
	if n == 0 {
		return nil, nil
	}
	return nodes[0], nodes
}

type cycleCase struct {
	n, loopTo int // loopTo < 0: no cycle
}

// describe returns a readable name for a test case.
func describe(c cycleCase) string {
	if c.loopTo < 0 {
		return fmt.Sprintf("%d nodes, no cycle", c.n)
	}
	return fmt.Sprintf("%d nodes, last -> node %d", c.n, c.loopTo)
}

var cases = []cycleCase{
	{0, -1},
	{1, -1},
	{2, -1},
	{5, -1},
	{1, 0},  // self-loop
	{2, 0},  // 1 -> 2 -> 1
	{2, 1},  // 2 points to itself
	{5, 2},  // the README example
	{5, 0},  // whole list is a cycle
	{5, 4},  // tail self-loop
	{10, 3}, // long tail before the cycle
	{7, 6},
	{100, 37},
}

func TestHasCycle(t *testing.T) {
	for _, c := range cases {
		head, _ := build(c.n, c.loopTo)
		want := c.loopTo >= 0
		if got := HasCycle(head); got != want {
			t.Errorf("HasCycle(%s) = %v, want %v", describe(c), got, want)
		}
	}
}

func TestCycleStart(t *testing.T) {
	for _, c := range cases {
		head, nodes := build(c.n, c.loopTo)
		got := CycleStart(head)
		if c.loopTo < 0 {
			if got != nil {
				t.Errorf("CycleStart(%s) = node %d, want nil", describe(c), got.Val)
			}
			continue
		}
		if got != nodes[c.loopTo] {
			t.Errorf("CycleStart(%s) returned the wrong node, want node holding %d", describe(c), nodes[c.loopTo].Val)
		}
	}
}

func TestCycleLen(t *testing.T) {
	for _, c := range cases {
		head, _ := build(c.n, c.loopTo)
		want := 0
		if c.loopTo >= 0 {
			want = c.n - c.loopTo
		}
		if got := CycleLen(head); got != want {
			t.Errorf("CycleLen(%s) = %d, want %d", describe(c), got, want)
		}
	}
}

func TestComparesByIdentity(t *testing.T) {
	// Every node holds the same value: only pointer comparison works.
	nodes := make([]*Node, 6)
	for i := range nodes {
		nodes[i] = &Node{Val: 7}
	}
	for i := 0; i+1 < len(nodes); i++ {
		nodes[i].Next = nodes[i+1]
	}
	if HasCycle(nodes[0]) {
		t.Errorf("HasCycle on 6 nodes all holding 7 (no cycle) = true, want false: compare nodes, not values")
	}
	nodes[5].Next = nodes[1]
	if !HasCycle(nodes[0]) {
		t.Errorf("HasCycle on 6 equal-valued nodes with last -> node 1 = false, want true")
	}
	if got := CycleStart(nodes[0]); got != nodes[1] {
		t.Errorf("CycleStart on equal-valued nodes returned the wrong node, want the second node")
	}
}

func TestDoesNotModify(t *testing.T) {
	head, nodes := build(8, 3)
	before := make([]*Node, len(nodes))
	for i, n := range nodes {
		before[i] = n.Next
	}
	HasCycle(head)
	CycleStart(head)
	CycleLen(head)
	for i, n := range nodes {
		if n.Next != before[i] || n.Val != i+1 {
			t.Fatalf("node %d was modified; the functions must only read the list", i+1)
		}
	}
}
