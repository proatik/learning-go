# Find, Insert & Delete in a Linked List

Edit a singly linked list by re-pointing `Next` fields. A `nil` `*Node` is the empty list.

```go
type Node struct {
	Val  int
	Next *Node
}

func Find(head *Node, v int) *Node
func InsertAfter(n *Node, v int)
func Delete(head *Node, v int) *Node
func DeleteAll(head *Node, v int) *Node
```

| Function | Behavior |
|---|---|
| `Find(head, v)` | the **first** node whose `Val` is `v`, or `nil` |
| `InsertAfter(n, v)` | links a new node holding `v` right after `n`; does nothing if `n` is nil |
| `Delete(head, v)` | unlinks the **first** node holding `v`; returns the (possibly new) head |
| `DeleteAll(head, v)` | unlinks **every** node holding `v`; returns the (possibly new) head |

| List | Call | Result |
|---|---|---|
| `1 → 2 → 3` | `Find(head, 2)` | the second node (same pointer, not a copy) |
| `1 → 3` | `InsertAfter(Find(head, 1), 2)` | `1 → 2 → 3` |
| `1 → 2 → 3` | `Delete(head, 1)` | `2 → 3` (head changed) |
| `1 → 2 → 2` | `Delete(head, 2)` | `1 → 2` |
| `1 → 2 → 3` | `Delete(head, 9)` | `1 → 2 → 3` (unchanged) |
| `2 → 1 → 2 → 2` | `DeleteAll(head, 2)` | `1` |
| `2 → 2` | `DeleteAll(head, 2)` | `nil` |

## Rules
- Callers always write `head = Delete(head, v)`: deleting the first node changes the head.
- Don't copy the list: keep the surviving nodes (the tests check pointer identity).
- `Find`, `Delete` and `DeleteAll` on a `nil` list return `nil`.
