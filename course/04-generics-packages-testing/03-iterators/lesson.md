# Iterators

Since Go 1.23, `for range` also works over **functions**. An iterator is a function that **pushes** values, one at a time, into a callback named `yield`. With generics, you can write iterators for your own types and reusable adapters (`Filter`, `Map`, `Take`, …) that work with every sequence.

| You want to… | Use |
|---|---|
| Loop over your own collection | An `All()` method returning `iter.Seq[V]` |
| Transform a sequence lazily | An adapter: `func(iter.Seq[In]) iter.Seq[Out]` |
| Turn a sequence into a slice / map | `slices.Collect`, `slices.Sorted`, `maps.Collect` |
| Read two sequences side by side | `iter.Pull` |

## 1. Range over a function

A `for range` loop accepts a function of one of these three shapes:

| Function type | Loop | Values per step |
|---|---|---|
| `func(yield func() bool)` | `for range f` | 0 |
| `func(yield func(V) bool)` | `for v := range f` | 1 |
| `func(yield func(K, V) bool)` | `for k, v := range f` | 2 |

```go
// Countdown yields n, n-1, …, 1.
func Countdown(n int) func(func(int) bool) {
	return func(yield func(int) bool) {
		for i := n; i >= 1; i-- {
			if !yield(i) { // the loop body ran a break: stop now
				return
			}
		}
	}
}

for v := range Countdown(3) {
	fmt.Print(v, " ") // 3 2 1
}
```

- `Countdown(3)` **builds** the iterator; nothing runs yet. The `for range` loop calls it.
- Each `yield(i)` runs the loop body once with `v = i`.
- `yield` returns `true` to ask for more, `false` when the body has stopped (`break`, `return`, `goto` out of the loop, or a panic).
- Range over functions needs `go 1.23` or later in `go.mod`.

## 2. `iter.Seq` and `iter.Seq2`

Package `iter` names the two common shapes, so signatures stay readable:

```go
package iter

type Seq[V any] func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)
```

```go
func Countdown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := n; i >= 1; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
```

- `Seq[V]` is a sequence of single values (`for v := range seq`).
- `Seq2[K, V]` is a sequence of pairs, usually index/value or key/value (`for k, v := range seq`).
- They are plain generic function types: a function literal of the right shape **is** a `Seq`, with no conversion.
- Functions that accept or return iterators should use these types, so they plug into `slices`, `maps` and other adapters.

## 3. The yield contract

| Rule | Why |
|---|---|
| Call `yield` once per value, in order | Each call runs the loop body once |
| **Stop as soon as `yield` returns `false`** | The loop body has finished: `break`, `return` or a panic |
| Never call `yield` again after it returned `false` | The synthesized `yield` panics |
| Return when the sequence ends | The loop ends when your function returns |

What happens if the iterator ignores the result of `yield`:

```go
func Bad(yield func(int) bool) {
	for i := range 3 {
		yield(i) // ❌ result ignored
	}
}

for v := range Bad {
	fmt.Println(v)
	break
}
```

```text
0
panic: runtime error: range function continued iteration after function for loop body returned false
```

- The rule is easy: **every** `yield` call is `if !yield(v) { return }`.
- An adapter that calls an inner sequence must also stop the inner loop (see Section 7).
- A recursive iterator must pass the "stop" signal back up the call stack (exercise 9).

## 4. `for range` vs calling the iterator yourself

An iterator is just a function, so you can call it with your own `yield`:

```go
Countdown(5)(func(v int) bool {
	fmt.Print(v, " ")
	return v > 3 // keep going while v > 3
})
// 5 4 3
```

With `for range`, the compiler writes the `yield` function for you. Roughly:

```go
// for v := range seq { body }
seq(func(v int) bool {
	body          // break → return false, continue → return true
	return true
})
```

| | `for range seq` | `seq(yield)` |
|---|---|---|
| `break`, `continue`, `return` | Work as in any loop | You return `false` / `true` yourself |
| `return` from the enclosing function | ✅ e.g. `return true` inside the loop | ❌ only returns from your callback |
| Misbehaving iterator | Panics with a clear message | Silently keeps calling you |
| Typical use | Almost always | Tests and adapters that count calls |

```go
func find(seq iter.Seq[int], target int) bool {
	for v := range seq {
		if v == target {
			return true // leaves find; the iterator's yield returns false
		}
	}
	return false
}
```

## 5. Iterator methods on your own types

Give a collection an `All` method; name other walks after the sequence or order they produce.

```go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// All returns an iterator over the items from top to bottom.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.items) - 1; i >= 0; i-- {
			if !yield(s.items[i]) {
				return
			}
		}
	}
}

var s Stack[string]
s.Push("a")
s.Push("b")
for v := range s.All() {
	fmt.Print(v) // ba
}
```

| Convention (from the `iter` docs) | Example |
|---|---|
| `All` walks every element | `func (s *Set[V]) All() iter.Seq[V]` |
| Name the sequence when there are several | `Keys()`, `Values()`, `Cities()` |
| Name the order when there are several | `All()` head→tail, `Backward()` tail→head |
| Extra configuration → arguments | `func (m *Map[K, V]) Scan(min, max K) iter.Seq2[K, V]` |

- Calling the iterator again walks the sequence again from the start.
- **Single-use iterators** read from a stream that can't rewind (a reader, a network connection). Say so in the doc comment: "It returns a single-use iterator."
- Don't change a collection while you iterate over it, unless your type documents what happens.

## 6. Laziness and infinite sequences

An iterator produces a value only when the consumer asks for it. So a sequence can be **infinite**: the consumer decides when to stop.

```go
// Naturals yields 0, 1, 2, … forever.
func Naturals() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

for n := range Naturals() {
	if n > 3 {
		break // without a break this loop never ends
	}
	fmt.Print(n, " ") // 0 1 2 3
}
```

- `slices.Collect(Naturals())` never returns: it asks for every value.
- Put a limiting adapter like `Take` (next section) in front of consumers that read everything.

## 7. Generic adapters

An adapter takes a sequence and returns a new one. It ranges over the input **inside** the returned function, so nothing runs until someone ranges over the result.

```go
// Filter yields the values of seq for which keep returns true.
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return // our consumer stopped: the return also stops seq
			}
		}
	}
}

// Map yields f(v) for every v in seq.
func Map[In, Out any](seq iter.Seq[In], f func(In) Out) iter.Seq[Out] {
	return func(yield func(Out) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}
```

The values flow **one at a time** through the whole pipeline:

```go
traced := func(yield func(int) bool) {
	for i := 1; i <= 4; i++ {
		fmt.Println("source", i)
		if !yield(i) {
			return
		}
	}
}
squares := Map(Filter(traced, func(n int) bool { return n%2 == 0 }),
	func(n int) int { return n * n })
fmt.Println("built")
for v := range squares {
	fmt.Println("got", v)
}
```

```text
built
source 1
source 2
got 4
source 3
source 4
got 16
```

- `return` inside `for v := range seq` makes the inner `yield` return `false`, so `seq` stops too. That is how a `break` at the end of a pipeline stops the source.
- `Take` limits a sequence; it must not ask for a value it won't use:

```go
// Take yields the first n values of seq.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		i := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			i++
			if i == n {
				return // stop BEFORE seq produces value n+1
			}
		}
	}
}

squares := Map(Naturals(), func(n int) int { return n * n })
fmt.Println(slices.Collect(Take(squares, 5))) // [0 1 4 9 16]
```

### Adapters with two values: `Enumerate`, `Concat`, `Zip`

```go
// Enumerate pairs each value with its index: 0, 1, 2, …
func Enumerate[V any](seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Concat yields every value of each sequence in turn.
func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return // return, not break: break would only leave the inner loop
				}
			}
		}
	}
}

for i, s := range Enumerate(slices.Values([]string{"a", "b"})) {
	fmt.Println(i, s) // 0 a, then 1 b
}
fmt.Println(slices.Collect(Concat(slices.Values([]int{1, 2}), slices.Values([]int{3})))) // [1 2 3]
```

- `Zip` (read two sequences in step) can't be written with two nested `for range` loops. It needs a pull iterator (Section 11).

## 8. Consuming sequences

| Function | Signature | Result |
|---|---|---|
| `slices.Collect` | `(seq iter.Seq[E]) []E` | All values, in order |
| `slices.AppendSeq` | `(s Slice, seq iter.Seq[E]) Slice` | `s` with the values appended |
| `slices.Sorted` | `(seq iter.Seq[E]) []E` (`E cmp.Ordered`) | Collected and sorted |
| `slices.SortedFunc` | `(seq iter.Seq[E], cmp func(E, E) int) []E` | Sorted by `cmp` |
| `maps.Collect` | `(seq iter.Seq2[K, V]) map[K]V` | A new map |
| `maps.Insert` | `(m Map, seq iter.Seq2[K, V])` | Adds the pairs to `m` |

```go
m := map[string]int{"b": 2, "a": 1, "c": 3}
fmt.Println(slices.Sorted(maps.Keys(m))) // [a b c]  ← the idiom for "map keys in order"

fmt.Println(maps.Collect(Enumerate(slices.Values([]string{"go", "rust"})))) // map[0:go 1:rust]
fmt.Println(slices.AppendSeq([]string{"x"}, strings.FieldsSeq("y z")))     // [x y z]
```

- Every consumer here reads the **whole** sequence. Never pass it an infinite one.

## 9. Iterators in the standard library

| Function | Yields | Notes |
|---|---|---|
| `slices.Values(s)` | `iter.Seq[E]`: elements | |
| `slices.All(s)` | `iter.Seq2[int, E]`: index, element | |
| `slices.Backward(s)` | `iter.Seq2[int, E]` | Descending indices |
| `slices.Chunk(s, n)` | `iter.Seq[Slice]`: sub-slices of up to `n` | Panics if `n < 1`; chunks share `s`'s memory |
| `maps.Keys(m)`, `maps.Values(m)` | `iter.Seq[K]`, `iter.Seq[V]` | **Random order**, like `range m` |
| `maps.All(m)` | `iter.Seq2[K, V]` | Random order |
| `strings.Lines(s)` | `iter.Seq[string]` | Lines **keep** their `"\n"` |
| `strings.SplitSeq(s, sep)` | `iter.Seq[string]` | Same pieces as `strings.Split` |
| `strings.FieldsSeq(s)` | `iter.Seq[string]` | Same pieces as `strings.Fields` |
| `bytes.Lines`, `bytes.SplitSeq`, `bytes.FieldsSeq` | `iter.Seq[[]byte]` | The `[]byte` versions |

```go
for i, v := range slices.Backward([]string{"a", "b", "c"}) {
	fmt.Print(i, v, " ") // 2c 1b 0a
}
for c := range slices.Chunk([]int{1, 2, 3, 4, 5}, 2) {
	fmt.Print(c, " ") // [1 2] [3 4] [5]
}
for l := range strings.Lines("a\nb\nc") {
	fmt.Printf("%q ", l) // "a\n" "b\n" "c"
}
for p := range strings.SplitSeq("a,b,,c", ",") {
	fmt.Printf("%q ", p) // "a" "b" "" "c"
}
for f := range strings.FieldsSeq("  go  is\tfun ") {
	fmt.Printf("%q ", f) // "go" "is" "fun"
}
```

- `strings.SplitSeq` and `FieldsSeq` don't allocate the result slice: the right choice when you only loop once.

## 10. Cleanup inside an iterator

The iterator function runs until it returns, so a `defer` inside it runs when the loop finishes, **breaks** or panics.

```go
func Rows() iter.Seq[int] {
	return func(yield func(int) bool) {
		fmt.Println("open")
		defer fmt.Println("close") // runs on every exit path
		for i := range 5 {
			if !yield(i) {
				return
			}
		}
	}
}

for v := range Rows() {
	if v == 1 {
		break
	}
	fmt.Println(v)
}
```

```text
open
0
close
```

- Open the resource **inside** the returned function, not in the constructor: then building an iterator costs nothing, and every walk opens and closes its own resource.
- This only works because the iterator returns promptly when `yield` returns `false`.

## 11. Pull iterators: `iter.Pull`

`for range` **pushes** values into your loop body. Sometimes you need to **pull** the next value when you are ready, e.g. to walk two sequences side by side.

```go
func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func())
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func())
```

```go
next, stop := iter.Pull(Countdown(3))
defer stop()
for {
	v, ok := next()
	if !ok {
		break // the sequence is over
	}
	fmt.Print(v, " ") // 3 2 1
}
```

| Function | Behaviour |
|---|---|
| `next()` | The next value and `true`; at the end the zero value and `false`, forever after |
| `stop()` | Ends the iteration: the iterator's `yield` returns `false`, so its deferred cleanup runs |

- The iterator does not start until the first `next()` call.
- **Always call `stop`** if you might quit before `next` reports `false`. Write `defer stop()` right after `iter.Pull`. Calling `stop` more than once, or after the end, is fine.
- Forgetting `stop` leaves the iterator paused forever: its `defer` never runs and its memory is never freed.

`Zip` pushes one sequence and pulls the other:

```go
// Zip pairs the values of as and bs; it stops when either one ends.
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextB, stop := iter.Pull(bs)
		defer stop()
		for a := range as {
			b, ok := nextB()
			if !ok || !yield(a, b) {
				return
			}
		}
	}
}

for s, n := range Zip(slices.Values([]string{"x", "y", "z"}), Naturals()) {
	fmt.Print(s, n, " ") // x0 y1 z2
}
```

- Prefer `for range` whenever it works: it is simpler and faster than `Pull`.

## 12. Iterators that can fail

`yield` can't return an error to the iterator, and `Seq` has no error result. Two common designs:

| Design | Use when | Example |
|---|---|---|
| `iter.Seq2[V, error]` | Each step can fail | Yield `(zero, err)` once, then return |
| An `Err() error` method | The source fails as a whole, like `bufio.Scanner` | Loop, then check `x.Err()` |

```go
// Records upper-cases each line; an empty line is an error.
func Records(lines []string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for i, l := range lines {
			if l == "" {
				yield("", fmt.Errorf("line %d: empty record", i+1))
				return // nothing comes after an error
			}
			if !yield(strings.ToUpper(l), nil) {
				return
			}
		}
	}
}

for rec, err := range Records([]string{"a", "b", "", "c"}) {
	if err != nil {
		fmt.Println("error:", err) // error: line 3: empty record
		break
	}
	fmt.Println(rec) // A, then B
}
```

- Calling `yield` and then returning without checking its result is fine: you never call it again.

## 13. Return a sequence or a slice?

| | Return `[]T` | Return `iter.Seq[T]` |
|---|---|---|
| Memory | All values at once | One value at a time |
| Early exit | All work is done up front | Work stops at `break` |
| Infinite / streaming data | ❌ | ✅ |
| `len`, indexing, walking twice cheaply | ✅ | ❌ (collect it first) |
| Simplicity for small, finite data | ✅ | More moving parts |

- Return a slice when the caller needs all the data or random access.
- Return a sequence for large, lazy, infinite or streaming data, and for walking your own types. The caller can always `slices.Collect` it.

## Common mistakes

| Mistake | Fix |
|---|---|
| `yield(v)` with the result ignored | `if !yield(v) { return }` every time |
| `break` inside nested loops of an iterator after `yield` returned `false` | `return`: `break` leaves only the inner loop |
| Recursive iterator ignores the helper's "stopped" result | Return a `bool` from the helper and stop when it is `false` |
| Adapter that reads one value too many (`Take` checks the count before yielding) | Stop right after the `n`-th `yield` |
| Doing the work in the constructor instead of inside the returned function | Keep the constructor cheap: return a closure that does the work |
| Returning a `nil` `iter.Seq` | Ranging over it panics; return `func(yield func(V) bool) {}` |
| `slices.Collect` / `slices.Sorted` on an infinite sequence | Limit it first with `Take` |
| `iter.Pull` without `stop` | `next, stop := iter.Pull(seq)` then `defer stop()` |
| Expecting `maps.Keys` in sorted order | `slices.Sorted(maps.Keys(m))` |
| Expecting `strings.Lines` to strip `"\n"` | It keeps it: trim, or use `bufio.Scanner` |
| Opening a resource in the constructor of a sequence | Open it inside the iterator and `defer` the close |

## Key takeaways

- An iterator is a function `func(yield func(V) bool)`; `for range` calls it and turns the loop body into `yield`.
- Use `iter.Seq[V]` and `iter.Seq2[K, V]` in signatures.
- **The contract:** stop as soon as `yield` returns `false`. Ignoring it makes `for range` panic.
- Collections get an `All()` method (plus `Backward()`, `Keys()`, …); constructors stay cheap and do the work inside the closure.
- Adapters (`Filter`, `Map`, `Take`, `Enumerate`, `Concat`) are lazy: values flow one at a time, and a `break` at the end stops the source.
- Infinite sequences are fine as long as something stops them.
- `slices.Collect`, `slices.Sorted`, `maps.Collect` consume sequences; `slices`, `maps`, `strings` and `bytes` produce them.
- `defer` inside an iterator runs on break, end and panic.
- `iter.Pull` turns a push iterator into `next`/`stop`; always `defer stop()`.
- Return a slice for small data that needs indexing; return a sequence for large, lazy or streaming data.
