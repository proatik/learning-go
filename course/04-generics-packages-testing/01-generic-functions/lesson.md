# Generic Functions & Constraints

## 1. Why generics: the duplication problem

Without generics, the same algorithm is written once per type:

```go
func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Float64Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
// ...and StringMax, Int64Max, DurationMax...
```

The old workaround is `any` plus type assertions:

```go
func MaxAny(a, b any) any {
	if a.(int) > b.(int) { // only works for int
		return a
	}
	return b
}

n := MaxAny(3, 7).(int) + 1 // the caller must assert the result back
MaxAny(3, 7.5)              // compiles, then panics:
// panic: interface conversion: interface {} is float64, not int
```

| | One func per type | `any` + assertions | Generics |
|---|---|---|---|
| Code written | N copies | 1 | 1 |
| Type checked | at compile time | at **run time** (panics) | at compile time |
| Result type | exact | `any`: caller must assert | exact |
| Speed | fast | boxing + assertions | fast |

With a **type parameter**, you write it once and the compiler checks every use:

```go
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

Max(3, 7)                // 7
Max(2.5, 1.0)            // 2.5
Max("apple", "banana")   // banana
```

You have already **used** generic functions: `slices.Sort`, `slices.Contains`, `maps.Clone`, `cmp.Compare`, `errors.AsType`. This lesson is about **writing** them.

## 2. Type-parameter syntax

```go
func Max[T cmp.Ordered](a, b T) T
//       │ │                 │  └ T used as the result type
//       │ │                 └ T used as a parameter type
//       │ └ the constraint: which types T may be
//       └ the type parameter's name
```

- The **type-parameter list** goes in square brackets between the name and the ordinary parameters.
- Every type parameter has a **constraint**. `any` means "no restriction".
- Inside the body, `T` is a real type: declare variables of it, make slices of it, return it.
- Names: a single capital letter by convention (`T`, `K`, `V`, `E`, `S`). Use a word when it helps (`From`, `To`).
- Parameters with the same constraint can share it, like ordinary parameters: `[T, U any]`.

```go
func Swap[T any](a, b T) (T, T) { return b, a }

func Repeat[T any](v T, n int) []T {
	out := make([]T, n)
	for i := range out {
		out[i] = v
	}
	return out
}
```

## 3. Instantiation and type inference

Calling a generic function **instantiates** it: every type parameter is replaced by a concrete **type argument**, and the result is checked like ordinary code.

```go
Max[int](3, 7)       // explicit type argument
Max(3, 7)            // inferred: T = int (from the arguments)
Max[float64](3, 2.5) // explicit: 3 becomes float64(3)
```

**Inference** reads the types of the arguments. You write `F[T]` yourself only when the compiler can't infer it:

| Situation | Example | Why |
|---|---|---|
| `T` appears only in the **result** | `Zero[string]()` | no argument mentions `T` |
| You want a **different** type than the arguments give | `Max[float64](3, 4)` | untyped `3, 4` would default to `int` |
| Using the function as a **value** | `f := Max[int]` | a generic function can't be used without instantiation |
| Some parameters only in the result | `Convert[float64](7)` | partial: give the first ones, the rest are inferred |

```go
func Zero[T any]() T {
	var zero T
	return zero
}

Zero()          // ❌ in call to Zero, cannot infer T
Zero[string]()  // ✅ ""

f := Max        // ❌ cannot use generic function Max without instantiation
f := Max[int]   // ✅ f is func(int, int) int
```

Arguments must agree on **one** type for `T`:

```go
var i int = 3
var f float64 = 2.5
Max(i, f)   // ❌ type float64 of f does not match inferred type int for T
Max(3, 2.5) // ✅ untyped constants: T = float64, prints 3
```

Inference also works through **function arguments**:

```go
Map([]int{1, 2, 3}, strconv.Itoa) // T = int, U = string (from strconv.Itoa's signature)
```

## 4. Constraints: `any`, `comparable`, `cmp.Ordered`

A constraint is an **interface**. It decides which types are allowed **and** which operations the body may use.

| Constraint | Allowed types | Operations in the body |
|---|---|---|
| `any` | every type | assign, pass, store, return: **no** `==`, `<`, `+` |
| `comparable` | types that support `==` (numbers, strings, bools, pointers, structs/arrays of those, interfaces) | `==`, `!=`, use as a **map key** |
| `cmp.Ordered` | integers, floats, strings (and named types of them) | `<`, `<=`, `>`, `>=`, `==`, `+` on numbers/strings |

The compiler only allows what **every** type in the constraint supports:

```go
func Max[T any](a, b T) T {
	if a > b { // ❌ invalid operation: a > b (type parameter T cannot use operator >)
	...

func Index[T any](xs []T, v T) int {
	if x == v { // ❌ invalid operation: x == v (incomparable types in type set)
	...
```

Fix: pick the **narrowest constraint that allows the operations you use**.

```go
func Index[T comparable](xs []T, target T) int {
	for i, x := range xs {
		if x == target {
			return i
		}
	}
	return -1
}

Index([]string{"a", "b"}, "b") // 1
```

- `comparable` is predeclared. `cmp.Ordered` lives in package `cmp`.
- Slices, maps and funcs are **not** comparable, so `Index([][]int{...}, ...)` doesn't compile.
- Interface types (like `any`) **do** satisfy `comparable`, but `==` panics at run time if the dynamic values aren't comparable:

```go
Index([]any{[]int{1}}, any([]int{1})) // panic: runtime error: comparing uncomparable type []int
```

## 5. Constraint interfaces: type sets and unions

You can write your own constraint. An interface lists the allowed types with `|` (a **union**):

```go
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // ✅ every type in the set supports +
	}
	return total
}

Sum([]int{1, 2, 3})     // 6
Sum([]float64{1.5, 2})  // 3.5
```

- The set of types an interface allows is its **type set**. An operation is allowed only if **every** type in the set supports it.
- Interfaces with a union (or `~`) are **constraints only**. They can't be the type of a variable:

```go
var n Number // ❌ cannot use type Number outside a type constraint: interface contains type constraints
```

- Conversions are allowed if every type can convert: `float64(x)` works for any `T Number`.
- Constraints compose by embedding: `type Real interface{ Integer | Float }`.

## 6. The `~` tilde: underlying types

Every named type has an **underlying type**. For `type Celsius float64`, the underlying type is `float64`.

| Term | Matches |
|---|---|
| `float64` | exactly `float64` |
| `~float64` | `float64` **and** every type whose underlying type is `float64`, like `Celsius` |
| `~int64` | `int64`, `time.Duration`, and any `type X int64` |

```go
type Celsius float64

type Exact interface{ int | int64 | float64 }

func SumExact[T Exact](xs []T) T { ... }

SumExact([]Celsius{20.5, 21})
// ❌ Celsius does not satisfy Exact (possibly missing ~ for float64 in Exact)

Sum([]Celsius{20.5, 21}) // ✅ 41.5, and the result type is Celsius
```

- Almost always write `~`: users of your function will have named types (`Celsius`, `UserID`, `Cents`).
- `cmp.Ordered` already uses `~` for every type: `~int | ~int8 | ... | ~float64 | ~string`.

**Keeping the caller's slice type: `S ~[]E`.** This is the pattern in `slices.Sort[S ~[]E, E cmp.Ordered](x S)`:

```go
type Temps []Celsius

func Double[S ~[]E, E Number](s S) S {
	out := make(S, len(s))
	for i, v := range s {
		out[i] = v * 2
	}
	return out
}

var t Temps = Temps{1, 2}
d := Double(t) // d is Temps, not []Celsius: the caller keeps its named type and its methods
```

## 7. Constraints with methods

A constraint can require methods too, alone or together with a type set:

```go
type Level int

func (l Level) String() string { return "L" + strconv.Itoa(int(l)) }

type IntStringer interface {
	~int            // an integer underneath: allows int(x), *, +
	String() string // and it has a String method
}

func Describe[T IntStringer](xs []T) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		out = append(out, fmt.Sprintf("%s(%d)", x.String(), int(x)*10))
	}
	return out
}

Describe([]Level{1, 2}) // [L1(10) L2(20)]
```

- A constraint with only methods (like `fmt.Stringer`) is an ordinary interface: you can use it as a constraint **and** as a normal type.
- `T fmt.Stringer` vs a plain `fmt.Stringer` parameter: §12.

## 8. Multiple type parameters

```go
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
```

- Each parameter gets the constraint **its** operations need: map keys need `comparable`, values can be `any`.
- Map keys are always `comparable`: `map[K]V` doesn't compile with `K any`.
- Order matters for explicit instantiation. Put the ones callers may need to write **first**:

```go
func Convert[To, From Number](x From) To { return To(x) }

Convert[float64](7) // To = float64 (given), From = int (inferred)
```

## 9. The zero value: `var zero T`

You don't know what `T` is, so you can't write `0`, `""` or `nil`. Declare a variable instead:

```go
func First[T any](xs []T) (T, bool) {
	if len(xs) == 0 {
		var zero T // 0, "", nil, false or an empty struct, depending on T
		return zero, false
	}
	return xs[0], true
}

First([]int{})        // 0 false
First([]string{"a"})  // a true
```

- Pair it with `ok bool` (or an `error`): a zero result alone can't tell "empty" from "the value is zero".
- `*new(T)` also works, but `var zero T` is the idiom.
- `v == zero` needs `T comparable`.

## 10. Generic helpers over slices and maps

The classic toolbox. Learn these shapes; they appear in every codebase.

```go
func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

func Reduce[T, A any](xs []T, init A, f func(A, T) A) A {
	acc := init
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}
```

```go
nums := []int{1, 2, 3, 4}
Filter(nums, func(n int) bool { return n%2 == 0 }) // [2 4]
Map(nums, strconv.Itoa)                            // ["1" "2" "3" "4"]

sumLen := func(n int, s string) int { return n + len(s) }
Reduce([]string{"a", "bb", "ccc"}, 0, sumLen)      // 6: T = string, A = int
```

- **Check the standard library first.** `slices` and `maps` already have `Index`, `Contains`, `Max`, `Min`, `Sort`, `SortFunc`, `Compact`, `Equal`, `Clone`, `maps.Keys` (an iterator, L3)...
- Don't modify the input: build a new slice. The caller may still use the old one.
- A key function `func(T) K` + `cmp.Ordered` / `comparable` is the reusable pattern for "sort by", "group by", "max by".
- A **comparator** `func(a, b T) int` (negative / 0 / positive) works for any `T`, even structs. `cmp.Compare[T]` is the ready-made one for ordered types:

```go
slices.SortFunc(people, func(a, b Person) int { return cmp.Compare(a.Age, b.Age) })
```

## 11. Go 1.26: self-referential constraints

Since Go 1.26 a generic type can mention **itself** in its own type-parameter list. This expresses "a type that can be added to its own type":

```go
type Adder[A Adder[A]] interface {
	Add(A) A
}

func SumAll[A Adder[A]](xs []A) A {
	var total A
	for _, x := range xs {
		total = total.Add(x)
	}
	return total
}

type Money struct{ Cents int }

func (m Money) Add(o Money) Money { return Money{m.Cents + o.Cents} }

SumAll([]Money{{100}, {250}}) // {350}
```

- `A Adder[A]` reads: "`A` must have a method `Add(A) A`", i.e. it takes and returns **its own** type, not some interface.
- Use it for operations on your own structs (`Money.Add`, `Vector.Add`, `Version.Less`). Built-in numbers use a type-set constraint instead.

## 12. When NOT to use generics

Generics remove **duplication across types**. If there is no such duplication, don't add them.

| Situation | Write |
|---|---|
| Only one type is ever used | a plain function: `func Total(orders []Order) int` |
| You only call methods on the value | an interface parameter: `func Print(w io.Writer)` |
| The code is different per type | different functions, or a type switch |
| The same algorithm for many types (containers, slice/map helpers) | ✅ generics |

```go
// ❌ the type parameter adds nothing: io.Writer already works for every writer
func Write[W io.Writer](w W, s string) { io.WriteString(w, s) }

// ✅
func Write(w io.Writer, s string) { io.WriteString(w, s) }
```

One case where a method constraint **does** help: `[]Level` is not a `[]fmt.Stringer`, so a `func(xs []fmt.Stringer)` can't take it. `func[T fmt.Stringer](xs []T)` can.

Rule of thumb (Go blog, *When To Use Generics*): **write the code first, and add a type parameter when you find yourself writing it a second time.**

## Common mistakes

| Mistake | Fix |
|---|---|
| `func Max[T any](a, b T)` then `a > b` | `[T cmp.Ordered]`: the constraint must allow the operator |
| `[T any]` then `x == v` or `map[T]...` | `[T comparable]` |
| `int \| float64` and a caller has `type Celsius float64` | `~int \| ~float64` |
| `var n Number` with a union constraint | union interfaces are constraints only; use `T` inside a generic func |
| `return 0` / `return nil` with result type `T` | `var zero T; return zero` |
| Returning only `T` when "empty" is possible | `(T, bool)`: zero can't signal "not found" |
| `Zero()` with `T` only in the result | `Zero[string]()` |
| `f := Max` | `f := Max[int]` |
| Summing `[]uint8` in `T` and dividing | overflow: convert each value to `float64` first |
| `Filter` returns `[]E` and the caller loses its `Temps` type | `[S ~[]E, E any](s S) S` |
| Rewriting `slices.Contains` / `slices.Max` | check `slices`, `maps`, `cmp` first |
| Generic wrapper around an interface (`[W io.Writer]`) | take the interface directly |

## Key takeaways

- A type parameter `[T C]` lets one function work for many types, checked at **compile time**, with no `any` assertions.
- Calls are **instantiated**. Inference picks `T` from the arguments; write `F[T]` when `T` is only in the result, for a function value, or to override.
- The constraint is an interface that controls both the allowed types and the allowed operations: `any`, `comparable`, `cmp.Ordered`, or your own union.
- `~T` = every type whose underlying type is `T`. Use it, so named types like `Celsius` work.
- Constraints can mix a type set with methods: `interface{ ~int; String() string }`.
- `var zero T` is the zero value of any `T`; pair it with an `ok bool`.
- `S ~[]E` keeps the caller's named slice type. A key func or a `func(a, b T) int` comparator makes helpers work for structs.
- Go 1.26: `type Adder[A Adder[A]] interface{ Add(A) A }` expresses "operates on its own type".
- Don't use generics for a single type or when an interface parameter is enough.
