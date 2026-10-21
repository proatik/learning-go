# Unwind

Watch a panic **unwind** the stack: every deferred call runs, innermost first, before `recover` stops it.

```go
func Unwind(w io.Writer, depth int)
func Acquire(w io.Writer, names []string, failAt int)
```

`Unwind(w, 3)` calls nested levels 1, 2, 3; level 3 panics:

```text
enter 1
enter 2
enter 3
exit 3
exit 2
exit 1
recovered: bottom at 3
```

`Acquire(w, []string{"db", "cache", "queue"}, 1)` fails on index 1:

```text
acquire db
acquire cache
release cache
release db
error: failed at cache
```

`Acquire(w, []string{"db", "cache"}, -1)` never fails:

```text
acquire db
acquire cache
all done
release cache
release db
```

## Rules
- Every line ends with `"\n"`; write with `fmt.Fprintln` / `fmt.Fprintf` to `w`.
- **`Unwind`**:
  - `depth <= 0`: print only `nothing to unwind`.
  - Otherwise level `i` (from 1) prints `enter i`, **defers** printing `exit i`, then either panics with the string `bottom at <depth>` (when `i == depth`) or calls level `i+1`.
  - `Unwind` recovers and prints `recovered: <value>` **after** all `exit` lines.
- **`Acquire`**:
  - For each name `i`: print `acquire <name>`, **defer** printing `release <name>`, then if `i == failAt` panic with `failed at <name>`.
  - With no panic, print `all done` after the last acquire.
  - On a panic, print `error: <value>` as the very **last** line (after all releases).
  - Use `defer` inside the loop (fine here: the loop is short and the releases must all wait for the function to end).
- No panic may escape either function.
