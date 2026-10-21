# Variables, Types & Constants

## 1. Declaring variables with `var`

```go
var name string             // type, no value  → zero value ""
var count int = 10          // type and value
var ratio = 0.5             // value only      → type inferred: float64
var x, y int = 1, 2         // several at once
var a, b = 3, "three"       // different inferred types are fine

var ( // grouped block — common at package level
	host    = "localhost"
	port    = 8080
	verbose bool
)
```

`var` works **everywhere**: at package level and inside functions.

## 2. Short declaration `:=`

```go
func main() {
	n := 42          // same as: var n = 42
	s, ok := "hi", true
	n, m := 7, 8     // OK: m is new, n is just assigned
	// n := 9        // ❌ no new variables on left side of :=
}
```

| Rule | Detail |
|---|---|
| Only **inside functions** | at package level you must use `var` |
| At least **one new** variable on the left | the others (in the same scope) are simply assigned |
| Type is always **inferred** | you can't write `x int := 1` |

Idiom: use `:=` for most locals; use `var x T` when you want the zero value or a specific type (`var f float64 = 1`).

`=` is plain **assignment** to existing variables. Multiple assignment evaluates all right-hand sides first, so swapping needs no temp:

```go
a, b = b, a
```

## 3. Zero values

A variable declared without a value is **never garbage** — it holds its type's zero value.

| Type | Zero value |
|---|---|
| numeric (`int`, `float64`, …) | `0` |
| `bool` | `false` |
| `string` | `""` (empty) |
| pointers, slices, maps, channels, funcs, interfaces | `nil` |

```go
var i int; var f float64; var s string; var ok bool
fmt.Printf("%v %v %q %v\n", i, f, s, ok) // 0 0 "" false
```

## 4. Basic types

| Type | Size | Range / notes |
|---|---|---|
| `bool` | — | `true`, `false` |
| `string` | — | immutable sequence of bytes, usually UTF-8 |
| `int8` | 8 bit | −128 … 127 |
| `int16` | 16 bit | −32 768 … 32 767 |
| `int32` | 32 bit | ≈ ±2.1 × 10⁹ |
| `int64` | 64 bit | ≈ ±9.2 × 10¹⁸ |
| `int` | 32 or 64 bit | platform word size (64 on all modern machines) — **the default integer** |
| `uint8` … `uint64`, `uint` | same sizes | 0 … 2ⁿ−1 (e.g. `uint8`: 0 … 255) |
| `byte` | 8 bit | alias for `uint8` — raw data |
| `rune` | 32 bit | alias for `int32` — one Unicode code point, e.g. `'é'` |
| `float32` | 32 bit | ~7 decimal digits |
| `float64` | 64 bit | ~15 digits — **the default float** |
| `complex128` | 128 bit | `complex(1, 2)` or `1 + 2i`; rarely needed |
| `uintptr` | word | low-level; ignore for now |

Use `int` and `float64` unless you have a reason (file formats, memory layout, APIs). The `math` package has the limits: `math.MaxInt8`, `math.MinInt64`, `math.MaxUint8`, `math.MaxInt`, `math.MaxFloat64`…

`%T` prints a value's type: `fmt.Printf("%T\n", 'a') // int32`.

## 5. Conversions: `T(v)`, never implicit

Go **never** converts between types automatically — not even `int` → `int64` or `int` → `float64`.

```go
i := 7
f := 2.5
// x := i * f            // ❌ mismatched types int and float64
x := float64(i) * f      // 17.5
n := int(f)              // 2  — float → int truncates toward zero
neg := -2.9
m := int(neg)            // -2 — toward zero, not down
// k := int(-2.9)        // ❌ "constant -2.9 truncated": constants must convert exactly
s := string(rune(65))    // "A" — but for numbers → text use strconv (section 11)
```

| From → To | Behaviour |
|---|---|
| float → int | drops the fraction (toward zero): with `f := 3.99`, `int(f) == 3`; with `f := -3.99`, `int(f) == -3` |
| int → float | exact for "small" ints, may round for huge ones |
| bigger int → smaller int | keeps the low bits: `int8(200) == -56` (variable only) |
| `string(i)` on an int | gives a **character**, not digits — `go vet` warns |

## 6. Constants

```go
const Pi = 3.14159             // untyped floating-point constant
const MaxUsers int = 100       // typed constant
const (
	Greeting = "hello"
	Answer   = 42
)
```

- Constant values must be known **at compile time**: numbers, strings, booleans, and expressions of them (`len("abc")`, `1 << 10`). No function calls like `math.Sqrt(2)`.
- Constants can't be changed and can't have their address taken.

### Untyped constants

An untyped constant has a **kind**, not yet a type, and is *exact* (arbitrary precision). It takes the type its context needs:

```go
const big = 1 << 100           // fine — never stored as a variable
const half = 0.5

var f32 float32 = half         // half becomes float32
var d time.Duration = 2        // 2 becomes time.Duration
fmt.Println(big >> 98)         // 4
```

If the context has no type (e.g. `x := c`), the **default type** is used:

| Untyped constant | Default type |
|---|---|
| `42` (integer) | `int` |
| `4.2`, `1e3` (floating-point) | `float64` |
| `'x'` (rune) | `rune` (`int32`) |
| `2i` (complex) | `complex128` |
| `"s"` | `string` |
| `true` | `bool` |

Integer-only constant expressions use **integer division**: `5 / 9 == 0` but `5.0 / 9 == 0.5555…`.

## 7. `iota` — enumerated constants

Inside a `const ( … )` block, `iota` is `0` on the first line and increases by 1 on each line. A line with no expression **repeats** the previous one.

```go
type Weekday int // a named type: its own type, underlying type int

const (
	Sunday Weekday = iota // 0
	Monday                // 1  (= Weekday(iota) repeated)
	Tuesday               // 2
	// …
)

const (
	_  = iota             // 0 — skip with the blank identifier
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20
	GB                    // 1 << 30
)

const (
	Read  = 1 << iota // 1   — bit flags
	Write             // 2
	Exec              // 4
)
```

Named types give type safety: a `Weekday` is not an `int`, so mixing them needs a conversion. Constants convert implicitly: `var d Weekday = 3` is fine, `var d Weekday = n` (`n int`) is not.

## 8. Integer arithmetic gotchas

### Overflow wraps around (no panic, no error)

```go
var x int8 = math.MaxInt8 // 127
x++
fmt.Println(x)            // -128
var u uint8 = 0
u--
fmt.Println(u)            // 255
```

Constants are checked at compile time instead: `var y int8 = 128` → *constant 128 overflows int8*.

To detect overflow, compute in a **wider type** and compare with the limits.

### Integer division and `%`

`/` on integers **truncates toward zero**; `%` takes the **sign of the left operand** (`a == (a/b)*b + a%b` always holds).

| Expression | Result |
|---|---|
| `7 / 2` | `3` |
| `-7 / 2` | `-3` |
| `7 % 3` | `1` |
| `-7 % 3` | `-1` |
| `7 % -3` | `1` |
| `x / 0` (variable) | runtime **panic** |

Need a float result? Convert *before* dividing: `float64(sum) / float64(n)`, not `float64(sum / n)`.

## 9. The blank identifier `_`

`_` is a write-only placeholder: assign to it to **discard** a value.

```go
_, err := strconv.Atoi("12") // keep only the error
n, _ := strconv.Atoi("12")   // keep only the number (only when failure is impossible!)
const (
	_ = iota                  // skip a constant value
	One
)
```

You can never *read* `_`; `x := _` is a compile error.

## 10. Scope and shadowing

Each `{ }` block opens a new scope. `:=` in an inner block **declares a new variable** that hides (shadows) an outer one of the same name.

```go
count := 0
if true {
	count := 5        // NEW variable, only inside this block
	fmt.Println(count) // 5
}
fmt.Println(count)     // 0  ← the outer one never changed
```

Fix: use `=` when you mean to assign to the outer variable. This bug is most common with `x, err := f()` inside an `if` block.

## 11. `strconv` and the `(value, error)` pattern

Numbers ↔ text uses package `strconv`:

| Function | Example | Result |
|---|---|---|
| `strconv.Itoa(i int) string` | `strconv.Itoa(-42)` | `"-42"` |
| `strconv.Atoi(s string) (int, error)` | `strconv.Atoi("42")` | `42, nil` |
| `strconv.ParseFloat(s, 64) (float64, error)` | `strconv.ParseFloat("2.5", 64)` | `2.5, nil` |
| `strconv.FormatInt(i int64, base int) string` | `strconv.FormatInt(255, 2)` | `"11111111"` |
| `strconv.ParseInt(s, base, bits) (int64, error)` | `strconv.ParseInt("ff", 16, 64)` | `255, nil` |
| `strconv.FormatFloat(f, 'f', prec, 64)` | `strconv.FormatFloat(3.14159, 'f', 2, 64)` | `"3.14"` |

Functions that can fail return an **`error` as the last result**. `nil` means success. Check it immediately:

```go
// Double parses s and returns twice its value.
func Double(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err // pass the error up; return the zero value for the rest
	}
	return n * 2, nil
}
```

Make your own error with `fmt.Errorf("port %d out of range", p)` (or `errors.New("msg")`). Errors get a full lesson later.

## Common mistakes

| Mistake | Fix |
|---|---|
| `x := 1` at package level | use `var x = 1` |
| `x := 2` again in the same scope | `x = 2` (assignment) |
| Declaring a variable and never using it | compile error — use it or delete it |
| `intVar * floatVar` | convert one side: `float64(intVar) * floatVar` |
| `float64(a / b)` with ints | `float64(a) / float64(b)` |
| `string(65)` to get `"65"` | `strconv.Itoa(65)` |
| Inner `v, err := …` hides outer `v` | declare `var err error` and use `=` |
| Expecting overflow to panic | it silently wraps — check limits yourself |
| `-7 % 3 == 2` (as in Python) | in Go it is `-1` |
| Ignoring the `error` from `Atoi` | always `if err != nil { … }` |

## Key takeaways

- `var` anywhere, `:=` only in functions; every variable starts at its **zero value**.
- `int`, `float64`, `string`, `bool`, `rune`, `byte` cover almost everything.
- **No implicit conversions** — write `T(v)`; float → int truncates toward zero.
- Constants are compile-time and untyped by default; `iota` builds enums and bit sizes.
- Integer `/` truncates, `%` follows the dividend's sign, overflow wraps silently.
- Watch for `:=` shadowing inside blocks.
- Fallible functions return `(value, error)`; check `err != nil` right away.
