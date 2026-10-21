# Functions, Closures & defer

## 1. Declaring a function

```go
func name(param1 type1, param2 type2) resultType {
	// body
}

func add(x int, y int) int { return x + y }
func add2(x, y int) int    { return x + y } // consecutive params of one type: group them
func greet()               { fmt.Println("hi") } // no params, no result
```

| Part | Notes |
|---|---|
| `func` keyword | Functions are declared at package level (top level), never inside another `func` |
| Name | Capitalized = exported (L1) |
| Parameters | `name type`; the type comes **after** the name |
| Result | Omitted if nothing is returned; `(T1, T2)` for several |

- There is **no overloading** and **no default parameter values**: every name is unique in its package and every call passes every argument.
- Order of declaration doesn't matter: a function may call one declared further down the file.

## 2. Pass by value

Every argument is **copied** into the parameter. Changing the parameter never changes the caller's variable.

```go
func double(n int) int {
	n *= 2 // modifies the local copy
	return n
}

x := 5
y := double(x) // x is still 5, y is 10
```

(Pointers, which let a function modify the caller's variable, come in a later module.)

## 3. Multiple return values

```go
func divmod(a, b int) (int, int) {
	return a / b, a % b
}

q, r := divmod(17, 5) // 3, 2
q, _ = divmod(9, 2)   // _ (blank identifier) discards a result you don't need
```

You must receive **all** results or none: `q := divmod(17, 5)` doesn't compile. Results can feed another call directly when they match its parameters: `fmt.Println(divmod(17, 5))`.

## 4. Named results and naked `return`

Results may be named. They are ordinary variables, initialized to their **zero value** at function entry.

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // "naked" return: returns the current x, y
}
```

| Use named results when… | Avoid naked `return` when… |
|---|---|
| The names document the results: `(lo, hi int)`, `(q, r int, err error)` | The function is longer than a few lines: the reader must scroll to see what is returned |
| A deferred closure needs to modify a result (§12) | Results get shadowed by an inner `:=` (compile error: "result parameter not in scope at return") |

Rule of thumb from the Go docs: naked returns are fine in **short** functions; otherwise write `return x, y` explicitly, even when the results are named.

## 5. The `(T, error)` convention

Functions that can fail return an `error` as the **last** result. `nil` means success. You met this with `strconv.Atoi` in L2.

```go
import (
	"errors"
	"fmt"
)

func sqrtInt(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("sqrtInt: negative input %d", n) // formatted message
	}
	if n > 1_000_000 {
		return 0, errors.New("sqrtInt: input too large") // fixed message
	}
	r := 0
	for (r+1)*(r+1) <= n {
		r++
	}
	return r, nil
}

r, err := sqrtInt(-4)
if err != nil {
	fmt.Println("error:", err) // error: sqrtInt: negative input -4
	return
}
fmt.Println(r)
```

- Check `err` **immediately**, handle it, and return early. The success path stays unindented.
- On error, return the **zero value** for the other results.
- Error strings start lowercase and have no trailing punctuation: they are often embedded in longer messages.
- `errors.New` for fixed text, `fmt.Errorf` when you need verbs. (Wrapping errors with `%w`, `errors.Is`/`As` come in a later module.)

## 6. Variadic functions and a slice preview

A final parameter `...T` accepts **zero or more** arguments. Inside the function it is a **slice** `[]T`: a list you can `range` over.

```go
func sum(nums ...int) int {
	total := 0
	for _, n := range nums { // range over a slice gives index, value
		total += n
	}
	return total
}

sum()        // 0 (nums is empty)
sum(1, 2, 3) // 6

xs := []int{1, 2, 3} // slice literal (slices are covered properly in a later lesson)
sum(xs...)           // spread an existing slice into the variadic parameter
```

- Only the **last** parameter may be variadic.
- `len(nums)` gives the number of arguments.
- To demand at least one argument, make it a regular parameter: `func max(first int, rest ...int) int`.
- `fmt.Println(a ...any)` is variadic: that is why it takes any number of arguments.

## 7. Function values and function types

Functions are **values**: assign them to variables, pass them, return them. The zero value of a function type is `nil`; calling a nil function crashes the program.

```go
func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }

type BinOp func(int, int) int // named function type

var op BinOp = add
fmt.Println(op(3, 4)) // 7
op = sub
fmt.Println(op(3, 4)) // -1

var f func(int) int  // nil
if f != nil {        // function values may only be compared to nil
	f(1)
}
```

A function's type is its **signature**: parameter and result types. Parameter names don't matter: `func(a, b int) int` and `func(int, int) int` are the same type.

## 8. Higher-order functions

```go
// apply takes a function as a parameter.
func apply(op BinOp, a, b int) int { return op(a, b) }

// adder returns a function.
func adder(n int) func(int) int {
	return func(x int) int { return x + n }
}

apply(add, 2, 3) // 5
add10 := adder(10)
add10(5)         // 15
```

## 9. Anonymous functions (function literals)

A function without a name, written as an expression. Use it inline, assign it, or call it immediately.

```go
square := func(x int) int { return x * x }
square(4) // 16

apply(func(a, b int) int { return a * b }, 6, 7) // 42

func() {
	fmt.Println("called immediately")
}() // the trailing () calls it
```

## 10. Closures

A function literal can refer to variables of the surrounding function. It **captures the variables themselves** (by reference), not copies of their values, and keeps them alive as long as the closure exists.

```go
func counter() func() int {
	count := 0 // lives on after counter returns
	return func() int {
		count++
		return count
	}
}

c1 := counter()
fmt.Println(c1(), c1(), c1()) // 1 2 3
c2 := counter()               // a new call → a new, independent count
fmt.Println(c2())             // 1
fmt.Println(c1())             // 4
```

```go
x := 1
show := func() { fmt.Println(x) }
x = 2
show() // 2 — the closure sees the variable, not the value at creation time
```

Since Go 1.22, each iteration of a `for` loop has its **own** loop variable, so closures created in a loop each see their own `i`.

## 11. Recursion

A function may call itself. Always have a **base case** that stops the recursion.

```go
func factorial(n int) int {
	if n <= 1 {
		return 1 // base case
	}
	return n * factorial(n-1) // recursive case: a smaller problem
}
```

- Go has **no tail-call optimization**: every call uses stack space. Go's stacks grow dynamically (up to 1 GB on 64-bit), so deep recursion works, but runaway recursion ends in `fatal error: stack overflow`.
- When a simple loop does the job, the loop is usually clearer and faster in Go.
- Naive recursion can repeat work: `fib(n-1) + fib(n-2)` is exponential. Prefer a loop or a closure that keeps state.

## 12. `defer`

`defer f(args)` schedules a call to run **when the surrounding function returns**: after the `return` statement sets the results, whether the function returns normally or panics.

```go
func main() {
	fmt.Println("start")
	defer fmt.Println("deferred 1")
	defer fmt.Println("deferred 2")
	fmt.Println("end")
}
// start
// end
// deferred 2
// deferred 1
```

Three rules (from the blog post *Defer, Panic, and Recover*):

| Rule | Example |
|---|---|
| 1. **Arguments are evaluated when `defer` runs**, not when the call runs | `i := 1; defer fmt.Println(i); i = 2` prints `1` |
| 2. Deferred calls run in **LIFO** order (last deferred, first run) | `for i := range 3 { defer fmt.Print(i) }` prints `210` |
| 3. A deferred **closure can read and modify named results** | see below |

```go
func triple() (n int) {
	defer func() { n *= 3 }() // runs after "return 5" set n = 5
	return 5
}
// triple() == 15
```

Typical uses: cleanup right next to the acquisition (`f, err := os.Open(...)` then `defer f.Close()`), unlocking, tracing, and annotating returned errors:

```go
func parsePort(s string) (port int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsePort(%q): %v", s, err)
		}
	}()
	port, err = strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return port, nil
}
```

`defer` runs at **function** exit, not block or loop-iteration exit: deferring inside a long loop piles up calls until the function returns.

## 13. `panic` and `recover` (preview)

`panic` stops normal execution (as a nil function call or an out-of-range index does); deferred calls still run; `recover`, called inside a deferred function, can stop the panic. They are for truly exceptional situations, not ordinary errors, and are covered in a later module. Use `error` results for anything that can be expected to fail.

## Common mistakes

| Mistake | Fix |
|---|---|
| `x := f()` when `f` returns two values | `x, err := f()` or `x, _ := f()` |
| Ignoring `err` with `_` "for now" | Check it: `if err != nil { return ..., err }` |
| Returning a non-zero value alongside an error | Return zero values: `return 0, err` |
| `Error strings like this.` | `error strings like this` (lowercase, no period) |
| Expecting a function to change the caller's `int` | Arguments are copies: return the new value instead |
| Calling a `nil` function value | Check `if f != nil` or make sure it is always set |
| Long function with naked `return` | Write the results explicitly |
| Forgetting the base case in recursion | Stack overflow; handle the smallest input first |
| `defer fmt.Println(x)` expecting the final `x` | Arguments are evaluated immediately; use `defer func() { fmt.Println(x) }()` |
| `defer` inside a loop of thousands of iterations | Move the loop body into its own function |
| `sum(xs)` where `sum(...int)` and `xs` is `[]int` | Spread it: `sum(xs...)` |

## Key takeaways

- `func name(a, b int) (int, error)`: types after names; group same-typed params; multiple results are normal.
- Arguments are passed **by value** (copied).
- Failure is a value: return `(T, error)`, check `err != nil` right away, use `errors.New` / `fmt.Errorf`.
- Named results document intent; keep naked `return` for short functions.
- `...T` collects arguments into a slice; `s...` spreads a slice into them.
- Functions are values with types like `func(int) int`; pass, return and store them.
- Closures capture variables by reference; each call to the outer function creates fresh state.
- Recursion needs a base case; Go has no tail-call optimization.
- `defer`: LIFO, arguments evaluated at once, can modify named results; runs at function exit.
