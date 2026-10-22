# Generic Types & Methods

## 1. Why generic types

Lesson 1 made **functions** generic. Data structures need the same thing: a stack of `int`, of `string`, of `*User`, all with one implementation.

| Approach | Problem |
|---|---|
| `type IntStack []int`, `type StringStack []string`, … | the same code copied per element type |
| `type Stack []any` | no compile-time checking; every `Pop` needs a type assertion that can fail at run time |
| `type Stack[T any] struct{ … }` | ✅ one implementation, fully type-checked, no assertions |

```go
// The []any version: anything goes in, so you must check what comes out.
var s []any
s = append(s, 1, "two") // mixing types compiles fine — a bug waiting to happen
n := s[0].(int)         // a type assertion every time you read
fmt.Println(n + 1)      // 2
```

## 2. Declaring a generic type

A **type parameter list** goes in square brackets right after the type name, exactly like a generic function.

```go
// Stack is a last-in, first-out stack of T values.
type Stack[T any] struct {
	items []T
}

// Pair holds two values of possibly different types.
type Pair[K comparable, V any] struct {
	Key K
	Val V
}

// Node is a linked-list node; it refers to its own instantiation.
type Node[T any] struct {
	Val  T
	Next *Node[T]
}
```

- `T`, `K`, `V` can be used anywhere inside: field types, slices `[]T`, maps `map[K]V`, pointers `*Node[T]`.
- Several parameters with the same constraint can share it: `type Pair[K, V any]`.
- Any type can be generic, not just structs: `type List[T any] []T`, `type Set[T comparable] map[T]struct{}`, `type Pred[T any] func(T) bool`.

## 3. Instantiation and zero values

A generic type is a **template**, not a type. You must **instantiate** it with type arguments before you can use it.

```go
var s Stack[string]              // zero value: items is nil, ready to use
nums := Stack[int]{}             // composite literal
p := Pair[string, int]{"age", 42}
n := &Node[string]{Val: "a", Next: &Node[string]{Val: "b"}}

fmt.Printf("%T\n", p)  // main.Pair[string,int]
fmt.Printf("%+v\n", p) // {Key:age Val:42}
fmt.Println(n.Next.Val) // b
```

| You write | Result |
|---|---|
| `var s Stack` | ❌ `cannot use generic type Stack[T any] without instantiation` |
| `p := Pair{"a", 1}` | ❌ same error: type arguments are **never inferred** for a type (only for function calls) |
| `var b Stack[string] = a` with `a` a `Stack[int]` | ❌ `cannot use a (variable of struct type Stack[int]) as Stack[string] value` |
| `type Pair[K comparable, V any]` then `Pair[[]int, string]{}` | ❌ `[]int does not satisfy comparable` |

- Each instantiation is a **distinct type**: `Stack[int]` and `Stack[string]` share code, not identity.
- Zero values follow the field types: `var z Pair[string, *int]` is `{Key:"", Val:nil}`.
- Inside generic code, the zero value of `T` is `var zero T` (from lesson 1). You'll write it in every `Pop` / `Get` that can fail.

## 4. Methods on a generic type

The receiver **re-states the type parameters** (names only, no constraints). Inside the method, `T` means "whatever this instance was created with".

```go
// Push puts v on top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top value; ok is false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Len returns the number of values on the stack.
func (s *Stack[T]) Len() int { return len(s.items) }

var s Stack[string]
s.Push("a")
s.Push("b")
v, ok := s.Pop()
fmt.Println(v, ok, s.Len()) // b true 1
```

| Receiver | Meaning |
|---|---|
| `func (s *Stack[T]) Push(v T)` | ✅ the usual form |
| `func (s *Stack[E]) Push(v E)` | ✅ the names may differ from the type declaration (but keep them the same) |
| `func (s *Stack[_]) Len() int` | ✅ `_` when the method doesn't use the parameter |
| `func (s *Stack[T comparable]) Push(v T)` | ❌ `syntax error: unexpected name comparable, expected ]`: constraints live **only** on the type declaration |
| `func (s Stack) Len() int` | ❌ `cannot use generic type Stack[T any] without instantiation` |

- Value vs pointer receivers work exactly as in Module 3: `Push` changes the stack, so every method gets `*Stack[T]` for consistency.
- ⚠️ **No specialization.** There is no "method only for `Stack[int]`". In `func (s Stack[int]) Sum() int`, `int` is just a **new name for the type parameter** that shadows the real `int`:

```go
func (s Stack[int]) Sum() int { // "int" here is a type parameter, not the int type!
	total := 0
	for _, v := range s.items {
		total += v // ❌ invalid operation: total += v (mismatched types int and int /* with int declared at … */)
	}
	return total
}
```

  Write a generic **function** with a tighter constraint instead (§5).

## 5. Constraints on generic types

A type's constraint limits what it can be instantiated with **and** what its methods may do with `T`.

```go
// Set is a set of comparable values (map keys must be comparable).
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T)      { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool { _, ok := s[v]; return ok }

// SortedList keeps its values in ascending order (needs <, so cmp.Ordered).
type SortedList[T cmp.Ordered] struct {
	items []T
}

// Insert adds v at its sorted position.
func (l *SortedList[T]) Insert(v T) {
	i, _ := slices.BinarySearch(l.items, v)
	l.items = slices.Insert(l.items, i, v)
}

s := Set[string]{}
s.Add("go")
fmt.Println(s.Has("go"), s.Has("rust")) // true false

var l SortedList[int]
for _, v := range []int{5, 1, 3} {
	l.Insert(v)
}
fmt.Println(l.items) // [1 3 5]
```

- Methods get **exactly** the type's constraint. With `any`, `+`, `<` and `==` are rejected: `invalid operation: operator + not defined on s (variable of type T constrained by any)`.
- A method **cannot tighten** the constraint. `Set[T comparable]` can't have a `Sorted()` method that needs `<` (`T does not satisfy cmp.Ordered`). Use a package-level function, which may ask for more:

```go
// Sorted returns the elements of s in ascending order.
func Sorted[T cmp.Ordered](s Set[T]) []T {
	out := make([]T, 0, len(s))
	for v := range s {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

fmt.Println(Sorted(Set[int]{3: {}, 1: {}, 2: {}})) // [1 2 3]
```

- `Sorted(…)` infers `T` from its argument: inference works for **function calls**, including calls with a generic type argument.

## 6. Constructors

A constructor is a generic function, so it follows lesson 1's inference rules.

```go
// NewStack returns a stack holding vals (bottom first).
func NewStack[T any](vals ...T) *Stack[T] {
	return &Stack[T]{items: slices.Clone(vals)}
}

a := NewStack(1, 2, 3)  // T inferred from the arguments: *Stack[int]
b := NewStack[string]() // no arguments → T must be explicit
fmt.Println(a.Len(), b.Len()) // 3 0
```

- `NewStack()` with no arguments → ❌ `in call to NewStack, cannot infer T`.
- Prefer a **useful zero value** (`var s Stack[int]` just works) and add a constructor only for validation or non-zero setup (capacity, comparison function, …).

## 7. Generic interfaces

Interfaces can have type parameters too. Each instantiation is an ordinary interface.

```go
// Container is anything that stores T values.
type Container[T any] interface {
	Add(T)
	Len() int
}

// Bag is an unordered collection.
type Bag[T any] struct{ items []T }

func (b *Bag[T]) Add(v T)  { b.items = append(b.items, v) }
func (b *Bag[T]) Len() int { return len(b.items) }

// Fill adds every value to c.
func Fill[T any](c Container[T], vals ...T) {
	for _, v := range vals {
		c.Add(v)
	}
}

b := &Bag[string]{}
Fill(b, "x", "y") // T = string inferred from *Bag[string]'s methods
var c Container[string] = b
fmt.Println(c.Len()) // 2

var _ Container[int] = (*Bag[int])(nil) // compile-time check that *Bag[int] implements Container[int]
```

- `*Bag[int]` implements `Container[int]`, not `Container[string]`: the method signatures must match **after** substitution.
- Don't confuse this with a **constraint** interface like lesson 1's `Number`: `Container[T]` lists methods, so it's an ordinary interface type you can store values in.

## 8. Embedding a generic type

Embed an **instantiated** type (or one instantiated with your own type parameter). Fields and methods are promoted as usual (Module 3).

```go
// Named embeds a concrete instantiation.
type Named struct {
	Stack[int] // the field is called Stack (no type arguments)
	Name       string
}

// Logged wraps any stack and shadows Push.
type Logged[T any] struct {
	*Stack[T]
	pushes int
}

// Push counts the push, then delegates to the embedded stack.
func (l *Logged[T]) Push(v T) {
	l.pushes++
	l.Stack.Push(v)
}

var n Named
n.Push(3)       // promoted from Stack[int]
n.Stack.Push(4) // explicit form
fmt.Println(n.Len()) // 2

lg := Logged[string]{Stack: &Stack[string]{}}
lg.Push("a")
lg.Push("b")
fmt.Println(lg.Len(), lg.pushes) // 2 2
```

- The field name is the type name **without** type arguments: `n.Stack`, `lg.Stack`.
- Embedding promotes methods, so `*Named` satisfies `interface{ Push(int) }` too.

## 9. Generic type aliases (Go 1.24)

An **alias** (`=`) gives another name to an existing type. Since Go 1.24 an alias may have type parameters.

```go
type Set[T comparable] = map[T]struct{} // alias: Set[string] IS map[string]struct{}

s := Set[string]{"a": {}}
var m map[string]struct{} = s // identical types: no conversion
fmt.Printf("%T %d\n", m, len(m)) // map[string]struct {} 1
```

| | `type Set[T comparable] map[T]struct{}` (defined type) | `type Set[T comparable] = map[T]struct{}` (alias) |
|---|---|---|
| New type? | yes, distinct from the map type | no, just another name |
| Can have methods? | ✅ | ❌ `cannot define new methods on generic alias type Set[T comparable]` |
| Typical use | your own data structure | shorter names; moving a type between packages without breaking callers |

Rule of thumb: you want methods, so you almost always want a **defined** type.

## 10. Generic methods (Go 1.27)

Since **Go 1.27** a method may declare **its own** type parameters, after the method name. This lets a `List[T]` have a `Map` that returns a `List[U]` of a **different** element type, a thing a method could not express before.

```go
// List is an immutable-style list of T values.
type List[T any] struct{ items []T }

// Map returns a new list holding f(v) for every v.
func (l List[T]) Map[U any](f func(T) U) List[U] {
	out := List[U]{items: make([]U, 0, len(l.items))}
	for _, v := range l.items {
		out.items = append(out.items, f(v))
	}
	return out
}

// Reduce folds the list into one value, starting from init.
func (l List[T]) Reduce[A any](init A, f func(A, T) A) A {
	acc := init
	for _, v := range l.items {
		acc = f(acc, v)
	}
	return acc
}

l := List[int]{items: []int{1, 2, 3}}
s := l.Map(strconv.Itoa) // U inferred: string
fmt.Printf("%T %q\n", s, s.items) // main.List[string] ["1" "2" "3"]

half := l.Map[float64](func(n int) float64 { return float64(n) / 2 }) // explicit U
fmt.Println(half.items) // [0.5 1 1.5]

sum := l.Reduce(0, func(acc, n int) int { return acc + n })
fmt.Println(sum) // 6
```

- `T` comes from the **receiver**, `U` / `A` from the **method**. A method type parameter can have any constraint: `func (l List[T]) GroupBy[K comparable](key func(T) K) map[K]List[T]`.
- It also works on non-generic types. The standard library's first one: `func (r *Rand) N[Int intType](n Int) Int` in `math/rand/v2`, so `r.N(int64(10))` returns an `int64`.
- Inference works like lesson 1; give the arguments explicitly when there's nothing to infer from: `l.Map[float64](…)`.
- Method values and expressions need every type parameter fixed:

```go
toStr := l.Map[string]       // ✅ method value: func(func(int) string) List[string]
mapB := List[int].Map[bool]  // ✅ method expression: func(List[int], func(int) bool) List[bool]
f := l.Map                   // ❌ cannot use generic function l.Map without instantiation
```

**The restrictions** (all verified with Go 1.27.3):

| You try | Result |
|---|---|
| a type parameter on an **interface method**: `interface{ Map[U any](func(int) U) U }` | ❌ `interface method must have no type parameters` |
| a generic method satisfying an interface: `List[int]` as `interface{ Map(func(int) string) List[string] }` | ❌ `List[int] does not implement …(wrong type for method Map)` — `have Map[U any](func(int) U) List[U]`, `want Map(func(int) string) List[string]` |
| reusing the receiver's name: `func (l List[T]) Bad[T any]()` | ❌ `T redeclared in this block` |
| `l.Map(nil)`: nothing to infer `U` from | ❌ `in call to l.Map, cannot infer U` |

Method or function?

- **Method** when the operation belongs to the type and reads well chained: `list.Filter(even).Map(strconv.Itoa)`.
- **Function** when it needs a tighter constraint on the receiver's `T` (`Sorted[T cmp.Ordered](Set[T])`), or when it's shared by several types (`slices.Sort`).
- A method type parameter is a **new** parameter: it can't tighten the receiver's `T` (`List[string]` can't get a `Sum`).
- Generic methods never take part in interfaces, so an interface can't promise "has a `Map`". Keep interface methods non-generic.

Next lesson: `All()` **iterators**, so callers can `range` over your generic collections.

## 11. When a generic type beats `[]any` (and when it doesn't)

| Use a **generic type** when… | Use an **interface** (or `any`) when… |
|---|---|
| the container stores values but never calls methods on them: stacks, queues, sets, caches, trees | the values really are different types that share **behavior** (`io.Reader`, `Shape`) |
| you want the compiler to reject wrong element types | one collection must hold mixed types on purpose (e.g. JSON-like data) |
| callers shouldn't need type assertions | the code calls methods and doesn't care about the concrete type |

Don't make a type generic "just in case": start concrete, add type parameters when a **second** element type shows up.

## Common mistakes

| Mistake | Fix |
|---|---|
| `var s Stack` / `Pair{"a", 1}` | Instantiate: `Stack[int]`, `Pair[string, int]{"a", 1}`; types are never inferred |
| Constraint in the receiver: `func (s *Stack[T any]) …` | Receivers list names only: `func (s *Stack[T]) …` |
| `func (s Stack[int]) Sum()` to "specialize" for int | `int` became a type-parameter name. Write `func Sum[T Number](s Stack[T]) T` |
| A method that needs `<` on a `Set[T comparable]` | Package function with a tighter constraint: `Sorted[T cmp.Ordered](s Set[T])` |
| Returning `nil` for "not found" from a `T` result | `var zero T; return zero, false` (`nil` isn't a valid `T` for `any`) |
| `NewStack()` with no arguments | `NewStack[int]()`, or rely on a useful zero value |
| Adding methods to a generic alias `type Set[T comparable] = map[T]struct{}` | Make it a defined type (no `=`) |
| Putting type parameters on an interface method | Not allowed; put them on the interface: `Mapper[T, U any]` |
| Expecting a generic method to implement an interface | It never does; add a non-generic method for the interface |
| `[]any` "to be flexible" | A generic type: same flexibility, checked at compile time |

## Key takeaways

- `type Name[T C] …` declares a generic type; **instantiate** it (`Stack[int]`) before use. Type arguments are inferred for function calls, never for types.
- Each instantiation is a distinct type; zero values follow the field types, and `var zero T` gives the zero `T`.
- Methods re-state the type parameters in the receiver, names only: `func (s *Stack[T]) Push(v T)`. No constraints, no specialization.
- A type's constraint is all its methods get. Need more? Write a generic **function** with a tighter constraint.
- Interfaces can be generic (`Container[T]`); a type implements a specific instantiation.
- Embed instantiated generic types; the field name is the bare type name.
- Generic **aliases** (`=`, Go 1.24) are other names for a type and can't have methods.
- **Go 1.27 generic methods** add their own type parameters (`Map[U any]`), but interface methods can't have them and generic methods never implement interfaces.
