# Arrays & Slices

## 1. Arrays: fixed-size values

An array type is `[N]T`: exactly `N` elements of type `T`. **The length is part of the type**: `[3]int` and `[4]int` are different, incompatible types. `N` must be a constant known at compile time.

```go
var a [5]int                     // [0 0 0 0 0]  — zero value: every element zeroed
b := [3]string{"go", "is", "fun"} // array literal
c := [...]int{10, 20, 30, 40}    // [...] lets the compiler count: type is [4]int
d := [5]int{1: 7, 3: 9}          // indexed literal: [0 7 0 9 0]

len(c)    // 4 — a constant
a[0] = 42 // index from 0 to len-1
```

### Value semantics

An array variable **is** the whole array, not a reference to it. Assignment and passing to a function **copy every element**.

```go
func zero(a [3]int) { a[0] = 0 } // modifies its own copy

x := [3]int{1, 2, 3}
y := x       // full copy
y[0] = 100
zero(x)
fmt.Println(x, y)                 // [1 2 3] [100 2 3]
fmt.Println(x == [3]int{1, 2, 3}) // true — arrays of comparable elements support == and !=
```

| Property | Array `[N]T` |
|---|---|
| Size | Fixed at compile time, part of the type |
| `=` / passing | Copies all elements |
| `==` | Allowed (element-by-element) if `T` is comparable |
| Typical use | Small fixed data: `[3]float64` coordinates, `[256]bool` lookup tables, `[32]byte` hashes; backing store for slices |

In everyday Go you rarely pass arrays around. You use **slices**.

## 2. Slices: a view into an array

A slice type is `[]T` (no length). A slice is a small **descriptor** that points into an underlying (backing) array.

```go
s := []int{1, 2, 3}   // slice literal: creates an array [3]int and a slice over it
m := make([]int, 3, 10) // len 3, cap 10: [0 0 0]
n := make([]int, 3)     // cap defaults to len

len(m) // 3  — elements you can index: m[0] .. m[2]
cap(m) // 10 — room in the backing array, counted from the slice's first element
```

### nil vs empty

```go
var ns []int   // nil slice: zero value, no backing array
es := []int{}  // empty, non-nil slice

ns == nil, es == nil // true, false
len(ns), len(es)     // 0, 0
fmt.Println(ns, es)  // [] []
ns = append(ns, 1)   // append works on nil: [1]
```

- `len`, `cap`, `range`, `append` and slicing all work on a nil slice: **treat nil and empty the same**. Check emptiness with `len(s) == 0`, not `s == nil`.
- Prefer `var s []int` when declaring an empty slice (Go Code Review Comments).
- Slices are **not comparable** with `==` (only `s == nil` is allowed). Compare elements in a loop (or with `slices.Equal`, next lesson).

## 3. Slice expressions

`a[lo:hi]` makes a new slice over elements `lo` up to but **not including** `hi`. Works on arrays, slices and strings.

```go
arr := [6]int{0, 1, 2, 3, 4, 5}
s := arr[1:4]          // [1 2 3]  len = 4-1 = 3, cap = 6-1 = 5
s[0] = 99              // writes through to arr
fmt.Println(arr)       // [0 99 2 3 4 5]
s = s[:cap(s)]         // re-extend up to cap: [99 2 3 4 5]
```

| Expression | Meaning |
|---|---|
| `s[lo:hi]` | elements `lo..hi-1`; `len = hi-lo`, `cap = cap(s)-lo` |
| `s[:hi]` | `lo` defaults to `0` |
| `s[lo:]` | `hi` defaults to `len(s)` |
| `s[:]` | the whole thing (turns an array into a slice) |
| `s[lo:hi:max]` | **full slice expression**: like `s[lo:hi]` but `cap = max-lo` |

Bounds: `0 <= lo <= hi <= max <= cap(s)`. Note `hi` may go **past `len`** up to `cap` (re-slicing); violating the bounds panics.

```go
t := arr[1:3:4] // [99 2]  len 2, cap 3 — cannot see arr[4] or arr[5]
```

The full slice expression limits what an `append` on the result may overwrite (§5).

## 4. The slice header and shared backing arrays

A slice value is a three-word struct: **pointer, length, capacity**.

```
arr := [6]int{0, 1, 2, 3, 4, 5}
s := arr[1:4]

         s (slice header)
        ┌─────┬─────┬─────┐
        │ ptr │ len │ cap │
        │  •  │  3  │  5  │
        └──┼──┴─────┴─────┘
           │
           ▼
arr: ┌───┬───┬───┬───┬───┬───┐
     │ 0 │ 1 │ 2 │ 3 │ 4 │ 5 │
     └───┴───┴───┴───┴───┴───┘
       0   1   2   3   4   5
           └─ len ───┘
           └─────── cap ─────┘
```

- **Copying a slice copies the header only** (`t := s`). Both headers point at the same array: a write through one is visible through the other.
- Re-slicing never copies elements; `s[1:2]`, `s[:0]`, etc. are all views of the same array.
- A small slice of a huge array keeps the **whole** array alive (no garbage collection). Copy the part you need if the rest is no longer used.

## 5. `append`

```go
s = append(s, 4)          // one element
s = append(s, 5, 6, 7)    // several
s = append(s, other...)   // all of another slice
```

How it works:
1. If `len(s)+n <= cap(s)`: writes into the existing backing array, returns a header with a larger `len`. **Same array**.
2. Otherwise: allocates a **new, bigger array**, copies the old elements, appends, returns a header pointing to the **new** array.

That is why you **must assign the result**: the returned header may have a different pointer and always has a different length. `append(s, x)` without assigning as a statement is a compile error (`append(s, x) (value of type []int) is not used`), but `t := append(s, x)` compiles and leaves `s` unchanged.

### Growth

```go
var g []int
prev := -1
for i := range 20 {
	g = append(g, i)
	if cap(g) != prev {
		fmt.Printf("len=%d cap=%d\n", len(g), cap(g))
		prev = cap(g)
	}
}
// One run on Go 1.27:
// len=1 cap=4
// len=5 cap=8
// len=9 cap=16
// len=17 cap=32
```

Capacity roughly doubles for small slices and grows by a smaller factor (≈1.25×) for large ones, rounded up to allocator size classes. The exact numbers are an **implementation detail**: never depend on them. Appending one element at a time is amortized O(1).

If you know the final size, preallocate:

```go
out := make([]int, 0, len(in)) // len 0, room for len(in): appends never reallocate
for _, v := range in {
	out = append(out, v*2)
}
```

### Aliasing surprises

When there is spare capacity, two appends to the **same** slice write to the **same** slot:

```go
base := make([]int, 3, 4)
p := append(base, 1) // fits in cap: writes arr[3] = 1
q := append(base, 2) // also fits: overwrites arr[3] = 2
fmt.Println(p, q)    // [0 0 0 2] [0 0 0 2]  ← p changed too!
```

When there is no spare capacity, `append` reallocates and the results are **independent**:

```go
full := []int{1, 2, 3}   // len 3, cap 3
r := append(full, 4)     // new array
r[0] = -1
fmt.Println(full, r)     // [1 2 3] [-1 2 3 4]
```

Whether a result aliases its input depends on capacity, which you usually don't control. Rules of thumb:
- Append back to the **same variable** (`s = append(s, ...)`) and don't keep old headers around.
- To hand out a sub-slice that callers may append to safely, cap it: `s[lo:hi:hi]` (cap = len, so the next append must copy).
- To get a truly independent slice, **copy** it (§6).

## 6. `copy`

`copy(dst, src)` copies `min(len(dst), len(src))` elements and returns that count. It never grows `dst`. Overlapping slices are fine.

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)
n := copy(dst, src)      // n = 3, dst = [1 2 3]

n = copy(src[1:], src)   // shift right by one: n = 4, src = [1 1 2 3 4]
```

An independent copy of a whole slice:

```go
dup := make([]int, len(s))
copy(dup, s)
// or, in one line: dup := append([]int(nil), s...)  (nil if s is empty)
```

`copy(dst, src)` into a `make([]T, 0)` copies **nothing**: `dst` needs a length, not just a capacity.

## 7. `range` over slices (and arrays)

```go
for i, v := range s { ... } // index and a COPY of the element
for i := range s { ... }    // index only
for _, v := range s { ... } // value only
```

`v` is a **copy**. Assigning to it doesn't change the slice:

```go
nums := []int{1, 2, 3}
for _, v := range nums {
	v *= 10 // changes the copy only
}
fmt.Println(nums) // [1 2 3]

for i := range nums {
	nums[i] *= 10 // write through the index
}
fmt.Println(nums) // [10 20 30]
```

- The range expression is evaluated **once**, before the loop: `len` is fixed at the start, so appending to `s` inside `for ... range s` doesn't extend the loop.
- Ranging over an **array** (not a slice) copies the whole array first; range over `arr[:]` or `&arr` to avoid that for big arrays.

## 8. Deleting and inserting by hand

The next lesson's `slices.Delete` / `slices.Insert` do this for you; knowing the mechanics explains their behaviour.

```go
// Delete element i, keep order (shifts the tail left; reuses the array)
s = append(s[:i], s[i+1:]...)

// Delete element i, order doesn't matter (O(1))
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// Insert v at index i (0 <= i <= len(s))
s = append(s, 0)     // grow by one (may reallocate)
copy(s[i+1:], s[i:]) // shift the tail right
s[i] = v
```

Deleting in place overwrites the original elements: any other slice sharing the array sees the shifted data.

## 9. Filtering in place: the `s[:0]` trick

`s[:0]` is an empty slice with the **same backing array and capacity**. Appending to it overwrites `s` from the front, which is safe because the write position never passes the read position.

```go
fs := []int{1, 2, 3, 4, 5, 6}
out := fs[:0]
for _, v := range fs {
	if v%2 == 0 {
		out = append(out, v) // never allocates
	}
}
fmt.Println(out) // [2 4 6]
fmt.Println(fs)  // [2 4 6 4 5 6]  — tail is stale leftovers
clear(fs[len(out):]) // optional: zero the tail → [2 4 6 0 0 0]
```

Zeroing the tail matters for slices of pointers/large values (lets the garbage collector free them).

## 10. 2D slices

A `[][]T` is a slice of slices. Each row is an independent slice that must be allocated.

```go
grid := make([][]int, 3) // 3 nil rows
for r := range grid {
	grid[r] = make([]int, 4)
}
grid[1][2] = 5 // [[0 0 0 0] [0 0 5 0] [0 0 0 0]]

rows, cols := len(grid), len(grid[0])
```

- Rows may have **different lengths** ("jagged"), e.g. Pascal's triangle.
- For fixed small sizes, an array of arrays works and is a single value: `var board [3][3]rune`.

## 11. Passing slices to functions

The **header is copied** (pass by value, as always), but it points to the **same backing array**.

| Inside the function | Visible to the caller? |
|---|---|
| `s[i] = x` (write an element) | **Yes** — shared array |
| `s = append(s, x)` | **No** — caller's header still has the old `len` (and maybe the old array) |
| `s = s[1:]` | **No** — only the local header changed |

```go
func zeroFirst(s []int) { s[0] = 0 }               // visible
func addOne(s []int)    { s = append(s, 1) }       // lost

zs := []int{1, 2, 3}
zeroFirst(zs)
addOne(zs)
fmt.Println(zs) // [0 2 3]
```

Idiom: a function that changes the length **returns the new slice**, just like `append`:

```go
func addOne(s []int) []int { return append(s, 1) }
zs = addOne(zs)
```

## 12. Variadic recap

A `...T` parameter is a `[]T` inside the function. `s...` passes an existing slice **without copying**: the function gets the same backing array.

```go
func sum(nums ...int) int { ... }

sum(1, 2, 3)
xs := []int{1, 2, 3}
sum(xs...)                        // nums shares xs's array
all := append([]int{0}, xs...)    // append is variadic too: [0 1 2 3]
```

## 13. `clear`

`clear(s)` (built-in, Go 1.21) sets every element up to `len(s)` to its zero value. It does **not** change `len` or `cap`.

```go
cs := []int{1, 2, 3}
clear(cs)               // [0 0 0], len still 3
cs = cs[:0]             // to empty a slice but keep its capacity
```

## Common mistakes

| Mistake | Fix |
|---|---|
| `s2 := append(s, x)` then still using `s`, or appending inside a function and dropping the result | Always `s = append(s, x)` |
| Two `append`s to the same base slice, expecting independent results | Copy first, or cap the base: `base[:len(base):len(base)]` |
| Returning a sub-slice and being surprised when the caller's append overwrites your data | Return `s[lo:hi:hi]` or a copy |
| `s[i]` with `i >= len(s)` → `panic: runtime error: index out of range [5] with length 3` | Check `i < len(s)`; remember `cap` doesn't make indices valid |
| `for _, v := range s { v = ... }` expecting `s` to change | `for i := range s { s[i] = ... }` |
| Function appends to its parameter and the caller doesn't see it | Return the new slice: `s = f(s)` |
| `copy(make([]int, 0, n), src)` copies nothing | `dst := make([]int, len(src))` |
| `if s == nil` to test emptiness | `if len(s) == 0` |
| Comparing slices with `==` (compile error) | Compare in a loop (or `slices.Equal`, next lesson) |
| Forgetting to `make` each row of a 2D slice → nil row → index panic | Allocate every row in a loop |
| Expecting `[3]int` arguments to be modified by a function | Arrays are copied; return the array or use a slice |

## Key takeaways

- Arrays `[N]T` are values: fixed length (part of the type), copied on assignment, comparable with `==`.
- Slices `[]T` are headers `{ptr, len, cap}` over a backing array; copying a slice shares the array.
- `s[lo:hi]` and `s[lo:hi:max]` re-slice without copying; `max` limits capacity.
- `append` may or may not reallocate: **always** `s = append(s, ...)`, and copy when you need independence.
- `copy(dst, src)` copies `min(len)` elements; `dst` needs a length.
- `range` gives a copy of each element; write via the index.
- Delete/insert/filter in place with `append` + `copy` + `s[:0]`.
- Functions see the caller's elements but not changes to its length: return the new slice.
- nil and empty slices behave the same; test with `len(s) == 0`.
