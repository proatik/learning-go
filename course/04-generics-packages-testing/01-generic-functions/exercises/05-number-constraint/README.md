# A `Number` Constraint with `~`

Write your own constraint, then use it for arithmetic that works on `int`, `float64` **and** named types like `Celsius`.

```go
type Number interface {
	// TODO: replace `any` with a union
}

func Sum[T Number](xs []T) T
func Abs[T Number](x T) T
func Scale[S ~[]E, E Number](s S, k E) S
```

## Examples

| Call | Result |
|---|---|
| `Sum([]int{1, 2, 3})` | `6` |
| `Sum([]Celsius{20.5, 21})` | `Celsius(41.5)` |
| `Sum([]int8{})` | `0` |
| `Abs(-7)`, `Abs(Celsius(-3.5))` | `7`, `Celsius(3.5)` |
| `Scale(Readings{1, 2.5}, 2)` | `Readings{2, 5}` (type `Readings`, not `[]Celsius`) |

```go
type Celsius float64
type Readings []Celsius
```

## Rules
- `Number` must allow **every signed integer and float type and every named type based on them**: `int`, `int8`, `int16`, `int32`, `int64`, `float32`, `float64`.
  - The starter's `any` is a placeholder so the file compiles. With `any`, `+` and `<` won't compile in your bodies.
  - The tests pass `Celsius` (`float64` underneath) and `Score` (`int` underneath).
- `Sum`: the total, computed in `T` (`0` for nil or empty).
- `Abs`: `-x` when `x < 0`, otherwise `x`.
- `Scale`: a **new** slice of type `S` where element `i` is `s[i] * k`. Its length is `len(s)`. Don't modify `s`.
