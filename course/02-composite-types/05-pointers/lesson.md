# Pointers

## 1. What a pointer is

Every variable lives somewhere in memory. A **pointer** holds the **address** of a variable. Through it you can read or change that variable without having a copy of it.

```go
x := 42
p := &x         // & = "address of": p points to x
fmt.Println(*p) // 42  — * = "dereference": the variable p points to
*p = 21         // write through the pointer
fmt.Println(x)  // 21  — x itself changed
fmt.Printf("%T %T\n", p, x) // *int int
```

| Syntax | Name | Meaning |
|---|---|---|
| `*T` (in a type) | pointer type | "pointer to a `T`": `*int`, `*Point`, `*[]string` |
| `&v` | address-of | a `*T` pointing to the variable `v` (`v` must be addressable: a variable, field, slice element, or a composite literal) |
| `*p` (in an expression) | dereference / indirection | the variable `p` points to; usable on both sides of `=` |

- `&` works on **variables**, not values: `&x` ✅, `&42` ❌, `&f()` ❌. (`&Point{...}` is a special, allowed case: §5.)
- A `*int` and an `int` are different types: no implicit conversion either way.

## 2. The zero value is `nil`

```go
var p *int          // nil: points to nothing
fmt.Println(p == nil, p) // true <nil>
fmt.Println(*p)     // panic: runtime error: invalid memory address or nil pointer dereference
```

- Dereferencing a nil pointer **panics** (crashes the program). Check `if p != nil` whenever a pointer may legitimately be nil.
- Functions that return a pointer document whether it can be nil; callers check.

```go
// SafeDeref returns *p, or def if p is nil.
func SafeDeref(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}
```

## 3. `new(T)` and `new(expr)` (Go 1.26+)

The built-in `new` allocates a variable and returns a pointer to it. The variable has no name; the pointer is the only way to reach it.

```go
p := new(int)    // *int pointing to 0           (argument is a type)
q := new(42)     // *int pointing to 42          (Go 1.26+: argument is an expression)
s := new("hi")   // *string pointing to "hi"
b := new(x > 10) // *bool pointing to true or false
```

| Form | Allocates | Value |
|---|---|---|
| `new(T)` | a variable of type `T` | the zero value of `T` |
| `new(expr)` (Go 1.26+) | a variable of the type of `expr` | the value of `expr`; an untyped constant uses its default type (`new(42)` is `*int`, `new(1.5)` is `*float64`) |

Before Go 1.26, getting a `*int` holding `36` took two lines (`age := 36; p := &age`) or a helper function. Now it is one expression, which shines in struct literals:

```go
type Person struct {
	Name string
	Age  *int // nil = "unknown"
}

func yearsSince(born int) int { return 2026 - born }

alice := Person{Name: "Alice", Age: new(yearsSince(1990))} // *alice.Age == 36
```

`new(T)` is the same as `&T{}` for structs; most Go code writes `&T{...}` because it can also set fields. `new(nil)` does not compile: `nil` has no type.

## 4. Pointers to structs: automatic dereference

```go
type Point struct{ X, Y int }

pt := Point{1, 2}
p := &pt
(*p).X = 10 // explicit dereference: legal but noisy
p.X = 10    // same thing: Go dereferences p for you
fmt.Println(pt) // {10 2}
fmt.Println(p)  // &{10 2} — fmt prints a struct pointer as &{...}
```

`p.X` on a pointer means `(*p).X`. If `p` is nil, `p.X` panics, exactly like `*p`.

The same shortcut works for **pointers to arrays**: `pa[i]`, `len(pa)` and `for i := range pa` work on a `*[N]T`. It does **not** work for pointers to slices or maps: with `ps *[]int` you must write `(*ps)[i]` and `len(*ps)`.

## 5. `&T{...}`: pointer to a composite literal

Taking the address of a composite literal allocates a new variable and returns its address, in one step:

```go
p := &Point{X: 1, Y: 2}    // *Point
n := &Node{Val: 1}         // the idiomatic way to create linked nodes
```

This is the most common way you will create pointers in Go.

## 6. Passing pointers to mutate

Go passes **everything by value** (FAQ): a function receives a copy of each argument. To let a function change the caller's variable, pass its address; the **pointer** is copied, but the copy points to the same variable.

```go
func scaleValue(p Point, k int) { p.X *= k; p.Y *= k } // changes a copy: useless
func scale(p *Point, k int)     { p.X *= k; p.Y *= k } // changes the caller's Point

pt := Point{1, 2}
scaleValue(pt, 10) // pt is still {1 2}
scale(&pt, 10)     // pt is now {10 20}
```

```go
func Swap(a, b *int) {
	*a, *b = *b, *a
}

x, y := 1, 2
Swap(&x, &y) // x == 2, y == 1
```

(Methods with pointer receivers, `func (p *Point) Scale(k int)`, come next module; they are this same idea with nicer call syntax.)

## 7. Returning a pointer to a local is safe

```go
func newPoint(x, y int) *Point {
	p := Point{X: x, Y: y}
	return &p // perfectly fine in Go
}
```

In C this would be a dangling pointer. In Go, the compiler's **escape analysis** notices that `p` outlives the call and allocates it on the **heap** instead of the stack (`go build -gcflags=-m` reports `moved to heap: p`). The **garbage collector** frees it once nothing points to it. You never choose stack vs heap, and you never free memory yourself; "where does it live?" is an optimization detail, not a correctness one.

## 8. Pointers vs slices and maps

A slice is a small header (pointer to a backing array, length, capacity); a map value is a pointer to the map's internal data. Passing them copies the header/pointer, so **the function already shares the elements**:

| You pass… | Callee can change elements / entries? | Callee can change length, or replace the whole thing? |
|---|---|---|
| `[]T` | ✅ `s[i] = v` is visible to the caller | ❌ `s = append(s, v)` only changes the callee's copy of the header |
| `*[]T` | ✅ `(*s)[i] = v` | ✅ `*s = append(*s, v)` |
| `map[K]V` | ✅ `m[k] = v`, `delete(m, k)` | ❌ `m = make(...)` is local |
| `[N]T` (array) | ❌ the whole array is copied | — |
| `*[N]T` | ✅ `a[i] = v` | — (arrays have a fixed length) |

```go
func appendBroken(s []int, v int) { s = append(s, v) }  // caller sees nothing
func appendTo(s *[]int, v int)    { *s = append(*s, v) } // caller's slice grows

var s []int
appendBroken(s, 1) // s == [] (len 0)
appendTo(&s, 1)    // s == [1]
```

- Idiomatic Go usually **returns** the new slice instead (`s = append(s, v)`, like `append` itself). Use `*[]T` when the function also returns something else (e.g. `Pop` returns the value) or when the slice is a field you mutate in place.
- A `*map[K]V` is almost never needed.

## 9. Optional values

Every Go type has a zero value, so a plain `int` field can't tell "set to 0" from "not set". A pointer can: `nil` means *absent*.

```go
type Config struct {
	Retries *int  // nil = use the default; new(0) = explicitly zero retries
	Verbose *bool // nil = use the default; new(false) = explicitly off
}

func retries(c Config) int {
	if c.Retries == nil {
		return 3 // default
	}
	return *c.Retries
}

retries(Config{})                // 3
retries(Config{Retries: new(0)}) // 0
```

**JSON preview** (`encoding/json` is covered in a later module): pointer fields map naturally to "missing vs present".

```go
type Person struct {
	Name string
	Age  *int `json:"age,omitempty"` // omitted when nil
}

json.Marshal(Person{Name: "Bob", Age: new(0)}) // {"Name":"Bob","age":0}
json.Marshal(Person{Name: "Carol"})            // {"Name":"Carol"}
```

Decoding works the other way: a missing `"age"` leaves `Age` nil; `"age": 0` makes it point to `0`.

## 10. Comparing pointers: identity, not contents

Two pointers are equal if they point to the **same variable** (or are both nil).

```go
p1 := &Point{1, 2}
p2 := &Point{1, 2}
fmt.Println(p1 == p2)   // false — two different variables
fmt.Println(*p1 == *p2) // true  — equal contents
fmt.Println(p1 == p1)   // true
```

Use `p == q` to ask "is this the same node?" (cycle detection, §14); compare `*p == *q` to ask "same value?".

## 11. No pointer arithmetic

`p++`, `p + 1`, `p[i]` on a `*int` don't compile. Go has no pointer arithmetic, which (together with the GC and bounds checks) makes memory corruption bugs impossible in ordinary code. To walk a sequence, use a slice index; to walk linked data, follow `.Next`. (Package `unsafe` exists for low-level code; you won't need it.)

## 12. Pointers into arrays and slices (and the `append` trap)

You can take the address of an element: `&arr[i]`, `&s[i]`, `&s[i].Field`.

```go
s := make([]int, 3) // len 3, cap 3
e := &s[0]          // points into s's backing array
s = append(s, 4)    // cap exceeded → append copies everything to a NEW array
*e = 99             // writes to the OLD array
fmt.Println(s[0], *e) // 0 99 — e no longer refers to s[0]
```

- A pointer to a slice element is only trustworthy while the slice is **not re-grown**. Don't keep `&s[i]` across an `append` to `s`.
- Pointers into a **map** are not allowed at all: `&m[k]` doesn't compile (maps move their entries as they grow).

## 13. `range` copies: modify through the index

The value variable in `for _, v := range s` is a **copy** of each element:

```go
pts := []Point{{1, 1}, {2, 2}}
for _, p := range pts {
	p.X = 100 // modifies the copy
}
fmt.Println(pts) // [{1 1} {2 2}]

for i := range pts {
	pts[i].X = 100 // modifies the element in the slice
}
fmt.Println(pts) // [{100 1} {100 2}]
```

Equivalent: `p := &pts[i]; p.X = 100`. Alternatively store pointers, `[]*Point`, when elements must be shared or are large; then `for _, p := range ps { p.X = 100 }` works because the copy is a pointer.

## 14. Linked data structures

Pointers let a value refer to another value of the **same type**: the basis of lists and trees. A struct can't contain itself (infinite size) but it can contain a pointer to itself.

```go
// Node is one element of a singly linked list. A nil *Node is the empty list.
type Node struct {
	Val  int
	Next *Node
}

// Push returns a new list with v in front of head.
func Push(head *Node, v int) *Node {
	return &Node{Val: v, Next: head}
}

var head *Node // empty list
for _, v := range []int{3, 2, 1} {
	head = Push(head, v)
}
for n := head; n != nil; n = n.Next { // the standard traversal loop
	fmt.Print(n.Val, " ") // 1 2 3
}
```

A binary search tree: smaller keys to the left, larger to the right.

```go
type Tree struct {
	Key         int
	Left, Right *Tree
}

// Insert returns the root of the tree with key added.
func Insert(t *Tree, key int) *Tree {
	if t == nil {
		return &Tree{Key: key} // empty spot found: new leaf
	}
	if key < t.Key {
		t.Left = Insert(t.Left, key)
	} else if key > t.Key {
		t.Right = Insert(t.Right, key)
	}
	return t
}
```

Patterns worth memorizing:
- **nil is the empty list/tree**, so functions take and return `*Node` and handle `nil` first.
- Functions that may change the head/root **return the new head/root**; the caller writes `head = Push(head, v)`.
- Recursion follows the shape of the data: "handle `nil`, then recurse on `Left`/`Right` or `Next`".

## 15. When to use a pointer

| Use a pointer (`*T`) when… | Use a value (`T`) when… |
|---|---|
| The function must **modify** the caller's variable | The function only reads it |
| The struct is **large** (many fields, big arrays) and copied often | It is small (a few words: `Point`, `time.Time`-sized) |
| The value is **optional**: nil means "absent" | Its zero value is a meaningful default |
| Several parts of the program must share **one identity** (a node, a single account, a cache) | Independent copies are what you want (safer: no one else can change it) |
| Building **linked** structures (`Next *Node`) | Plain data stored in slices and maps |

Values are the default in Go. Reach for a pointer for one of the reasons on the left, not "for speed" by reflex: copying a small struct is cheap and keeps it on the stack.

## Common mistakes

| Mistake | Fix |
|---|---|
| Dereferencing a pointer that may be nil (`*p`, `p.X`, `n.Next.Val`) | Check `if p != nil` first; loop with `for n != nil` |
| `var p *Point; p.X = 1` | `p` is nil: allocate first with `p := &Point{}` or `new(Point)` |
| Expecting `func f(p Point)` to change the caller's struct | Take `*Point` and call `f(&pt)` |
| `for _, v := range s { v.X = 1 }` | `for i := range s { s[i].X = 1 }` |
| Keeping `&s[i]` and then `append`ing to `s` | Re-take the address after appending, store indexes, or use `[]*T` |
| `s = append(s, v)` inside a function taking `[]T`, expecting the caller to see it | Return the slice, or take `*[]T` and write `*s = append(*s, v)` |
| `ps[i]` / `len(ps)` with `ps *[]int` | `(*ps)[i]`, `len(*ps)` (auto-dereference only applies to struct fields and arrays) |
| Comparing `p1 == p2` to check equal contents | That compares identity; use `*p1 == *p2` |
| Using `int` for an optional field, then guessing whether `0` meant "unset" | Use `*int` (`nil` = unset) and `new(v)` to set it |
| Pre-Go 1.22 code taking `&v` of the range variable (every pointer was the same variable) | Since Go 1.22 each iteration has its own `v`; still, `&s[i]` is clearer when you want the element itself |

## Key takeaways

- `*T` is "pointer to T"; `&v` takes an address; `*p` reads or writes the variable it points to.
- The zero value is `nil`; dereferencing nil panics. Check before you dereference.
- `new(T)` gives a pointer to a zero `T`; since Go 1.26 `new(expr)` gives a pointer to a copy of `expr`'s value: `Age: new(36)`.
- `p.X` auto-dereferences struct pointers; `&T{...}` is the idiomatic way to allocate.
- Everything is passed by value; pass a pointer to let the callee mutate the caller's variable.
- Returning `&local` is safe: escape analysis moves it to the heap and the GC frees it.
- Slices and maps already share their elements; you need `*[]T` only to change a slice's length for the caller.
- Pointers make optional fields possible (`nil` = absent) and compare by identity.
- No pointer arithmetic. `&s[i]` is invalidated by a growing `append`. `range` copies: modify via `s[i]`.
- Lists and trees are structs with `*Node` fields; `nil` is the empty structure; return the new head/root.
