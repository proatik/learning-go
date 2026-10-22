# Lazy Pipeline Library

Write a small, reusable package of **lazy iterator adapters** (`pipe/`), then build an order report in the root package **using only those adapters**.

```go
// package pipe — learngo/course/04-generics-packages-testing/06-module-challenge/exercises/06-lazy-pipeline/pipe
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U]
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T]
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T]
func Batch[T any](seq iter.Seq[T], size int) iter.Seq[[]T]
func Fold[T, A any](seq iter.Seq[T], init A, f func(A, T) A) A

// package report — the exercise root
type Order struct {
	ID       int
	Customer string
	Cents    int64 // never negative
	Refunded bool
}

func Revenue(orders iter.Seq[Order]) int64
func LargeOrderIDs(orders iter.Seq[Order], minCents int64, n int) []int
func Line(o Order) string
func Pages(orders iter.Seq[Order], size int) iter.Seq[[]string]
```

## Examples

| Code | Result |
|---|---|
| `Take(Filter(Map(1..100, square), >10), 2)` | `[16 25]` — only 5 values pulled and squared |
| `Batch(1..7, 3)` | `[1 2 3]`, `[4 5 6]`, `[7]` |
| `Fold([1 2 3], "x", concat)` | `"x123"` |
| `Revenue(orders)` | sum of `Cents` of orders that are **not** refunded |
| `LargeOrderIDs(orders, 1000, 2)` | IDs of the first 2 paid orders with `Cents >= 1000` |
| `Line(Order{42, "bo", 7, false})` | `"#42 bo $0.07"` |
| `Line(Order{9, "cy", 100000, true})` | `"#9 cy $1000.00 (refunded)"` |
| `Pages(orders, 2)` | pages of 2 `Line`s, the last page may be shorter |

## Rules
- **Lazy:** building a pipeline does no work. Values are pulled one at a time only while the consumer keeps ranging.
- Every adapter honours the yield contract and stops pulling from its source as soon as its consumer stops.
- **Take(seq, n)** pulls **at most** `n` values (none if `n <= 0`) — it must not pull an `n+1`th value just to discover it is done.
- **Batch** yields a **fresh** slice each time (collected batches must not change later); `size < 1` yields nothing.
- **report** imports `pipe` and builds each function from it: `Revenue` = `Fold` over `Filter`; `LargeOrderIDs` = `Filter` → `Map` to the ID → `Take` → `slices.Collect`; `Pages` = `Batch` over `Map(orders, Line)`.
- **LargeOrderIDs** skips refunded orders, keeps input order, stops reading after `n` matches, and returns an empty, non-nil slice when there are none.
- **Line** is `"#<ID> <Customer> $<dollars>.<2-digit cents>"`, plus `" (refunded)"` if refunded.
