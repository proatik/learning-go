# Detect a Cycle (Floyd's Tortoise and Hare)

A buggy list may loop back on itself: some node's `Next` points to an **earlier** node, so `for n != nil` never ends. Detect that using **pointer identity** (`==` on `*Node`) and constant extra memory.

```go
type Node struct {
	Val  int
	Next *Node
}

func HasCycle(head *Node) bool
func CycleStart(head *Node) *Node
func CycleLen(head *Node) int
```

```
1 → 2 → 3 → 4 → 5
        ↑       │
        └───────┘      5.Next == node 3
```

| Function | For the list above | For a list without a cycle |
|---|---|---|
| `HasCycle(head)` | `true` | `false` |
| `CycleStart(head)` | node `3` (the first node reached twice) | `nil` |
| `CycleLen(head)` | `3` (nodes `3, 4, 5`) | `0` |

Floyd's algorithm:
1. `slow` moves 1 step, `fast` moves 2 steps. If `fast` reaches `nil`, there is no cycle. If `slow == fast`, they met **inside** the cycle.
2. **Start:** move one pointer back to `head`; advance both 1 step at a time. Where they meet is the cycle's first node.
3. **Length:** from the meeting point, walk the cycle once and count steps until you are back.

## Rules
- Use `O(1)` extra memory: a few pointer variables, **no map or slice of visited nodes**.
- Compare nodes by pointer (`slow == fast`), never by `Val` (values may repeat).
- Don't modify the list.
- `nil` and a single node without a self-loop have no cycle; a node with `n.Next == n` is a cycle of length 1.
