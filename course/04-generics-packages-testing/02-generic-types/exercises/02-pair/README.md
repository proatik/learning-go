# Pairs and Zip

A struct with **two** type parameters, a method whose result swaps them, and helpers that zip two slices together.

```go
type Pair[A, B any] struct {
	First  A
	Second B
}

func New[A, B any](a A, b B) Pair[A, B]
func (p Pair[A, B]) Swap() Pair[B, A]
func (p Pair[A, B]) String() string

func Zip[A, B any](as []A, bs []B) []Pair[A, B]
func Unzip[A, B any](ps []Pair[A, B]) ([]A, []B)
```

## Examples

| Code | Result |
|---|---|
| `New("age", 42)` | `Pair[string, int]{First: "age", Second: 42}` |
| `New("age", 42).Swap()` | `Pair[int, string]{First: 42, Second: "age"}` |
| `New("age", 42).String()` | `"(age, 42)"` |
| `New(1.5, true).String()` | `"(1.5, true)"` |
| `Zip([]int{1, 2, 3}, []string{"a", "b"})` | `[(1, a) (2, b)]` |
| `Zip([]int{1}, []string(nil))` | empty |
| `Unzip([]Pair[int, string]{{1, "a"}, {2, "b"}})` | `[1 2]`, `["a" "b"]` |

## Rules
- `String` formats both fields with `%v`, as `"(First, Second)"`.
- `Zip` pairs elements by index and stops at the **shorter** slice; the result has exactly `min(len(as), len(bs))` elements (empty or nil when that's 0).
- `Unzip` is the inverse: two slices of `len(ps)` elements each (empty or nil for no pairs).
- `Zip` and `Unzip` must not modify their inputs.
