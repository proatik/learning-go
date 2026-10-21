# Maps

## 1. The map type

A map is an **unordered** collection of key/value pairs with fast lookup, insert and delete by key (a hash table).

```go
map[KeyType]ValueType

map[string]int          // word → count
map[rune][]int          // character → positions
map[string]map[string]int // student → (subject → score)
```

| Part | Requirement |
|---|---|
| **Key type** | Must be **comparable**: `==` and `!=` are defined for it |
| **Value type** | Anything: numbers, strings, slices, maps, functions… |

| Allowed as keys ✅ | Not allowed ❌ (compile error `invalid map key type`) |
|---|---|
| `bool`, all numeric types, `string`, `rune`/`byte` | slices `[]T` |
| arrays of comparable types, e.g. `[2]int` | maps |
| (later: pointers, structs of comparable fields, channels) | functions |

```go
grid := map[[2]int]string{{0, 0}: "origin"} // an array key: fine
fmt.Println(grid[[2]int{0, 0}])             // origin

var bad map[[]int]string // compile error: invalid map key type []int
```

Need a slice as a key? Convert it to something comparable first, e.g. a `string` (`string(runes)`) or a fixed-size array.

## 2. Creating maps

```go
// 1. Map literal: keys and values listed up front
ages := map[string]int{
	"alice": 31,
	"bob":   28, // trailing comma required when the } is on its own line
}

// 2. Empty literal
counts := map[string]int{}

// 3. make, optionally with a size hint
index := make(map[string]int)
big := make(map[string]int, 1000) // room for ~1000 entries before it must grow; len(big) is still 0
```

- In a literal, inner composite types may be elided: `map[string][]int{"a": {1, 2}}` means `"a": []int{1, 2}`.
- A duplicate constant key in a literal is a compile error.
- The size hint is only a performance tip; the map always grows as needed. Use it when you know the final size (e.g. building a map from a slice of known length).

## 3. Read, write, delete, len

```go
ages := map[string]int{"alice": 31, "bob": 28}

ages["carol"] = 40     // insert
ages["bob"]++          // update in place: 29 (read, add, write back)
fmt.Println(ages["bob"], len(ages)) // 29 3

delete(ages, "alice")  // remove
delete(ages, "zzz")    // deleting a missing key is a no-op, no error
fmt.Println(ages)      // map[bob:29 carol:40]
```

| Operation | Syntax | Notes |
|---|---|---|
| Read | `v := m[k]` | Missing key → **zero value** of the value type |
| Write | `m[k] = v` | Inserts or overwrites |
| Delete | `delete(m, k)` | No-op if `k` is absent |
| Size | `len(m)` | Number of entries (`cap` doesn't work on maps) |
| Remove all | `clear(m)` | Built-in since Go 1.21; `len(m)` becomes 0 |

## 4. Missing keys and the comma-ok idiom

A lookup **never fails**: a missing key yields the zero value. That is often exactly what you want:

```go
counts := map[string]int{}
counts["go"]++ // counts["go"] was 0 → now 1. No "if exists" check needed.
```

When you must tell "absent" from "present with zero value", use the **two-value** form:

```go
v, ok := ages["dave"] // ok is true only if the key is present
fmt.Println(v, ok)    // 0 false

if age, ok := ages["bob"]; ok {
	fmt.Println("bob is", age)
}

_, exists := ages["carol"] // only test presence
```

The second result is conventionally named `ok`, hence "comma ok".

## 5. Iterating with `range`

```go
for k, v := range m { ... } // key and value
for k := range m { ... }    // keys only
for _, v := range m { ... } // values only
```

**Iteration order is unspecified and deliberately randomized**: two loops over the same map can visit keys in different orders. Never rely on it, not even in tests.

For deterministic output, **sort the keys** first:

```go
import (
	"maps"
	"slices"
)

m := map[string]int{"b": 2, "a": 1, "c": 3}
for _, k := range slices.Sorted(maps.Keys(m)) { // maps.Keys returns an iterator; slices.Sorted collects + sorts it
	fmt.Println(k, m[k])
}
// a 1
// b 2
// c 3
```

- `fmt.Println(m)` / `%v` is the one exception: the `fmt` package **prints maps sorted by key** (`map[a:1 b:2 c:3]`), so printing a whole map is stable. `range` is not.
- **Deleting** entries during `range` is safe: a deleted entry not yet reached will not be produced.
- **Adding** entries during `range` is allowed, but a new entry may or may not be visited. Avoid it.

```go
for k, v := range counts { // remove every rare word
	if v < 2 {
		delete(counts, k) // safe
	}
}
```

## 6. The nil map

The zero value of a map type is `nil`. A nil map behaves like an **empty, read-only** map:

```go
var m map[string]int // nil

m["x"]          // 0  — reading is fine
len(m)          // 0
for range m {}  // zero iterations
delete(m, "x")  // no-op
m == nil        // true (maps can only be compared to nil)

m["x"] = 1 // panic: assignment to entry in nil map
```

Rule: **initialize before writing**, with a literal or `make`. A function that returns a map usually returns an empty, initialized map (`map[K]V{}`) rather than `nil`, so callers can write to it.

## 7. Maps are reference-like

A map value is a small handle pointing to the underlying table. Assigning a map or passing it to a function **copies the handle, not the data**: both refer to the same entries.

```go
func addOne(m map[string]int) { m["added"] = 1 }

x := map[string]int{}
addOne(x)    // the function sees and changes x's data
y := x       // y and x are the same map
y["y"] = 2
fmt.Println(x) // map[added:1 y:2]
```

- No pointer (`*map[...]`) is needed to modify a map in a function.
- To get an **independent copy**, use `maps.Clone(m)` (§11). This matters when a function must **not mutate** its input.
- `clear(m)` empties the shared table, so every holder of the map sees it empty. It does **not** make `m` nil.

## 8. Map of slices

Because a missing key yields a nil slice, and `append` works on nil slices, grouping needs no initialization:

```go
groups := map[string][]string{}
groups["fruit"] = append(groups["fruit"], "apple") // groups["fruit"] was nil
groups["fruit"] = append(groups["fruit"], "pear")
fmt.Println(groups, len(groups["veg"])) // map[fruit:[apple pear]] 0
```

Always assign the result of `append` back to `m[k]`: `append(m[k], x)` alone changes nothing in the map.

## 9. Nested maps

The inner map is **not** created automatically. Its zero value is `nil`, and writing to it panics:

```go
grades := map[string]map[string]int{}
grades["alice"]["math"] = 90 // panic: grades["alice"] is a nil map

// Fix: create the inner map on first use
inner, ok := grades["alice"]
if !ok {
	inner = map[string]int{}
	grades["alice"] = inner // store the handle; inner and grades["alice"] are the same map
}
inner["math"] = 90
fmt.Println(grades) // map[alice:map[math:90]]
```

Reading through a missing outer key is safe: `grades["bob"]["math"]` is `0` (indexing a nil inner map).

## 10. Idioms: sets and counting

Go has no set type; a map with the elements as keys is the set.

| Set type | Add | Test | Trade-off |
|---|---|---|---|
| `map[string]bool` | `s[x] = true` | `if s[x] { … }` | Reads nicely; absent → `false` for free |
| `map[string]struct{}` | `s[x] = struct{}{}` | `if _, ok := s[x]; ok { … }` | Values take **zero bytes**; makes "set" intent explicit |

```go
seen := map[string]bool{}
seen["go"] = true
fmt.Println(seen["go"], seen["rust"]) // true false

set := map[string]struct{}{} // struct{} is the empty type; struct{}{} its only value
set["go"] = struct{}{}
_, has := set["go"] // true
```

Counting relies on the zero value:

```go
counts := map[string]int{}
for _, w := range strings.Fields("the cat the hat") {
	counts[w]++
}
fmt.Println(counts) // map[cat:1 hat:1 the:2]
```

## 11. The `maps` package

Generic helpers that work on any map type (`import "maps"`). The iterator-returning ones (`Keys`, `Values`, `All`) pair with `range` and with `slices.Sorted` / `slices.Collect` from the previous lesson.

| Function | Does |
|---|---|
| `maps.Keys(m)` | Iterator (`iter.Seq[K]`) over the keys, in unspecified order |
| `maps.Values(m)` | Iterator (`iter.Seq[V]`) over the values |
| `maps.All(m)` | Iterator (`iter.Seq2[K, V]`) over key/value pairs |
| `maps.Collect(seq)` | Builds a new map from an `iter.Seq2[K, V]` |
| `maps.Insert(m, seq)` | Adds the pairs from `seq` into `m`, overwriting |
| `maps.Clone(m)` | Shallow copy (a nil map clones to nil) |
| `maps.Copy(dst, src)` | Copies every pair of `src` into `dst`, overwriting existing keys |
| `maps.Equal(m1, m2)` | Same keys with equal values (values must be comparable) |
| `maps.EqualFunc(m1, m2, eq)` | Like `Equal`, comparing values with `eq` |
| `maps.DeleteFunc(m, del)` | Deletes every pair for which `del(k, v)` returns `true` |

```go
src := map[string]int{"a": 1, "b": 2}

c := maps.Clone(src)
c["a"] = 100
fmt.Println(src, c, maps.Equal(src, c)) // map[a:1 b:2] map[a:100 b:2] false

dst := map[string]int{"a": 0, "z": 26}
maps.Copy(dst, src)
fmt.Println(dst) // map[a:1 b:2 z:26]

maps.DeleteFunc(dst, func(k string, v int) bool { return v > 10 })
fmt.Println(dst) // map[a:1 b:2]

fmt.Println(slices.Sorted(maps.Values(src))) // [1 2]

for k, v := range maps.All(src) { ... } // same as ranging over src directly
```

- `Clone` is **shallow**: for `map[string][]int` the clone's slices share backing arrays with the original; for nested maps, the inner maps are shared.
- `==` between two maps does not compile; use `maps.Equal` (or `reflect.DeepEqual` for nested values, mostly in tests).

## 12. Concurrency (preview)

Maps are **not safe for concurrent use**: if one goroutine writes while another reads or writes, the program can crash (`fatal error: concurrent map writes`). Synchronization (`sync.Mutex`, `sync.Map`) comes in the concurrency module.

## Common mistakes

| Mistake | Fix |
|---|---|
| `var m map[string]int; m["x"] = 1` → panic | Initialize first: `m := map[string]int{}` or `make(map[string]int)` |
| Relying on `range` order (printing, building a result slice, tests) | Sort the keys: `slices.Sorted(maps.Keys(m))` |
| `grades["alice"]["math"] = 90` with no inner map → panic | Check with comma-ok and create the inner map first (§9) |
| `p := &m["a"]` | Compile error `cannot take address of m["a"]`: map entries move as the map grows. Read, modify, write back: `v := m["a"]; v += 1; m["a"] = v` |
| `if m[k] == 0` to test presence | `0` may be a stored value: use `_, ok := m[k]` |
| `append(m[k], x)` without assigning | `m[k] = append(m[k], x)` |
| Assuming a function got a copy of the map | Maps share data; `maps.Clone` if the input must stay untouched |
| `m1 == m2` | Only `m == nil` compiles; use `maps.Equal(m1, m2)` |
| Using a slice as a key | Use a `string` or an array (`[N]T`) key instead |
| Thinking `maps.Clone` deep-copies nested slices/maps | It is shallow: clone the inner values yourself if you will mutate them |

## Key takeaways

- `map[K]V`: K must be comparable (no slices, maps, funcs); V can be anything.
- Create with a literal or `make(map[K]V, hint)`; the nil map is readable but **panics on write**.
- Missing key → zero value. Use `v, ok := m[k]` to distinguish absent from zero.
- `m[k]++` counts, `m[k] = append(m[k], x)` groups, `map[T]struct{}` / `map[T]bool` are sets.
- `range` order is random: sort keys with `slices.Sorted(maps.Keys(m))` for deterministic output. Deleting while ranging is safe.
- Maps are reference-like: passing or assigning shares the data; `maps.Clone` copies (shallowly).
- Nested maps need their inner map created before writing.
- `delete`, `clear`, `len` are built-ins; `maps.Keys/Values/All/Collect/Insert/Clone/Copy/Equal/DeleteFunc` cover the rest.
- Map elements are not addressable, and maps are not safe for concurrent writes.
