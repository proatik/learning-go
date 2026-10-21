# Methods

## 1. Declaring a method

A **method** is a function with a **receiver**: an extra parameter written before the name. It attaches the function to a type.

```go
type Rect struct {
	W, H float64
}

// Area is a method on Rect; r is the receiver.
func (r Rect) Area() float64 {
	return r.W * r.H
}

r := Rect{W: 3, H: 4}
fmt.Println(r.Area()) // 12 — call it with the selector syntax value.Method()
```

| Part | Example | Notes |
|---|---|---|
| Receiver | `(r Rect)` | one parameter only; name it with a short abbreviation of the type (`r`, `c`, `acc`), **never** `this` or `self` |
| Name | `Area` | capitalized = exported, like functions and fields |
| Params / results | `() float64` | exactly like a function |

- A method is just a function with nicer call syntax: `r.Area()` behaves like `Area(r)`.
- Different types may have methods with the same name (`Rect.Area`, `Circle.Area`). Within one type, names must be unique (no overloading), and a method can't share a name with a field.
- Methods are what **interfaces** are built on (next lesson): an interface is a set of method signatures.

## 2. Methods on any named type of your package

The receiver type can be **any named type declared in the same package**, not just structs.

```go
type Celsius float64
type IntList []int

func (c Celsius) Fahrenheit() float64 { return float64(c)*9/5 + 32 }

func (l IntList) Sum() int {
	total := 0
	for _, v := range l {
		total += v
	}
	return total
}

fmt.Println(Celsius(100).Fahrenheit())  // 212
fmt.Println(IntList{1, 2, 3}.Sum())    // 6
```

| You try | Result |
|---|---|
| `func (i int) Double() int` | ❌ `cannot define new methods on non-local type int` |
| `func (t time.Duration) Pretty() string` | ❌ same: the type belongs to package `time` |
| `type MyInt int` then `func (m MyInt) Double() MyInt` | ✅ define your own type first |
| receiver of pointer type `type P *int` or an interface type | ❌ the base type can't itself be a pointer or interface |

A named type gets **none** of its underlying type's methods, but keeps its operators: `Celsius` supports `+`, `<`, conversion `float64(c)`, etc.

## 3. Value receivers vs pointer receivers

A receiver is a parameter, so a **value receiver gets a copy**. To change the caller's value, use a **pointer receiver** `*T` (same idea as passing `&x` in the pointers lesson).

```go
func (r Rect) ScaleBroken(k float64) { r.W *= k; r.H *= k } // changes the copy only
func (r *Rect) Scale(k float64)      { r.W *= k; r.H *= k } // changes the caller's Rect

r := Rect{W: 2, H: 3}
r.ScaleBroken(10)
fmt.Println(r) // {2 3}  — nothing happened
r.Scale(10)
fmt.Println(r) // {20 30}
```

When to use which:

| Use a **pointer receiver** `(t *T)` when… | A **value receiver** `(t T)` is fine when… |
|---|---|
| the method **modifies** the receiver | the type is small and immutable-feeling (`Celsius`, `Point`, `Money`) |
| the struct is **large** (copying every call is wasteful) | it's a map, func or chan type (already reference-like) |
| the struct contains something that must not be copied (e.g. a `sync.Mutex`, later module) | it's a slice type whose methods don't reslice or append |
| **any other method** of the type already uses a pointer receiver | |

**Consistency rule** (Go Code Review Comments): if **any** method needs a pointer receiver, give **all** methods of that type pointer receivers. Mixed receiver kinds confuse readers and cause trouble with interfaces (next lesson).

Slice types are a classic trap: `append` may move the backing array, so a method that grows the slice needs a pointer.

```go
type Stack []int

func (s *Stack) Push(v int) { *s = append(*s, v) } // must update the caller's header
func (s *Stack) Len() int   { return len(*s) }     // pointer too, for consistency
```

## 4. Automatic `&` and `*` when calling

Go inserts the `&` or `*` for you, so you rarely think about it at the call site.

```go
type Counter struct{ n int }

func (c *Counter) Inc()      { c.n++ }
func (c Counter) Value() int { return c.n }

var c Counter
c.Inc()            // c is addressable → Go rewrites it as (&c).Inc()
p := &c
fmt.Println(p.Value()) // 1 — Go rewrites it as (*p).Value()
```

| Receiver declared | Call on a value variable `v` | Call on a pointer `p` |
|---|---|---|
| `(t T)` | ✅ `v.M()` | ✅ `p.M()` → `(*p).M()` (panics if `p` is nil) |
| `(t *T)` | ✅ `v.M()` → `(&v).M()` **only if `v` is addressable** | ✅ `p.M()` |

**Addressable** = a variable, a field of an addressable struct, a slice element, or `*p`. Not addressable: **map elements**, **literals** and **function results**.

```go
m := map[string]Counter{}
m["a"].Inc()     // ❌ cannot call pointer method Inc on Counter
Counter{}.Inc()  // ❌ cannot call pointer method Inc on Counter
newCounter().Inc() // ❌ same, if newCounter returns a Counter

pm := map[string]*Counter{"a": {}}
pm["a"].Inc()    // ✅ the element is already a pointer

cs := []Counter{{}, {}}
cs[0].Inc()      // ✅ slice elements are addressable
(&Counter{}).Inc() // ✅ &T{} is allowed
```

## 5. Nil receivers can be valid

A pointer-receiver method can be called on a **nil** pointer; it only panics if it dereferences it. Handling nil inside the method makes the zero value (a nil pointer) usable.

```go
type List struct {
	Val  int
	Next *List
}

// Len reports the number of nodes; a nil *List is the empty list.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return 1 + l.Next.Len() // l.Next may be nil: still fine
}

var l *List
fmt.Println(l.Len()) // 0 — no panic
l = &List{Val: 1, Next: &List{Val: 2}}
fmt.Println(l.Len()) // 2
```

This is the Module 2 linked list, now as methods: `l.Len()` instead of `Len(l)`. Document nil-safety in the doc comment so callers can rely on it.

## 6. Method sets (one-line preview)

The **method set** of `T` contains its value-receiver methods; the method set of `*T` contains **both** value and pointer-receiver methods. This decides which types satisfy an interface: details in the next lesson.

## 7. Methods promoted through embedding

Embedding a type promotes its **methods** as well as its fields (the structs lesson showed fields).

```go
type Animal struct{ name string }

func (a Animal) Name() string     { return a.name }
func (a Animal) Describe() string { return a.name + " is an animal" }
func (a *Animal) Rename(n string) { a.name = n }

type Dog struct {
	Animal // embedded
	Breed  string
}

// Describe shadows (overrides) the promoted Animal.Describe.
func (d Dog) Describe() string { return d.Name() + " is a dog" }

d := Dog{Animal: Animal{name: "Rex"}, Breed: "beagle"}
fmt.Println(d.Name())            // Rex — promoted from Animal
fmt.Println(d.Describe())        // Rex is a dog — Dog's own method wins (shallower)
fmt.Println(d.Animal.Describe()) // Rex is an animal — the embedded one is still reachable
d.Rename("Max")                  // promoted pointer method: (&d.Animal).Rename("Max")
fmt.Println(d.Name())            // Max
```

- The **shallowest** method wins, exactly like fields; two promoted methods with the same name at the same depth → `ambiguous selector` (only when you call it).
- ⚠️ No virtual dispatch: inside `Animal.Describe`, the receiver is an `Animal`. If an `Animal` method calls `a.Describe()`, it always runs `Animal.Describe`, even when reached through a `Dog`. Embedding is composition, not inheritance.
- To extend rather than replace, call the embedded method explicitly: `return d.Animal.Describe() + " (" + d.Breed + ")"`.

## 8. Method values and method expressions

| Form | Example | Type | Meaning |
|---|---|---|---|
| **Method value** | `f := c.Inc` | `func()` | a function **bound** to the receiver `c` |
| **Method expression** | `g := (*Counter).Inc` | `func(*Counter)` | an unbound function; the receiver becomes the **first parameter** |
| | `h := Counter.Value` | `func(Counter) int` | |

```go
var c Counter
inc := c.Inc // receiver evaluated now: &c is stored inside inc
inc()
inc()
fmt.Println(c.Value()) // 2

val := c.Value // value receiver: a COPY of c is stored now
c.Inc()
fmt.Println(c.Value(), val()) // 3 2 — val still sees the old copy

g := (*Counter).Inc
g(&c)          // same as c.Inc()
h := Counter.Value
fmt.Println(h(c)) // 4
```

Method values are handy wherever a `func` is expected, e.g. a slice of operations:

```go
var c Counter
ops := []func(){c.Inc, c.Inc, c.Inc}
for _, op := range ops {
	op()
}
fmt.Println(c.Value()) // 3
```

## 9. Constructors and encapsulation

Unexported fields plus exported methods give you a type whose invariants **other packages cannot break**.

```go
package bank

// Account is a bank account whose balance never goes negative.
type Account struct {
	owner   string // unexported: only package bank can touch these
	balance int64
}

// NewAccount returns an account for owner with a zero balance.
func NewAccount(owner string) *Account {
	return &Account{owner: owner}
}

// Balance returns the current balance in cents.
func (a *Account) Balance() int64 { return a.balance }

// Withdraw takes amount out and reports whether it succeeded.
func (a *Account) Withdraw(amount int64) bool {
	if amount <= 0 || amount > a.balance {
		return false
	}
	a.balance -= amount
	return true
}
```

- Name constructors `NewT` (or just `New` when the package has one main type: `list.New`). Return `*T` when methods use pointer receivers.
- Getters are named after the field, **without `Get`**: `a.Balance()`, not `a.GetBalance()`. Setters use `Set`: `a.SetOwner(n)`.
- Add getters/setters **only when useful**. A plain data struct (`Point`, `Config`) should just export its fields.
- Prefer a useful zero value (§5) when you can; a constructor is for validation or non-zero defaults.

## 10. `String() string`: controlling how fmt prints

If a type has a method `String() string`, `fmt` uses it for `%v`, `%s` and `Println`.

```go
type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

t := Celsius(21.5)
fmt.Println(t)                // 21.5°C
fmt.Printf("%v | %s\n", t, t) // 21.5°C | 21.5°C
fmt.Printf("%.2f\n", t)       // 21.50 — numeric verbs ignore String
```

- ⚠️ Inside `String`, never format the receiver itself with `%v`/`%s` (`fmt.Sprintf("%v", c)`): that calls `String` again → infinite recursion. Convert first: `float64(c)`.
- If `String` has a **pointer** receiver, only `*T` values get it: `fmt.Println(v)` on a plain `T` value ignores it. Use a value receiver for `String` on small types.
- Why this works (it's the `fmt.Stringer` interface) is in lesson 3.

**Go 1.27:** methods may declare their own type parameters (generic methods, e.g. `func (b Box) Map[T any](f func(int) T) T`); covered in the generics module.

## Common mistakes

| Mistake | Fix |
|---|---|
| Value receiver on a method that modifies: the change "doesn't save" | Use a pointer receiver `(t *T)` |
| Mixing `(t T)` and `(t *T)` methods on one type | If one method needs a pointer, make them all pointers |
| `m[k].Inc()` with a pointer method on a map element | `c := m[k]; c.Inc(); m[k] = c`, or store pointers `map[string]*Counter` |
| `T{}.Inc()` on a literal / function result | Assign to a variable first, or use `(&T{}).Inc()` |
| `func (i int) Double()` | Define a local type: `type MyInt int` |
| `Push` on `type Stack []int` with a value receiver | `func (s *Stack) Push(v int) { *s = append(*s, v) }` |
| Receiver named `this` / `self` | A short name from the type: `func (s *Stack)` |
| Pointer method that panics when the receiver is nil | Check `if l == nil` first when nil means "empty" |
| `String()` calling `fmt.Sprintf("%v", t)` on itself | Format the underlying value: `float64(t)`, or the fields |
| Expecting an embedded type's method to call the outer type's override | No virtual dispatch: pass what's needed explicitly |
| `GetBalance()` | `Balance()` |

## Key takeaways

- `func (r T) Name(...)`: a method is a function with a receiver; declare it on any named type **of your package**.
- Value receivers get a **copy**; pointer receivers can **mutate** and avoid copying large structs. Be consistent per type.
- Go adds `&`/`*` automatically, but pointer methods need an **addressable** value: not map elements, literals or function results.
- A pointer method can handle a **nil receiver**, making nil a valid "empty" value.
- Embedding **promotes methods**; the outer type can shadow them and still call `x.Embedded.M()`. No inheritance, no virtual dispatch.
- `x.M` is a bound **method value**; `T.M` / `(*T).M` is a **method expression** taking the receiver as its first argument.
- Unexported fields + `NewT` constructor + methods = encapsulation. Getters have no `Get` prefix.
- `String() string` controls how `fmt` prints your type.
