# Bug Hunt: Stack and Queue

The code is written — your job is the **tests**. Write table-driven tests with **subtests** in `containers_test.go` that pass on the real generic `Stack` and `Queue` and fail for each of 9 planted bugs.

```go
// Provided in check_test.go (read it in the Tests tab).
func NewStack[T any]() *Stack[T]
func (s *Stack[T]) Push(v T)
func (s *Stack[T]) Pop() (T, bool)  // top item; zero value and false when empty
func (s *Stack[T]) Peek() (T, bool) // top item, not removed
func (s *Stack[T]) Len() int
func (s *Stack[T]) All() iter.Seq[T] // top to bottom

func NewQueue[T any]() *Queue[T]     // ring buffer: room for 4, doubles when full
func (q *Queue[T]) Push(v T)
func (q *Queue[T]) Pop() (T, bool)   // front item; zero value and false when empty
func (q *Queue[T]) Len() int
func (q *Queue[T]) All() iter.Seq[T] // front to back
```

## Examples

| Your tests on… | Must |
|---|---|
| the correct code | pass |
| `stack Pop takes the bottom item` | fail |
| `queue grow scrambles a wrapped buffer` | fail |
| … every other bug listed in `check_test.go` | fail |

## Rules
- Write at least one `Test…` function, and use **at least one** `t.Run` subtest. The checker runs your tests once per bug.
- Drive the containers through their **methods** only — don't read `items`, `buf`, `head` or `n`.
- Cover: empty containers, order of `Pop`, `Peek` leaving the item in place, `Len` after pops, `All` order, `break` inside `range All()`, and a queue that **wraps around** (pop some, then push past the end) and **grows** while wrapped.
- Check both results of `Pop`: the value **and** `ok`. An empty container must return the **zero value**, even after items have come and gone.
- A failing run lists each surviving bug as `Bug "…" = not caught`.
