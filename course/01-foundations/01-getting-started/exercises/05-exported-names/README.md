# Exported Names

The tests in this exercise live in a **separate test package** (`package shapes_test`). That means they can only use what package `shapes` **exports**, exactly as any other package would.

Right now the code doesn't even compile against the tests. Fix `shapes.go` so that package `shapes` exports:

```go
const Pi = 3.14159
func SquareArea(side float64) float64
func RectangleArea(width, height float64) float64
func CircleArea(radius float64) float64 // Pi * r * r
```

## Rules
- Keep the helper `square` **unexported** — it's an implementation detail. Use it inside `SquareArea` and `CircleArea`.
- Read the compiler error first: `undefined: shapes.SquareArea` tells you exactly what's missing.
