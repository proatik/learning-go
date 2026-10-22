# Generic Binary Search Tree

A binary search tree for any ordered type. A recursive generic node type (`*node[T]` inside `node[T]`), `cmp.Ordered` for comparisons, and the full set of operations, including delete.

```go
type node[T cmp.Ordered] struct {
	val         T
	left, right *node[T]
}

type Tree[T cmp.Ordered] struct {
	root *node[T]
	size int
}

func From[T cmp.Ordered](vals ...T) *Tree[T]

func (t *Tree[T]) Insert(v T) bool
func (t *Tree[T]) Contains(v T) bool
func (t *Tree[T]) Delete(v T) bool
func (t *Tree[T]) Len() int
func (t *Tree[T]) Height() int
func (t *Tree[T]) Min() (T, bool)
func (t *Tree[T]) Max() (T, bool)
func (t *Tree[T]) Floor(v T) (T, bool)
func (t *Tree[T]) Ceil(v T) (T, bool)
func (t *Tree[T]) InOrder() []T
```

## Examples

| Code (`t := From(5, 3, 8, 1, 4)`) | Result |
|---|---|
| `t.InOrder()` | `[1 3 4 5 8]` |
| `t.Insert(4)` / `t.Insert(6)` | `false` (duplicate) / `true` |
| `t.Contains(8)`, `t.Contains(7)` | `true`, `false` |
| `t.Height()` | `3` (5 → 3 → 1) |
| `t.Min()` / `t.Max()` | `1, true` / `8, true` |
| `t.Floor(7)` / `t.Floor(0)` | `5, true` / `0, false` |
| `t.Ceil(6)` / `t.Ceil(9)` | `8, true` / `0, false` |
| `t.Delete(5)` then `t.InOrder()` | `true`, `[1 3 4 8]` |
| `From("pear", "fig", "apple").InOrder()` | `["apple" "fig" "pear"]` |

## Rules
- The zero `Tree` is empty and usable. `From` inserts the values in order (duplicates ignored).
- Smaller values go **left**, larger go **right**; duplicates are not stored (`Insert` returns `false`).
- Don't rebalance: the tree's shape is decided purely by insertion order (the tests check `Height` and the node layout).
- `Height` counts **nodes** on the longest root-to-leaf path: empty → `0`, one node → `1`.
- `Floor(v)` is the largest value `<= v`; `Ceil(v)` the smallest value `>= v`. When none exists, or the tree is empty, return the zero value of `T` and `false`; same for `Min` / `Max`.
- `Delete` returns whether `v` was present:
  - a node with **0 or 1** child is replaced by that child (or nil);
  - a node with **2** children takes the value of its **in-order successor** (the minimum of its right subtree), and that successor node is deleted from the right subtree.
- `Len` is O(1): keep `size` up to date.
