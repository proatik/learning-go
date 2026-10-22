# Three-Stage Pipeline

Build the classic `Gen → Square → Sum` pipeline: each stage runs in its own goroutine and talks to the next through a channel.

```go
func Gen(nums ...int) <-chan int
func Square(in <-chan int) <-chan int
func Sum(in <-chan int) int
```

```text
Gen(1, 2, 3) ──1,2,3──▶ Square ──1,4,9──▶ Sum ──▶ 14
```

## Examples

| Code | Result |
|---|---|
| `Gen(1, 2, 3)` | yields `1`, `2`, `3`, then is closed |
| `Gen()` | closed straight away, yields nothing |
| `Square(Gen(-4, 0, 10))` | yields `16`, `0`, `100`, then is closed |
| `Sum(Square(Gen(1, 2, 3)))` | `14` |
| `Sum(Square(Gen()))` | `0` |

## Rules
- `Gen` and `Square` return **at once** and do their sending from a goroutine they start.
- Each stage **owns** its output channel: it creates it, is the only sender, and closes it when done (`Gen` after its last value, `Square` after `in` is closed and drained).
- Values keep their order through every stage.
- `Square` **streams**: it emits `n*n` as soon as it receives `n`; it never waits for the whole input.
- `Sum` receives until `in` is closed, then returns the total (`0` for no values).
