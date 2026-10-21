# Module 2 Challenge: Modelling Data

No new concepts here. This page is a one-page recap of Module 2; the ten exercises each combine several of these topics to model real data. Behaviour is written as **plain functions** that take structs or pointers to structs — methods come in Module 3.

## 1. Arrays vs slices

| | Array `[N]T` | Slice `[]T` |
|---|---|---|
| Length | part of the type, fixed | dynamic: `len(s)`, capacity `cap(s)` |
| Assign / pass | **copies** every element | copies a 3-word header; **shares** the backing array |
| Zero value | N zero elements | `nil` (`len 0`, `append` works) |
| Comparable `==` | yes (if `T` is) | no (only `== nil`); use `slices.Equal` |
| Map key | yes | no |
| Use for | fixed shapes: `[3][3]rune`, `[9][9]int` | almost everything else |

```go
a := [...]int{1, 2, 3}          // array, length inferred
s := []int{1, 2, 3}             // slice literal
s = make([]int, 0, 10)          // len 0, cap 10
t := a[1:3]                     // slice of array: shares a's memory, len 2, cap 2
u := s[lo:hi:max]               // full slice expression: cap = max-lo
grid := make([][]int, rows)     // 2D: allocate each row
for i := range grid {
	grid[i] = make([]int, cols)
}
```

## 2. append / copy rules

| Rule | Consequence |
|---|---|
| `s = append(s, v...)` — **always** reassign | `append` may return a new backing array |
| Fits in `cap` → writes in place | other slices of the same array **see** the write |
| Exceeds `cap` → allocates and copies | the result no longer aliases the original |
| `n := copy(dst, src)` copies `min(len(dst), len(src))` | `dst` must already have length: `make([]T, len(src))` |
| Independent copy | `slices.Clone(s)` or `append([]T(nil), s...)` |
| Delete from a slice while keeping order | `s = slices.Delete(s, i, i+1)` |

## 3. `slices`, `maps`, `cmp`

| Function | Does |
|---|---|
| `slices.Sort(s)` / `slices.SortFunc(s, cmp)` / `slices.SortStableFunc` | sort in place |
| `slices.BinarySearch(s, v)` → `(i, found)` | search a sorted slice |
| `slices.Index`, `slices.Contains`, `slices.Equal`, `slices.Reverse`, `slices.Max` | everyday helpers |
| `slices.Insert(s, i, v...)`, `slices.Delete(s, i, j)`, `slices.Compact(s)` | return the new slice — reassign |
| `slices.Collect(seq)`, `slices.Sorted(seq)` | iterator → slice |
| `maps.Keys(m)`, `maps.Values(m)` | iterators over a map (random order) |
| `maps.Clone(m)`, `maps.Copy(dst, src)`, `maps.DeleteFunc` | whole-map operations |
| `cmp.Compare(a, b)` → `-1/0/+1` | building block for `SortFunc` |
| `cmp.Or(a, b, ...)` | first non-zero value — chains tie-breakers |

```go
// Sort by count descending, then word ascending.
slices.SortFunc(ws, func(a, b WordCount) int {
	return cmp.Or(
		cmp.Compare(b.Count, a.Count), // reversed operands = descending
		cmp.Compare(a.Word, b.Word),
	)
})
keys := slices.Sorted(maps.Keys(m)) // deterministic iteration over a map
```

## 4. Map idioms

```go
m := map[string]int{}            // or make(map[string]int, hint); a nil map panics on write
v, ok := m[k]                    // comma-ok: missing key → zero value, false
m[w]++                           // counting: missing key starts at 0
delete(m, k)                     // no-op if absent
seen := map[string]bool{}        // set (or map[string]struct{})
groups := map[string][]Item{}    // group-by: groups[k] = append(groups[k], it)
nested := map[string]map[string]int{} // inner map must be created before writing
```

- Iteration order is **random** — sort keys when order matters.
- Map values are not addressable: `m[k].Field = 1` does not compile. Store `map[K]*T`, or read–modify–write.
- Maps are reference-like: passing one to a function lets it modify the caller's map.

## 5. Structs and embedding

```go
type Point struct{ X, Y int }
p := Point{X: 1, Y: 2}            // keyed literal (preferred); omitted fields are zero
q := p                            // copy: changing q does not change p

type Named struct {
	Point                         // embedded: promoted fields n.X, n.Y
	Name string `json:"name"`     // tag read by encoding/json
}
n := Named{Point: Point{1, 2}, Name: "A"}
```

## 6. Pointers

| Rule | Code |
|---|---|
| Take an address / dereference | `p := &v`, `*p = 3` |
| Field access auto-dereferences | `p.X` means `(*p).X` |
| Allocate | `&T{...}`, `new(T)`, `new(expr)` (Go 1.26+: `new(42)` is an `*int` to 42) |
| Zero value is `nil` | check before dereferencing: `if p == nil { ... }` |
| Returning `&local` is safe | escape analysis moves it to the heap |
| Linked structures | `type node struct { val int; next *node }` |

### Value or pointer?

| Situation | Use |
|---|---|
| Small, read-only data (`Point`, `Money`) | value |
| Function must modify the caller's struct | `*T` parameter |
| Identity matters / shared by several owners (an account, a book) | `*T`, often `map[string]*T` |
| Large struct passed often | `*T` |
| Nodes of a list / tree / graph | `*node` |
| Returning data the caller must not mutate | value or a fresh copy |

## 7. Choosing the right data structure

- **Fixed size known at compile time?** Array (`[9][9]int`) — free copies, comparable.
- **Ordered collection that grows?** Slice. Need sorted? `slices.Sort` + `slices.BinarySearch`.
- **Lookup by key / membership / counting?** Map (`map[K]V`, `map[K]bool`, `map[K]int`).
- **Several related fields?** Struct. Many of them? `[]T` for order, `map[ID]*T` for lookup.
- **Need O(1) lookup *and* O(1) reordering?** Map + doubly linked list (the LRU pattern).
- **Grid?** `[][]T` or `[]string`; check bounds before indexing neighbours.
- **Graph?** `map[string][]string` adjacency list; BFS with a slice as queue and a map as visited set.
- **Results must be deterministic?** Never rely on map order — sort before returning.

## Key takeaways

- Slices share memory, arrays copy it; always reassign `append`'s result.
- Maps need `make` (or a literal) before writing and have no order.
- `cmp.Or` + `cmp.Compare` express multi-key sorts in one expression.
- Use pointers when something has identity or must be mutated; values otherwise.
