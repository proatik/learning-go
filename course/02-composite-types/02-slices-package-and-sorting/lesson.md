# The slices Package & Sorting

## 1. Why a package?

Before Go 1.21, "does this slice contain x?" meant writing a loop every time. Go 1.18 added **generics** (type parameters), and Go 1.21 added the standard packages `slices` and `cmp`: helpers that work on a slice of **any** element type.

```go
import "slices"

slices.Contains([]int{1, 2, 3}, 2)        // true
slices.Contains([]string{"a", "b"}, "z")  // false
```

The signature in the docs looks scary, but you don't have to write it:

```go
func Contains[S ~[]E, E comparable](s S, v E) bool
```

The part in `[...]` lists **type parameters**. The compiler **infers** them from your arguments (`S = []int`, `E = int`), so you call it like any other function. Writing your own generic functions comes in a later module.

| Constraint in the docs | Means `E` must be… | Examples |
|---|---|---|
| `any` | anything | all types |
| `comparable` | usable with `==` | numbers, strings, bools, arrays of those |
| `cmp.Ordered` | usable with `<` | integers, floats, strings |

## 2. The `slices` toolbox

| Function | Does | Mutates `s`? |
|---|---|---|
| `Contains(s, v)` | Is `v` in `s`? | no |
| `Index(s, v)` | First index of `v`, or `-1` | no |
| `ContainsFunc(s, f)` | Does any element satisfy `f`? | no |
| `IndexFunc(s, f)` | First index where `f(e)` is true, or `-1` | no |
| `Equal(a, b)` | Same length and same elements in order? | no |
| `Compare(a, b)` | Lexicographic compare: `-1`, `0`, `+1` | no |
| `Min(s)`, `Max(s)` | Smallest / largest element (**panics if empty**) | no |
| `Reverse(s)` | Reverses in place | **yes** |
| `Sort(s)` | Sorts ascending in place | **yes** |
| `SortFunc(s, cmp)` | Sorts by your comparison function (not stable) | **yes** |
| `SortStableFunc(s, cmp)` | Like `SortFunc`, but equal elements keep their order | **yes** |
| `IsSorted(s)` | Is `s` in ascending order? | no |
| `BinarySearch(s, v)` | Position of `v` in a **sorted** slice + found? | no |
| `Insert(s, i, vs...)` | Inserts values at `i`, **returns** the new slice | may |
| `Delete(s, i, j)` | Removes `s[i:j]`, **returns** the new slice | **yes** |
| `DeleteFunc(s, f)` | Removes every element where `f(e)` is true, **returns** the new slice | **yes** |
| `Compact(s)` | Collapses runs of equal **adjacent** elements, **returns** the new slice | **yes** |
| `Clone(s)` | Independent copy (shallow) | no |
| `Concat(ss...)` | New slice joining several slices | no |
| `Repeat(s, n)` | New slice: `s` repeated `n` times | no |
| `Chunk(s, n)` | Iterator over sub-slices of up to `n` elements (§9) | no |

Many have a `Func` variant taking your own function: `IndexFunc`, `ContainsFunc`, `EqualFunc`, `MinFunc`, `MaxFunc`, `CompactFunc`, `BinarySearchFunc`. Run `go doc slices` for the full list.

## 3. Searching and comparing

```go
words := []string{"go", "rust", "zig", "go"}

slices.Contains(words, "zig") // true
slices.Index(words, "go")     // 0 (the first match)
slices.Index(words, "c")      // -1

long := func(s string) bool { return len(s) > 3 }
slices.IndexFunc(words, long)    // 1 ("rust")
slices.ContainsFunc(words, long) // true

slices.Equal([]int{1, 2}, []int{1, 2})          // true
slices.Equal([]int{}, nil)                      // true: both have length 0
slices.Compare([]int{1, 2, 3}, []int{1, 3})     // -1: first difference 2 < 3
slices.Compare([]int{1, 2}, []int{1, 2, 0})     // -1: a prefix is smaller
slices.Min([]int{4, 1, 7})                      // 1
slices.Max([]int{4, 1, 7})                      // 7
```

- `==` doesn't work on slices (only `s == nil`). Use `slices.Equal`.
- `Contains`/`Index` are **linear** (O(n)): fine for small slices; for many lookups, sort + `BinarySearch` (§7) or use a map (next lesson).
- Guard `Min`/`Max`: `if len(s) == 0 { ... }` first.

## 4. Sorting: `Sort`, `IsSorted`, `Reverse`

```go
nums := []int{5, 2, 8, 1}
slices.Sort(nums)              // nums is now [1 2 5 8]
slices.IsSorted(nums)          // true
slices.Reverse(nums)           // [8 5 2 1]: descending, the simple way

names := []string{"bob", "Alice", "carol"}
slices.Sort(names)             // [Alice bob carol]: byte order, 'A' (65) < 'b' (98)
```

- Sorting is **in place**: no new slice, nothing returned. The caller's backing array changes (aliasing, L1).
- To keep the original intact, sort a copy: `sorted := slices.Clone(s); slices.Sort(sorted)`.
- Strings sort by **bytes**, so all uppercase letters come before lowercase ones.

## 5. Custom order: `SortFunc` and comparison functions

A **comparison function** `cmp(a, b E) int` returns:

| Result | Meaning |
|---|---|
| negative | `a` comes **before** `b` |
| `0` | equal for sorting purposes |
| positive | `a` comes **after** `b` |

Don't hand-write `if a < b { return -1 }...`. Use the helpers:

| Helper | Returns |
|---|---|
| `cmp.Compare(x, y)` | `-1`, `0`, `+1` for any ordered type (ints, floats, strings) |
| `strings.Compare(a, b)` | `-1`, `0`, `+1` for strings |
| `cmp.Or(v1, v2, ...)` | the first argument that is **not the zero value** |

```go
import (
	"cmp"
	"slices"
	"strings"
)

// Descending: swap the arguments.
d := []int{3, 1, 2}
slices.SortFunc(d, func(a, b int) int { return cmp.Compare(b, a) }) // [3 2 1]

// Two keys: by length, then alphabetically.
ws := []string{"pear", "fig", "banana", "kiwi", "apple"}
slices.SortFunc(ws, func(a, b string) int {
	return cmp.Or(
		cmp.Compare(len(a), len(b)), // primary key
		strings.Compare(a, b),       // tie-breaker: only used when lengths are equal
	)
})
// [fig kiwi pear apple banana]
```

Why `cmp.Or` works: a comparison returns `0` for "tied". `cmp.Or` skips zeros, so it returns the first key that **decides** the order. Add as many keys as you need; flip one key's arguments to make **that key** descending:

```go
// score descending, then name ascending
cmp.Or(cmp.Compare(scoreB, scoreA), strings.Compare(nameA, nameB))
```

> All arguments of `cmp.Or` are evaluated (it's a normal function call). That's fine for cheap comparisons like these.

**Rule:** the function must be a **strict weak ordering**: consistent (`cmp(a, b) < 0` ⇔ `cmp(b, a) > 0`) and transitive. Breaking it (e.g. returning `-1` for both orders, or random results) gives a garbled order, not an error.

## 6. Stability

`Sort` and `SortFunc` are **not stable**: elements that compare equal may end up in any order. `SortStableFunc` keeps equal elements in their **original** order.

```go
// [name, team] pairs, sort by team only
people := [][2]string{{"ann", "B"}, {"bob", "A"}, {"cat", "B"}, {"dan", "A"}}
slices.SortStableFunc(people, func(a, b [2]string) int {
	return strings.Compare(a[1], b[1])
})
// [[bob A] [dan A] [ann B] [cat B]]: within each team, input order is kept
```

| Use | When |
|---|---|
| `SortFunc` + tie-breaker keys | You can fully define the order (usually best: deterministic **and** fast) |
| `SortStableFunc` | The existing order is meaningful and should survive ties |

## 7. Binary search

`BinarySearch` finds a value in a **sorted** slice in O(log n). It returns the index where the value is (or **would be inserted**) and whether it was found.

```go
sorted := []int{10, 20, 30, 40}
slices.BinarySearch(sorted, 30) // 2, true
slices.BinarySearch(sorted, 25) // 2, false: 25 belongs at index 2

// Insert while keeping the slice sorted:
i, found := slices.BinarySearch(sorted, 25)
if !found {
	sorted = slices.Insert(sorted, i, 25) // [10 20 25 30 40]
}
```

On an unsorted slice the result is meaningless: `BinarySearch([]int{3, 1, 2}, 3)` returns `3, false` although `3` is right there.

## 8. Modifying: `Insert`, `Delete`, `DeleteFunc`, `Compact`

These change length, so, like `append`, they **return the new slice**. Always assign it back.

```go
s := []int{0, 1, 2, 3, 4}
s = slices.Insert(s, 2, 90, 91) // [0 1 90 91 2 3 4]
s = slices.Delete(s, 1, 3)      // [0 91 2 3 4]: removes s[1:3]
isEven := func(n int) bool { return n%2 == 0 }
s = slices.DeleteFunc(s, isEven) // [91 3]

c := []int{1, 1, 2, 2, 2, 3, 1}
c = slices.Compact(c) // [1 2 3 1]: only ADJACENT duplicates
```

`Compact` only merges neighbours. For **unique** values: sort first, then compact.

```go
u := slices.Compact(slices.Sorted(slices.Values(xs))) // sorted unique copy of xs (see §9)
```

**They reuse the original backing array.** Since Go 1.22, `Delete`, `DeleteFunc`, `Compact` (and friends) **zero the elements** between the new length and the old length:

```go
orig := []int{1, 2, 3, 4, 5}
shorter := slices.Delete(orig, 1, 3)
fmt.Println(shorter, orig) // [1 4 5] [1 4 5 0 0]
```

So after `s2 := slices.Delete(s, ...)`, treat `s` as **dead**. If you need the original, `Clone` it first.

`go vet` reports a dropped result: `slices.Insert(x, 0, 0)` alone gives *"result of slices.Insert call not used"*.

## 9. Building: `Clone`, `Concat`, `Repeat`

```go
a := []int{1, 2}
b := slices.Clone(a)
b[0] = 99                                   // a is still [1 2]

slices.Concat([]int{1}, []int{2, 3}, nil, []int{4}) // [1 2 3 4]: always a new slice
slices.Repeat([]string{"ab"}, 3)                    // [ab ab ab]
```

- `Clone(nil)` returns `nil`; otherwise a new backing array with the same elements.
- The copy is **shallow**: for `[][]int`, the inner slices are still shared.

## 10. Iterators (Go 1.23): `range` over functions

An **iterator** is a function that produces values one at a time. `range` works on two iterator types from the `iter` package:

| Type | `range` gives | Example source |
|---|---|---|
| `iter.Seq[V]` | one value: `for v := range seq` | `slices.Values(s)`, `slices.Chunk(s, n)` |
| `iter.Seq2[K, V]` | two values: `for k, v := range seq` | `slices.All(s)`, `slices.Backward(s)` |

```go
letters := []string{"a", "b", "c"}

for i, v := range slices.All(letters) { ... }      // 0 a, 1 b, 2 c (same as range letters)
for v := range slices.Values(letters) { ... }      // a, b, c
for i, v := range slices.Backward(letters) { ... } // 2 c, 1 b, 0 a: no index arithmetic

for chunk := range slices.Chunk([]int{1, 2, 3, 4, 5}, 2) {
	fmt.Print(chunk, " ")                          // [1 2] [3 4] [5]
}
```

Iterators connect producers and consumers. Functions that **consume** a sequence and build a slice:

| Function | Returns |
|---|---|
| `slices.Collect(seq)` | a new slice with every value, in order |
| `slices.Sorted(seq)` | a new, **sorted** slice |
| `slices.SortedFunc(seq, cmp)` | a new slice sorted by `cmp` |
| `slices.AppendSeq(s, seq)` | `s` with every value appended |

```go
slices.Collect(slices.Values([]int{3, 1, 2})) // [3 1 2]: a copy
slices.Sorted(slices.Values([]int{3, 1, 2}))  // [1 2 3]: sorted copy, input untouched
```

- `break` inside the loop works as usual: the iterator stops producing.
- Iterators are **lazy**: no slice is built until you range over them or `Collect` them.
- More iterators exist in the standard library, e.g. `strings.SplitSeq(s, ",")` and `strings.FieldsSeq(s)` yield the parts one by one.
- Here you only **use** iterators; writing your own (`func(yield func(V) bool)`) comes later.

## 11. Legacy: the `sort` package

Code written before Go 1.21 uses package `sort`. You will read it often:

| Old (`sort`) | New (`slices`) |
|---|---|
| `sort.Ints(s)`, `sort.Strings(s)`, `sort.Float64s(s)` | `slices.Sort(s)` |
| `sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })` | `slices.SortFunc(s, func(a, b T) int { return cmp.Compare(a, b) })` |
| `sort.SliceStable(...)` | `slices.SortStableFunc(...)` |
| `sort.SearchInts(s, v)` | `slices.BinarySearch(s, v)` |

`sort.Slice` takes a **less** function on **indices** (`bool`), while `slices.SortFunc` takes a **compare** function on **elements** (`int`). The `sort` docs point to the `slices` versions: `sort.Ints` "simply calls `slices.Sort`" since Go 1.22, and `slices.SortFunc` "is more ergonomic and runs faster" than `sort.Slice`. Write new code with `slices`.

## Common mistakes

| Mistake | Fix |
|---|---|
| `slices.BinarySearch` on an unsorted slice | Sort first (or check `slices.IsSorted` in tests) |
| `slices.Insert(s, i, v)` / `slices.Delete(s, i, j)` without assigning | `s = slices.Insert(s, i, v)`; `go vet` catches the unused result |
| Using the old slice after `Delete`/`Compact` | Its tail is zeroed; use only the returned slice, or `Clone` first |
| `Compact` to remove all duplicates | It only merges **adjacent** ones: sort first |
| Sorting a parameter the caller still needs | `c := slices.Clone(s)` then sort `c` |
| Comparator `func(a, b int) int { return -1 }` or one that's inconsistent | Build it from `cmp.Compare` / `strings.Compare` / `cmp.Or` |
| Descending via `cmp.Compare(a, b) * -1` with extra keys tangled | Swap the arguments of just that key: `cmp.Compare(b, a)` |
| `slices.Max(s)` on an empty slice | Panics: check `len(s) == 0` first |
| `slices.Delete(s, i, i+1)` inside `for i, v := range s` | Indices shift and elements get skipped: `for` over `{1, 2, 2, 3}` deleting 2s leaves `[1 2 3]`. Use `slices.DeleteFunc` |
| `a == b` on slices | `slices.Equal(a, b)` |
| Reaching for `sort.Slice` in new code | `slices.SortFunc` |

## Key takeaways

- `slices` (Go 1.21+) gives generic helpers for any slice; type parameters are **inferred**, so call them like normal functions.
- `Sort`, `SortFunc`, `Reverse` mutate in place. `Clone` first when the input must stay untouched.
- Comparators return negative/zero/positive. Build them with `cmp.Compare`, `strings.Compare`, and chain keys with `cmp.Or`; swap arguments for descending.
- `SortFunc` is not stable; use tie-breaker keys or `SortStableFunc`.
- `BinarySearch` needs sorted input and returns the insertion point too.
- `Insert`, `Delete`, `DeleteFunc`, `Compact` **return** the new slice; the old one may be changed and its tail zeroed.
- Sort + `Compact` = unique values.
- Iterators (Go 1.23): `for v := range slices.Values(s)`, `slices.All`, `slices.Backward`, `slices.Chunk`; finish with `slices.Collect` or `slices.Sorted`.
- `sort.Ints` / `sort.Slice` are legacy: read them, write `slices`.
