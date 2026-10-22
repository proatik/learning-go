# Tree Walk

Write **recursive** iterators over a generic binary search tree. The hard part: when the loop body stops, the stop signal must travel back up through every level of recursion.

```go
type Tree[K cmp.Ordered] struct {
	root *node[K]
	size int
}

type node[K cmp.Ordered] struct {
	key         K
	left, right *node[K] // left: smaller keys, right: larger keys
}

func (t *Tree[K]) Insert(k K) bool // provided: false if k was already present
func (t *Tree[K]) Len() int        // provided

func (t *Tree[K]) All() iter.Seq[K]          // keys in ascending order
func (t *Tree[K]) Backward() iter.Seq[K]     // keys in descending order
func (t *Tree[K]) Range(lo, hi K) iter.Seq[K] // keys with lo <= k <= hi, ascending
```

## Examples

Insert `5 3 8 1 4 9` into `t`:

| Code | Result |
|---|---|
| `slices.Collect(t.All())` | `[1 3 4 5 8 9]` |
| `slices.Collect(t.Backward())` | `[9 8 5 4 3 1]` |
| `slices.Collect(t.Range(3, 8))` | `[3 4 5 8]` |
| `slices.Collect(t.Range(6, 7))` | `[]` |
| `slices.Collect(t.Range(8, 3))` | `[]` (`lo > hi`) |
| `for k := range t.All() { if k == 3 { break } }` | Visits `1`, `3`, then stops without a panic |

## Rules
- In-order walk: left subtree, node, right subtree. `Backward` mirrors it: right, node, left.
- An empty tree (the zero value `Tree[K]{}`) yields nothing.
- `Range` includes both ends. It should skip subtrees that can't contain keys in range: don't go left when `n.key < lo`, don't go right when `n.key > hi`.
- When `yield` returns `false`, stop the **whole** walk: no more `yield` calls at any level.
- Recursion is the natural fit: write a helper like `func walk(n *node[K], yield func(K) bool) bool` that returns `false` once the walk must stop.
