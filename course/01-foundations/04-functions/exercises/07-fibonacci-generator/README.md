# Fibonacci Generator

The Fibonacci sequence starts `0, 1` and each next number is the sum of the previous two: `0, 1, 1, 2, 3, 5, 8, 13, ...`

```go
func Fibonacci() func() int
func Nth(n int) (int, error)
```

- `Fibonacci` returns a **generator**: a closure that returns the next Fibonacci number each time it is called.
- `Nth` returns F(n), the n-th number counting from `F(0) = 0`, by **using a generator** from `Fibonacci`.

| Code | Result |
|---|---|
| `f := Fibonacci()`; 7 calls | `0, 1, 1, 2, 3, 5, 8` |
| `Nth(0)` | `0, nil` |
| `Nth(10)` | `55, nil` |
| `Nth(90)` | `2880067194370816120, nil` |
| `Nth(-1)` | `0, <error>` |

## Rules
- The generator keeps two captured variables, `a` and `b`, and advances them with a tuple assignment: `a, b = b, a+b`.
- `Nth(90)` must return instantly: the naive recursive `fib(n-1) + fib(n-2)` would take years.
- Negative `n` is an error.
