# Mutexes & Atomics

## 1. Data races

Channels (lesson 2) share memory **by communicating**. Sometimes you just want several goroutines to update one value, e.g. a counter or a cache. Doing that without synchronization is a **data race**.

```go
n := 0
var wg sync.WaitGroup
for range 1000 {
	wg.Go(func() { n++ }) // 1000 goroutines write n with no synchronization
}
wg.Wait()
fmt.Println(n) // 1000 on one run, 897 on the next: the result is not predictable
```

- A **data race** is two goroutines accessing the same variable at the same time, where **at least one writes**, with nothing ordering the two accesses.
- `n++` is three steps: read `n`, add 1, write `n`. Two goroutines can both read 41 and both write 42, so one increment is **lost**.
- A racy program has **no guaranteed behaviour**. Lost updates are the mild case. You can also get stale reads, or half-written values for multi-word types like strings, slices and interfaces.
- Maps have a built-in check. Concurrent writes crash the whole program with `fatal error: concurrent map writes`, and `recover` can't catch it.
- Reads are only safe together. Any number of goroutines may **read** a value no one is writing.

**The race detector.** `go test -race` and `go run -race` instrument every memory access. When two goroutines touch the same memory without synchronization, they print `WARNING: DATA RACE` with both stack traces.

| | |
|---|---|
| How it works | finds races that **actually happen** during that run, so run it on tests that exercise concurrency |
| Cost | typically 2–20× slower and 5–10× more memory, so use it in CI and local runs, not in production builds |
| Requirement | **cgo**, so a C toolchain must be installed |
| On this platform | cgo is off: `go test -race` prints `go: -race requires cgo; enable cgo by setting CGO_ENABLED=1` |

Instead, the tests here **hammer** your code (e.g. 50 goroutines × 1000 operations) and check exact totals. Unsynchronized code almost always loses updates and fails. On a machine with a C compiler, also run your own projects with `-race`; it is the standard check in Go CI.

## 2. `sync.Mutex`

A **mutex** (mutual exclusion lock) lets only one goroutine at a time run the code between `Lock` and `Unlock`. That code is the **critical section**.

```go
// Counter is safe for concurrent use.
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

var c Counter
var wg sync.WaitGroup
for range 1000 {
	wg.Go(c.Inc)
}
wg.Wait()
fmt.Println(c.Value()) // 1000, every time
```

| Method | What it does |
|---|---|
| `mu.Lock()` | waits until the mutex is free, then takes it |
| `mu.Unlock()` | releases it; one waiting goroutine can then take it. Unlocking an unlocked mutex is a fatal error |
| `mu.TryLock() bool` | takes it only if it's free right now. Rarely the right tool (the tests here use it to peek at your locking) |

- **The zero value is an unlocked mutex.** No constructor is needed, and `var c Counter` is ready to use.
- Write **`defer mu.Unlock()` on the line after `Lock()`.** Then every `return` and even a panic releases the lock. A forgotten unlock freezes every later `Lock` forever.
- **Reads lock too.** `Value` must lock, or it races with `Inc`.
- Unlocking also **publishes** the changes. Whatever a goroutine wrote before `Unlock` is visible to the next goroutine after its `Lock`.
- **Go mutexes are not reentrant.** Calling `Lock` on a mutex you already hold deadlocks. So a locked method must not call another method that locks the same mutex:

```go
func (c *Counter) IncTwice() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Inc() // ❌ Inc calls c.mu.Lock() again: this goroutine waits for itself forever
}
```

  If you need shared logic, put it in an unexported helper that **assumes** the lock is held, e.g. `incLocked()`, and call it from both methods.

## 3. Guarding struct fields

```go
// Inventory is safe for concurrent use.
type Inventory struct {
	mu    sync.Mutex     // guards items
	items map[string]int // name → quantity
}

// Snapshot returns a copy of the items.
func (inv *Inventory) Snapshot() map[string]int {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	return maps.Clone(inv.items) // never hand out the guarded map itself
}

inv := &Inventory{items: map[string]int{"apple": 3}}
snap := inv.Snapshot()
snap["apple"] = 0
fmt.Println(inv.Snapshot()["apple"]) // 3
```

The conventions:

- **Put the mutex right above the fields it guards**, with a `// guards x, y` comment. Fields above the mutex, or in another struct, are not protected by it.
- **Keep guarded fields unexported.** If callers can reach `inv.items` directly, no lock can protect it.
- **Use pointer receivers** on every method. A value receiver gets a **copy** of the struct and its mutex, so it locks nothing (§9).
- **Don't let references escape.** Returning the guarded map or slice lets the caller read or write it later **without the lock**. Return a copy (`maps.Clone`, `slices.Clone`) or a single value instead.
- Keep critical sections **short**. Lock, touch the data, unlock. Don't do I/O, sleep, send on a channel or call unknown callbacks while holding a lock.

## 4. Check-then-act: one decision, one critical section

Locking each access separately is not enough when a decision depends on what you read:

```go
// ❌ two critical sections: another goroutine can withdraw between them
if a.Balance() >= n { // lock, read, unlock
	a.subtract(n) // lock, write, unlock
}
```

Two goroutines can both see a balance of 1, both decide "enough", and both subtract, leaving `-1`. The **check and the act** must happen under one lock:

```go
var ErrInsufficient = errors.New("insufficient funds")

type Account struct {
	mu      sync.Mutex // guards balance
	balance int
}

func (a *Account) Withdraw(n int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < n {
		return ErrInsufficient
	}
	a.balance -= n
	return nil
}

// 150 goroutines each Withdraw(1) from a balance of 100
fmt.Println(successes, a.balance) // 100 0
```

The same rule applies to "add if absent", "create on first use" and "increment and return the new value". A method's result must come from the **same** critical section as its change.

## 5. `sync.RWMutex`

An `RWMutex` has two modes: many **readers** at once, or one **writer** alone.

```go
// Cache is safe for concurrent use.
type Cache struct {
	mu    sync.RWMutex // guards items
	items map[string]string
}

func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock() // shared: other readers may hold it too
	defer c.mu.RUnlock()
	v, ok := c.items[k]
	return v, ok
}

func (c *Cache) Set(k, v string) {
	c.mu.Lock() // exclusive: waits for every reader and writer to leave
	defer c.mu.Unlock()
	c.items[k] = v
}
```

| Held by | `RLock` | `Lock` |
|---|---|---|
| nobody | ✅ gets it | ✅ gets it |
| one or more readers | ✅ gets it (unless a writer is already waiting) | ⏳ waits |
| a writer | ⏳ waits | ⏳ waits |

- Pair each lock with its own unlock: `RLock` with `RUnlock`, `Lock` with `Unlock`.
- **A waiting writer blocks new readers**, so writers can't starve. For the same reason, taking `RLock` twice in one goroutine can deadlock.
- **You can't upgrade.** Calling `Lock` while holding `RLock` deadlocks. Release the read lock, take `Lock`, then **check again**: something may have changed in between.
- It pays off only for **read-heavy** data with non-trivial reads. An `RWMutex` costs more than a `Mutex`, so for short critical sections or balanced reads and writes a plain `Mutex` is often as fast. Start with `Mutex`, and switch when profiling shows contention.

## 6. Run once: `sync.Once`, `OnceFunc`, `OnceValue`, `OnceValues`

Lazy initialization ("load the config on first use") is check-then-act on a flag. The `sync` package has it ready-made:

```go
var once sync.Once
for range 3 {
	once.Do(func() { fmt.Println("init") }) // prints "init" once
}

calls := 0
port := sync.OnceValue(func() int { calls++; return 8080 })
fmt.Println(port(), port(), calls) // 8080 8080 1

load := sync.OnceValues(func() (Config, error) {
	return Config{}, errors.New("no config file")
})
_, err1 := load()
_, err2 := load()
fmt.Println(err1, err1 == err2) // no config file true

greet := sync.OnceFunc(func() { fmt.Println("hello once") })
greet()
greet() // prints nothing
```

| API | You get | Use for |
|---|---|---|
| `var once sync.Once` + `once.Do(f)` | a field or variable whose `Do` runs `f` the first time only | initialization inside a struct method |
| `sync.OnceFunc(f func()) func()` | a function that runs `f` once | "close once", "log the warning once" |
| `sync.OnceValue(f func() T) func() T` | a function returning `f`'s cached result | a lazily computed value |
| `sync.OnceValues(f func() (T1, T2)) func() (T1, T2)` | the same for two results | the `(value, error)` loader |

- **Concurrent callers wait.** While `f` runs, every other caller blocks until it finishes, then all of them see its result. That makes it safe for any number of goroutines.
- **Errors are cached, not retried.** `OnceValues` returns the same error forever. If you need a retry, `Once` is the wrong tool.
- **Panics:** if `f` panics, the `OnceFunc` / `OnceValue(s)` wrappers panic **with the same value on every call**. `once.Do` counts a panicking `f` as done, and later calls do nothing.
- Calling `once.Do` from **inside** its own `f` deadlocks.
- One `Once` per thing. Two `App` values need two `Once`s, so store it (or the `OnceValues` function) in the struct, not in a package variable.

## 7. `sync/atomic` typed values

For a **single** number, flag or pointer, the `sync/atomic` types are lock-free and faster than a mutex. Each method is one indivisible step.

```go
var hits atomic.Int64
for range 1000 {
	wg.Go(func() { hits.Add(1) })
}
wg.Wait()
fmt.Println(hits.Load()) // 1000

var n atomic.Int64
fmt.Println(n.Add(5), n.Swap(1), n.Load()) // 5 5 1: Add returns the new value, Swap the old one

var ready atomic.Bool
fmt.Println(ready.CompareAndSwap(false, true), ready.CompareAndSwap(false, true)) // true false

var cfg atomic.Pointer[Config] // the zero value holds nil
cfg.Store(&Config{Port: 80})
old := cfg.Swap(&Config{Port: 443}) // hot-swap the config for every reader
fmt.Println(old.Port, cfg.Load().Port) // 80 443
```

| Type | Methods |
|---|---|
| `atomic.Int32`, `Int64`, `Uint32`, `Uint64`, `Uintptr` | `Load`, `Store`, `Add`, `Swap`, `CompareAndSwap`, `And`, `Or` |
| `atomic.Bool` | `Load`, `Store`, `Swap`, `CompareAndSwap` |
| `atomic.Pointer[T]` | `Load`, `Store`, `Swap`, `CompareAndSwap` (on `*T`) |

- The zero value is 0 / `false` / `nil` and ready to use. Like a mutex, an atomic **must not be copied** after first use.
- **Use the return value.** `x := n.Add(1)` is *your* new value. A separate `n.Load()` afterwards may already include other goroutines' adds.
- **Atomics only cover one value.** Two atomics updated one after the other are not one atomic step together. A reader can see the first change without the second. When several fields must change together (a balance and a history, a map and its size), use a mutex.
- `atomic.Pointer[T]` is the way to swap a whole immutable value at once (config, routing table). Readers `Load()` and use the pointer without locking. Writers build a **new** value and `Store` it, and never modify the old one.
- Older code uses functions like `atomic.AddInt64(&x, 1)` on plain `int64` fields. The typed values are preferred: you can't accidentally read the field without an atomic.

## 8. Compare-and-swap loops

`CompareAndSwap(old, new)` stores `new` **only if** the value is still `old`, and reports whether it did. It is the building block for any "read, compute, write" update that `Add` can't express, e.g. "keep the maximum":

```go
// StoreMax sets m to v if v is larger, and reports whether it did.
func StoreMax(m *atomic.Int64, v int64) bool {
	for {
		old := m.Load()
		if v <= old {
			return false // nothing to do
		}
		if m.CompareAndSwap(old, v) {
			return true // nobody changed m between our Load and our CAS
		}
		// someone else stored first: read the new value and decide again
	}
}

var max atomic.Int64
for i := range 100 {
	wg.Go(func() { StoreMax(&max, int64(i)) })
}
wg.Wait()
fmt.Println(max.Load()) // 99
```

- `old := m.Load(); m.Store(compute(old))` is **check-then-act** again. Another goroutine can store in between, and your `Store` overwrites its value.
- The loop is **optimistic**. It assumes no one interferes, and retries when someone did. Under heavy contention it may spin a few times, but it never blocks.
- Keep the loop body tiny and side-effect free, because it may run more than once.

## 9. `sync.Map` vs a map with a mutex

`sync.Map` is a concurrent map built into the standard library.

```go
var m sync.Map // zero value is empty; keys and values are any
m.Store("go", 1)
v, ok := m.Load("go")
fmt.Println(v.(int)+1, ok) // 2 true: every value needs a type assertion

actual, loaded := m.LoadOrStore("go", 99) // atomic "get or insert"
fmt.Println(actual, loaded)               // 1 true

m.Range(func(k, v any) bool { // like a for loop over the map; return false to stop
	fmt.Println(k, v) // go 1
	return true
})
```

Other methods: `Delete`, `LoadAndDelete`, `Swap`, `CompareAndSwap`, `CompareAndDelete`, `Clear`.

| | `map[K]V` + `Mutex` / `RWMutex` | `sync.Map` |
|---|---|---|
| Type safety | ✅ `K`, `V` checked at compile time | ❌ `any` in, type assertion out |
| Invariants across several fields (map + counter + list) | ✅ one lock covers them all | ❌ only the map itself is protected |
| `len` | ✅ `len(m)` | ❌ none; count with `Range` |
| Best for | **most code** | (1) keys written once and read many times (a cache that only grows); (2) many goroutines working on **disjoint** sets of keys |

The `sync` docs say it plainly: *"Most code should use a plain Go map instead, with separate locking or coordination."* Reach for `sync.Map` only when profiling shows lock contention in one of those two patterns.

## 10. Never copy a mutex

A `Mutex`, `RWMutex`, `Once`, `WaitGroup`, `sync.Map` or any `atomic` type **must not be copied after first use**. The copy is a separate lock, so the two halves of your program stop excluding each other. Copies hide in value receivers, value parameters, assignments and `range` loops. `go vet`'s **copylocks** check finds them:

```go
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c Counter) Inc() { c.mu.Lock(); c.n++; c.mu.Unlock() } // value receiver
func show(c Counter)   { fmt.Println(c.n) }

d := c                       // copy
for _, c := range counters { // copy per element
	fmt.Println(c.n)
}
```

```text
Inc passes lock by value: main.Counter contains sync.Mutex
show passes lock by value: main.Counter contains sync.Mutex
assignment copies lock value to d: main.Counter contains sync.Mutex
range var c copies lock: main.Counter contains sync.Mutex
```

The fixes:

- **Pointer receivers** on every method of a type that holds a lock (`func (c *Counter) Inc()`).
- **Pass and return pointers**: `func NewCache() *Cache`, `func show(c *Counter)`.
- Range by index or over pointers: `for i := range counters { counters[i].Inc() }`, or use a `[]*Counter`.
- Run `go vet ./...`. It is part of every normal Go workflow (the editor runs it too). Note that `go test` runs only a few vet checks, and copylocks is **not** one of them.

## 11. Deadlocks and lock ordering

A **deadlock** happens when goroutines wait for each other in a cycle. The classic case is two locks taken in **opposite orders**:

```go
// transfer locks from, then to: the order depends on the caller.
func transfer(from, to *Account, n int) {
	from.mu.Lock()
	defer from.mu.Unlock()
	to.mu.Lock()
	defer to.mu.Unlock()
	from.balance -= n
	to.balance += n
}

wg.Go(func() { transfer(a, b, 10) }) // locks a, waits for b
wg.Go(func() { transfer(b, a, 10) }) // locks b, waits for a
wg.Wait()
```

```text
fatal error: all goroutines are asleep - deadlock!
```

It doesn't happen every run. Only when the timing lines up, which is what makes deadlocks dangerous. The runtime reports it only when **every** goroutine is stuck. In a server with other goroutines alive, the two requests simply hang forever.

The fixes, best first:

| Fix | How |
|---|---|
| **Hold one lock at a time** | finish with one object before locking the next, if the logic allows it |
| **A global lock order** | always lock in the same order, e.g. by ascending account ID. The cycle can't form: whoever waits holds nothing the other needs |
| **Never block while holding a lock** | no channel sends or receives, no `wg.Wait`, no network calls, no callbacks into unknown code inside a critical section |

```go
first, second := from, to
if second.ID < first.ID {
	first, second = second, first // smaller ID always first, whatever the direction
}
first.mu.Lock()
defer first.mu.Unlock()
second.mu.Lock()
defer second.mu.Unlock()
```

- Also guard **`from == to`**: locking the same mutex twice deadlocks on its own.
- To lock *many* objects (e.g. to sum all balances at one consistent moment), sort a copy by the same key and lock in that order.

## 12. Choosing a tool

| Situation | Use |
|---|---|
| Hand data or ownership to another goroutine; signal "done" | **channel** (lesson 2) |
| Several fields that must change together; maps; check-then-act | **`sync.Mutex`** |
| Read-heavy shared data with real read work | **`sync.RWMutex`** |
| One counter, flag or pointer | **`sync/atomic`** typed value |
| Initialize once, lazily | **`sync.OnceValue(s)`** / `OnceFunc` / `Once` |
| A cache that only grows, or disjoint keys per goroutine, under measured contention | `sync.Map` |

A common Go proverb is *"Don't communicate by sharing memory; share memory by communicating."* It's a guideline, not a ban. A mutex around a struct is often the simplest correct design, and the standard library uses mutexes everywhere.

## Common mistakes

| Mistake | Fix |
|---|---|
| Reading shared data without the lock ("it's only a read") | Readers lock too (`Lock` or `RLock`); a read racing a write is a data race |
| `Lock()` with no `Unlock()` on an early `return` | `defer mu.Unlock()` right after `Lock()` |
| Checking under one lock, acting under another | One critical section for the check **and** the change |
| A locked method calling another locked method | Mutexes aren't reentrant; use an unexported `…Locked` helper |
| Returning the guarded map / slice | Return a copy (`maps.Clone`, `slices.Clone`) |
| Value receiver on a type with a mutex | Pointer receivers; run `go vet` (copylocks) |
| `RLock` then `Lock` to "upgrade" | Release, `Lock`, re-check |
| `n.Add(1)` then `n.Load()` to get "my" value | Use `Add`'s return value |
| `Load` + `Store` for read-modify-write | `CompareAndSwap` loop (or a mutex) |
| Two atomics expected to change together | A mutex around both fields |
| Locking two objects in caller-dependent order | A global order (e.g. by ID); guard `a == b` |
| Blocking (channel op, I/O, `Wait`) while holding a lock | Copy what you need, unlock, then block |
| `sync.Map` as the default concurrent map | `map` + `Mutex`; `sync.Map` only for its two specific patterns |
| Package-level `sync.Once` for per-instance setup | Put the `Once` / `OnceValues` func in the struct |

## Key takeaways

- A **data race** is unsynchronized concurrent access with at least one write. The result is undefined, and concurrent map writes crash the program. `go test -race` finds races but needs cgo, so here the tests hammer the code and check exact totals.
- `sync.Mutex`: `Lock` + `defer Unlock`. The zero value is ready. Put it next to the fields it guards, keep those unexported, use pointer receivers, and don't leak references.
- Make each decision (check-then-act) inside **one** critical section. Go mutexes aren't reentrant.
- `sync.RWMutex` allows many readers or one writer. It has no upgrade, and is worth it only for read-heavy data.
- `sync.OnceFunc` / `OnceValue` / `OnceValues` / `Once` run initialization exactly once. Concurrent callers wait, and errors and panics are cached.
- `sync/atomic` typed values (`Int64`, `Bool`, `Pointer[T]`) suit single values. Use their return values, and `CompareAndSwap` loops for read-modify-write.
- Prefer `map` + mutex over `sync.Map` unless profiling says otherwise.
- Never copy a lock or an atomic. `go vet` copylocks catches it.
- Avoid deadlocks with a global lock order. Hold locks briefly, and never block while holding one.
