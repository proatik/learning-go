# Binary Search Tree

A binary search tree (BST) keeps, for every node, **smaller keys in the left subtree and larger keys in the right subtree**. A `nil` `*Tree` is the empty tree.

```go
type Tree struct {
	Key         int
	Left, Right *Tree
}

func Insert(t *Tree, key int) *Tree
func Contains(t *Tree, key int) bool
func InOrder(t *Tree) []int
func Height(t *Tree) int
func Min(t *Tree) (int, bool)
```

| Function | Behavior |
|---|---|
| `Insert(t, key)` | adds `key` as a new leaf in the right place and returns the root; a duplicate key is ignored |
| `Contains(t, key)` | reports whether `key` is in the tree, visiting only one path (left or right at each node) |
| `InOrder(t)` | keys in **in-order** (left subtree, node, right subtree), which is ascending; `nil` for the empty tree |
| `Height(t)` | nodes on the longest root-to-leaf path: `0` for nil, `1` for a single node |
| `Min(t)` | the smallest key and `true`; `0, false` for the empty tree |

```go
var root *Tree
for _, k := range []int{5, 3, 8, 1, 4} {
	root = Insert(root, k)
}
//       5
//      / \
//     3   8
//    / \
//   1   4
InOrder(root)     // [1 3 4 5 8]
Height(root)      // 3
Contains(root, 4) // true
Min(root)         // 1, true
```

## Rules
- Callers write `root = Insert(root, k)`: inserting into an empty tree creates the root.
- Insert must not rebuild or rebalance: existing nodes stay where they are.
- Recursion is the natural fit: handle `nil` first, then recurse into `Left` or `Right`.
