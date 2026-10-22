# Generic Methods

Go 1.27 lets a method declare **its own** type parameters. Use them to give `List[T]` chainable `Map`, `FlatMap`, `Reduce` and `GroupBy` methods that return lists of **other** element types. `Sum` needs a tighter constraint on `T` itself, so it stays a function.

```go
type List[T any] struct {
	items []T
}

func Of[T any](vals ...T) List[T]

func (l List[T]) Items() []T
func (l List[T]) Len() int
func (l List[T]) Filter(keep func(T) bool) List[T]
func (l List[T]) Map[U any](f func(T) U) List[U]
func (l List[T]) FlatMap[U any](f func(T) []U) List[U]
func (l List[T]) Reduce[A any](init A, f func(A, T) A) A
func (l List[T]) GroupBy[K comparable](key func(T) K) map[K]List[T]

type Number interface {
	~int | ~int64 | ~float64
}

func Sum[N Number](l List[N]) N
```

## Examples

| Code | Result |
|---|---|
| `Of(1, 2, 3, 4).Filter(isEven).Items()` | `[2 4]` |
| `Of(1, 2, 3).Map(strconv.Itoa).Items()` | `["1" "2" "3"]` (a `List[string]`) |
| `Of(1, 2).Map[float64](half).Items()` | `[0.5 1]` (explicit type argument) |
| `Of("ab", "c").FlatMap(splitChars).Items()` | `["a" "b" "c"]` |
| `Of(1, 2, 3).Reduce("", join)` with `join` appending digits | `"123"` |
| `Of("go", "c", "zig").GroupBy(byLen)` | `map[1:[c] 2:[go] 3:[zig]]` |
| `Sum(Of(1.5, 2.5))` | `4` |
| `Of(1, 2, 3, 4).Filter(isEven).Map(strconv.Itoa).Items()` | `["2" "4"]` (chained) |

## Rules
- `List` is immutable-style: `Of` copies its arguments, and **no method modifies** `l`. Every method returns a new list (or value).
- `Items` returns a **new** slice (empty or nil for an empty list).
- `Map` / `FlatMap` / `Filter` keep the original order. `FlatMap` concatenates the slices `f` returns.
- `Reduce` calls `f` left to right, starting from `init`; an empty list returns `init`.
- `GroupBy` keeps the original order inside each group; an empty list gives an empty (or nil) map.
- `Sum` of an empty list is `0`.
- The zero `List[T]` is an empty list.
