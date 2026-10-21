# Transactions

Write the classic `WithTx` helper: begin a transaction, run a callback, **commit** on success and **roll back** on error *or panic* — exactly once, guaranteed by `defer`.

```go
type Tx interface {
	Exec(stmt string) error
	Commit() error
	Rollback() error
}

type DB interface {
	Begin() (Tx, error)
}

func WithTx(db DB, fn func(tx Tx) error) (err error)
func ExecAll(db DB, stmts []string) error
```

| Situation in `WithTx` | Calls after `Begin` | Result |
|---|---|---|
| `fn` returns `nil`, `Commit` ok | `Commit` | `nil` |
| `fn` returns `nil`, `Commit` fails with `c` | `Commit` | `commit: <c>` |
| `fn` returns `e` | `Rollback` | `e` itself |
| `fn` returns `e`, `Rollback` fails with `r` | `Rollback` | `errors.Join(e, rollback: <r>)` |
| `fn` **panics** with `v` | `Rollback` | panics again with the same `v` |
| `Begin` fails with `b` | none (`fn` not called) | `begin: <b>` |

## Rules
- **`WithTx`**:
  - After a successful `Begin`, exactly **one** of `Commit` / `Rollback` is called (never both, never twice). Don't roll back after a failed `Commit`.
  - Callback error: return it unchanged (`err == e`), unless `Rollback` also fails: then `errors.Join(e, fmt.Errorf("rollback: %w", rerr))`.
  - Panic in `fn`: call `Rollback` (ignore its error), then **re-panic** with the original value. Never swallow the panic.
  - Wrap with `%w`: `begin: %w`, `commit: %w`.
- **`ExecAll`**:
  - Uses `WithTx` to `Exec` every statement in order inside one transaction.
  - Stop at the first failing statement and return `fmt.Errorf("exec %q: %w", stmt, err)` (so the transaction is rolled back).
  - No statements → begin and commit an empty transaction.
- Recommended: a `defer` registered right after `Begin` that rolls back unless the commit path was reached.
