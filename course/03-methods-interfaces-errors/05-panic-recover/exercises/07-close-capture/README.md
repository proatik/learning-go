# Close Error Capture

Write lines to a store that must **always be closed**, without losing the error from `Close`, and without piling up open stores in a loop.

```go
type Store interface {
	Write(line string) error
	Close() error
}

func Save(open func() (Store, error), lines []string) (err error)
func SaveBatches(open func(i int) (Store, error), batches [][]string) error
```

| Situation in `Save` | Returned error | `Close` called? |
|---|---|---|
| `open` fails with `e` | `open: <e>` | no |
| all writes succeed, `Close` returns `nil` | `nil` | once |
| all writes succeed, `Close` fails with `c` | `close: <c>` | once |
| `Write("b")` fails with `e` | `write "b": <e>` (Close error, if any, is dropped) | once |
| `Write` **panics** | the panic keeps going (don't recover) | once, before the panic leaves `Save` |

## Rules
- **`Save`**:
  - Call `open()` once. If it fails, return `fmt.Errorf("open: %w", err)`.
  - Write `lines` in order; stop at the first failing write and return `fmt.Errorf("write %q: %w", line, err)`.
  - `Close` exactly once whenever `open` succeeded, using `defer`.
  - Report `fmt.Errorf("close: %w", cerr)` **only if** no earlier error happened (the first error wins).
  - All errors wrap with `%w` (tests use `errors.Is`).
- **`SaveBatches`**:
  - For each batch `i`, save it with `Save` and `open(i)`.
  - Each store must be **closed before the next one is opened** (no `defer` piling up in the loop: let `Save` do it).
  - Stop at the first failing batch and return `fmt.Errorf("batch %d: %w", i, err)`.
  - Return `nil` when all batches succeed (also for no batches).
