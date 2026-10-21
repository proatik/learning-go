# Singly Linked List

Build the basic toolkit for a singly linked list. A `nil` `*Node` is the empty list.

```go
type Node struct {
	Val  int
	Next *Node
}

func Push(head *Node, v int) *Node
func FromSlice(vals []int) *Node
func Len(head *Node) int
func ToSlice(head *Node) []int
func Reverse(head *Node) *Node
```

| Function | Behavior |
|---|---|
| `Push(head, v)` | returns a new head: a node holding `v` whose `Next` is `head` |
| `FromSlice(vals)` | builds a list with the values **in the same order**; `nil` for an empty slice |
| `Len(head)` | number of nodes |
| `ToSlice(head)` | the values in order; `nil` for the empty list |
| `Reverse(head)` | reverses the list **in place** and returns the new head |

| Example | Result |
|---|---|
| `Push(Push(nil, 2), 1)` | `1 → 2` |
| `ToSlice(FromSlice([]int{1, 2, 3}))` | `[1 2 3]` |
| `Len(nil)` | `0` |
| `ToSlice(Reverse(FromSlice([]int{1, 2, 3})))` | `[3 2 1]` |
| `Reverse(nil)` | `nil` |

## Rules
- Traverse with `for n := head; n != nil; n = n.Next`.
- `Reverse` must reuse the existing nodes (re-point `Next` fields), not allocate new ones.
- No recursion needed; loops are fine for every function.
