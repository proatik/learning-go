# Interfaces

## 1. An interface type is a method set

An **interface type** lists method signatures. A value of an interface type can hold **any** value whose type has all of those methods.

```go
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rect struct{ W, H float64 }

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

type Circle struct{ R float64 }

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

var s Shape     // a variable of interface type
s = Rect{3, 4}  // ✅ Rect has Area and Perimeter
fmt.Println(s.Area()) // 12
s = Circle{1}   // ✅ the same variable can hold a different type
fmt.Printf("%.2f\n", s.Area()) // 3.14
```

- An interface describes **behaviour** ("what it can do"), not data. It has no fields.
- Calling `s.Area()` runs the method of whatever concrete type `s` currently holds: **dynamic dispatch**.
- Only the methods listed in the interface are callable through it: `s.W` doesn't compile even when `s` holds a `Rect`.

## 2. Implicit satisfaction: no `implements`

A type **satisfies** (implements) an interface simply by having the methods. There is no `implements` keyword and no declaration linking the two.

```go
type Square struct{ S float64 }

func (q Square) Area() float64      { return q.S * q.S }
func (q Square) Perimeter() float64 { return 4 * q.S }

var s Shape = Square{2} // works: Square was written without ever mentioning Shape
```

| Consequence | Why it matters |
|---|---|
| Types and interfaces are **decoupled** | A package can define an interface that types from *other* packages (even the standard library) already satisfy |
| Satisfaction is checked **where a value is assigned** | `var s Shape = Square{2}`, passing to a `Shape` parameter, returning as `Shape`… |
| Names and signatures must match **exactly** | `Area() float32` or `area() float64` is a different method: no match |

Compiler message when something is missing:

```text
cannot use Point{} (value of struct type Point) as Shape value in variable declaration:
Point does not implement Shape (missing method Area)
```

## 3. Compile-time assertion

Because satisfaction is implicit, a typo in a method name is only caught where the type is *used* as the interface. Make the promise explicit with a blank variable:

```go
var _ Shape = Rect{}          // value type
var _ Shape = (*Circle)(nil)  // pointer type, without allocating anything
```

- `_` discards the value; the line exists only to make the compiler check `Rect` / `*Circle` against `Shape`.
- `(*Circle)(nil)` is a **conversion** of `nil` to the type `*Circle`: a typed nil pointer, which costs nothing.
- Put it next to the type declaration. The build breaks the moment a method is renamed or its signature changes.

## 4. Method sets: `T` vs `*T`

Recall from the Methods lesson:

| Type | Its method set contains |
|---|---|
| `T` | methods with **value** receivers `(t T)` |
| `*T` | methods with value receivers **and** pointer receivers `(t *T)` |

A type satisfies an interface only if the interface's methods are all in **its method set**. So when some methods use pointer receivers, only `*T` satisfies the interface:

```go
type Counter interface {
	Inc()
	Count() int
}

type Tally struct{ n int }

func (t *Tally) Inc()       { t.n++ } // pointer receiver: must mutate
func (t *Tally) Count() int { return t.n }

var c Counter = &Tally{} // ✅ *Tally has Inc and Count
var d Counter = Tally{}  // ❌ Tally does not implement Counter (method Inc has pointer receiver)
```

Why the rule? An interface stores a **copy** of the value you put in it. Calling a pointer method on that hidden copy would mutate the copy, not your variable — a silent bug — so Go forbids it. (`t.Inc()` on an addressable *variable* `t` still works: that's just the `(&t).Inc()` shorthand.)

| Receivers on `T`'s methods | `T` satisfies? | `*T` satisfies? |
|---|---|---|
| all value | ✅ | ✅ |
| some or all pointer | ❌ | ✅ |

Rule of thumb: if any method needs a pointer receiver, store **pointers** in the interface: `[]Counter{&Tally{}, &Tally{}}`.

## 5. Interface values: a (type, value) pair

At run time an interface value is a pair: the **dynamic type** and the **dynamic value** it holds.

```go
var s Shape                  // (nil, nil)
fmt.Printf("%v %T\n", s, s)  // <nil> <nil>

s = Rect{3, 4}               // (Rect, Rect{3 4})
fmt.Printf("%v %T\n", s, s)  // {3 4} main.Rect

s = &Circle{1}               // (*Circle, pointer to Circle{1})
fmt.Printf("%T\n", s)        // *main.Circle
```

- `%T` prints the dynamic type; `%v` prints the dynamic value.
- Assigning a value to an interface **copies** it (a struct is copied; a pointer copies only the address, so the pointee is shared).
- The zero value of any interface type is `nil`: **both** halves are nil. Calling a method on it panics: `invalid memory address or nil pointer dereference`.

## 6. The nil-interface trap

An interface is `== nil` **only when both the type and the value are nil**. An interface that holds a nil *pointer* has a type, so it is **not** nil.

```go
type Problem interface{ Problem() string }

type FieldProblem struct{ Field string }

func (p *FieldProblem) Problem() string { return p.Field + " is invalid" }

func check(name string) Problem {
	var p *FieldProblem // nil *FieldProblem
	if name == "" {
		p = &FieldProblem{Field: "name"}
	}
	return p // ⚠️ always returns (type *FieldProblem, value nil-or-pointer)
}

fmt.Println(check("ann") == nil) // false!  the interface holds (*FieldProblem, nil)
```

| Interface value | `== nil` |
|---|---|
| `var p Problem` → (nil, nil) | ✅ true |
| `Problem((*FieldProblem)(nil))` → (*FieldProblem, nil) | ❌ false |

Fix: **return the literal `nil`** on the success path; never return a typed nil pointer through an interface result.

```go
func check(name string) Problem {
	if name == "" {
		return &FieldProblem{Field: "name"}
	}
	return nil // (nil, nil): compares equal to nil
}
```

You will meet this trap again with `error` (itself an interface) in the Errors lesson.

## 7. `any`: the empty interface

`interface{}` has no methods, so **every** type satisfies it. `any` is a built-in alias for it.

```go
var x any
x = 42
x = "hello"
x = []int{1, 2}
fmt.Println(x) // [1 2]   (fmt.Println's parameters are ...any)

vals := []any{1, "two", 3.0, nil}
```

- `any` says nothing about the value, so you can't do anything with it directly: `x + 1` doesn't compile. You must first get the concrete value back (§8, §9).
- Use `any` sparingly: heterogeneous data, printing, decoding. A specific interface (or a concrete type) gives you compile-time checking.

## 8. Type assertions

`x.(T)` extracts the dynamic value from an interface `x`, asserting its dynamic type is `T`.

```go
var x any = 42

n, ok := x.(int)    // 42 true
s, ok := x.(string) // "" false   (zero value of string)
n2 := x.(int)       // 42
s2 := x.(string)    // panic: interface conversion: interface {} is int, not string
```

| Form | On mismatch |
|---|---|
| `v, ok := x.(T)` ("comma-ok") | `v` = zero value of `T`, `ok` = false; never panics |
| `v := x.(T)` | **panics** |

`T` may also be an **interface** type: then the assertion asks "does the dynamic value *also* have these methods?" and returns the value as that interface.

```go
type Namer interface{ Name() string }

var s Shape = Rect{1, 2}
if n, ok := s.(Namer); ok { // Rect has no Name method → ok == false
	fmt.Println(n.Name())
}
```

- Asserting on a nil interface fails (`ok == false`, or a panic without `ok`).
- Use the comma-ok form unless a mismatch is truly a programming error.

## 9. Type switches

A **type switch** is a chain of type assertions in one statement.

```go
func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nothing"
	case int, int64: // several types: x keeps the type of v (any)
		return fmt.Sprintf("integer %v", x)
	case string:
		return fmt.Sprintf("string of length %d", len(x)) // x is a string here
	case []int:
		return fmt.Sprintf("%d ints", len(x))
	case fmt.Stringer: // interface case: any type with String() string
		return "stringer " + x.String()
	default:
		return fmt.Sprintf("unknown %T", x) // x has the type of v
	}
}
```

| Case form | Type of `x` inside the case |
|---|---|
| `case T:` (one type) | `T` |
| `case T1, T2:` (several types) | the type of the switched expression (here `any`) |
| `case nil:` | matches a **nil interface** only; `x` is the original nil interface |
| `default:` | the type of the switched expression |

- Cases are tried **top to bottom**; the first match wins. Put specific types before interface cases that would also match.
- `x := v.(type)` is only legal inside `switch`. Without a name (`switch v.(type)`) you only branch on the type.
- `fmt.Stringer` (types with a `String() string` method) is covered properly in the next lesson.

## 10. Embedding interfaces

An interface can **embed** other interfaces: its method set is the union.

```go
type Getter interface {
	Get(key string) (string, bool)
}

type Setter interface {
	Set(key, value string)
}

type Store interface { // like io.ReadWriter = io.Reader + io.Writer
	Getter
	Setter
}
```

- Any `Store` is automatically a `Getter` and a `Setter`: `var g Getter = myStore` works without an assertion.
- Going the other way (`Getter` → `Store`) requires an assertion: `st, ok := g.(Store)`.
- The same method appearing in two embedded interfaces is fine if the signatures are identical.

## 11. Small interfaces, `-er` names

Go favours tiny interfaces, often with **one** method, named after it with an **-er** suffix.

| Interface | Method |
|---|---|
| `fmt.Stringer` | `String() string` |
| `io.Reader` / `io.Writer` | `Read(p []byte)` / `Write(p []byte)` (next lesson) |
| `sort.Interface` | `Len`, `Less`, `Swap` (an older, larger one) |

- "The bigger the interface, the weaker the abstraction" (Go proverb): a 1-method interface is easy to satisfy, easy to fake in tests and easy to compose by embedding.
- Name interfaces after behaviour (`Notifier`, `Handler`), not after the types that implement them.

A one-method interface can even be satisfied by a **function type**: give the func type a method that calls itself (the `http.HandlerFunc` pattern).

```go
type Handler interface{ Handle(args []string) string }

type HandlerFunc func(args []string) string

func (f HandlerFunc) Handle(args []string) string { return f(args) }

var h Handler = HandlerFunc(func(args []string) string { // a conversion, not a call
	return strings.Join(args, "+")
})
fmt.Println(h.Handle([]string{"a", "b"})) // a+b
```

## 12. Accept interfaces, return structs

```go
// ✅ the parameter asks only for what the function needs
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// ✅ constructors return the concrete type: callers get every method and field
func NewMemStore() *MemStore { return &MemStore{data: map[string]string{}} }
```

| Guideline | Reason |
|---|---|
| **Parameters**: the smallest interface that works | Callers can pass any type, including test fakes |
| **Results**: concrete types | Callers aren't limited to the interface; adding methods later doesn't break anyone |
| **Define interfaces where they're consumed** | The package that *uses* the behaviour declares the interface it needs; the implementing package just provides methods |

```go
package report // the consumer

type AreaFinder interface{ Area() float64 } // declared here, where it's needed

func Big(xs []AreaFinder) int { ... }
```

Don't create an interface "just in case" with a single implementation; introduce it when a second implementation (or a test fake) appears.

## 13. Polymorphism with slices of interfaces

A `[]Shape` can hold different concrete types; one loop handles them all.

```go
shapes := []Shape{Rect{3, 4}, Circle{1}, Square{2}}
for _, s := range shapes {
	fmt.Printf("%T %.2f\n", s, s.Area())
}
// main.Rect 12.00
// main.Circle 3.14
// main.Square 4.00
```

⚠️ A `[]Rect` is **not** a `[]Shape`: the element layouts differ, so there's no conversion. Copy element by element:

```go
rects := []Rect{{1, 1}, {2, 2}}
// TotalArea(rects) → cannot use rects (variable of type []Rect) as []Shape value
shapes := make([]Shape, len(rects))
for i, r := range rects {
	shapes[i] = r
}
```

## 14. Sorting interface values

`slices.SortFunc` works on a `[]Shape` like any slice: compare through the interface methods.

```go
slices.SortFunc(shapes, func(a, b Shape) int {
	return cmp.Compare(a.Area(), b.Area())
})
fmt.Println(shapes) // [{1} {2} {3 4}]  — Circle(3.14), Square(4), Rect(12)
```

- Add a tie-breaker with `cmp.Or(cmp.Compare(a.Area(), b.Area()), cmp.Compare(a.Name(), b.Name()))` when the interface has a `Name` method.
- Sort `slices.Clone(shapes)` if the caller's order must be kept.

## 15. Comparing interface values

Two interface values are `==` when their **dynamic types are identical** and their **dynamic values are equal**.

```go
var a, b any = 1, 1
fmt.Println(a == b)            // true
fmt.Println(a == any(int64(1))) // false: int vs int64
var s Shape = Rect{1, 2}
fmt.Println(s == Rect{1, 2})   // true: Rect{1, 2} is converted to Shape first
```

⚠️ If the dynamic type is **not comparable** (slice, map, func), `==` compiles but **panics** at run time:

```go
var x, y any = []int{1}, []int{1}
fmt.Println(x == y) // panic: runtime error: comparing uncomparable type []int
```

The same applies to interface values used as **map keys** (`map[any]int`). Only compare interfaces whose dynamic types you know are comparable.

## Common mistakes

| Mistake | Fix |
|---|---|
| Looking for `implements` | Just write the methods; add `var _ I = T{}` / `(*T)(nil)` to document it |
| `var c Counter = Tally{}` with pointer-receiver methods | Store a pointer: `&Tally{}` |
| Returning a typed nil pointer as an interface (`var p *T; return p`) | `return nil` explicitly on the success path |
| `if x == nil` to detect a nil pointer *inside* an interface | Compare before converting, or assert: `p, ok := x.(*T); ok && p == nil` |
| `v := x.(T)` on untrusted input | `v, ok := x.(T)` — or a type switch |
| `case int, int64:` and then using `x` as an `int` | Multi-type cases keep the interface type: split the cases |
| Interface case before a specific type that also matches it | Order cases from most specific to most general |
| Passing a `[]Rect` where `[]Shape` is expected | Build a `[]Shape` in a loop |
| Huge interfaces defined next to their only implementation | Small, consumer-side interfaces; return concrete types |
| `==` on interfaces that may hold slices/maps | Only compare comparable dynamic types; otherwise compare manually |
| Using `any` everywhere "for flexibility" | Use a concrete type or a specific interface; keep the compiler's help |

## Key takeaways

- An interface is a **method set**; a type satisfies it **implicitly** by having the methods.
- `var _ I = (*T)(nil)` turns the promise into a compile-time check.
- `T`'s method set has only value-receiver methods; `*T` has both. Pointer-receiver methods ⇒ put `*T` in the interface.
- An interface value is a **(type, value)** pair; it's `nil` only when both are nil. Return a literal `nil`, never a typed nil pointer.
- `any` holds anything; get values back with `v, ok := x.(T)` or a **type switch** (`case nil`, multi-type cases keep the interface type).
- Embed interfaces to combine them; keep them small and named `-er`.
- **Accept interfaces, return structs**, and declare interfaces where they're used.
- `[]Shape` gives polymorphism; sort it with `slices.SortFunc` through interface methods.
- Interface `==` compares type and value; it panics on uncomparable dynamic types.
