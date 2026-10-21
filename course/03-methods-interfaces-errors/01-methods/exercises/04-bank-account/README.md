# Bank Account

An encapsulated type: unexported fields, a constructor, and methods that protect the invariant **balance ≥ 0**.

```go
type Account struct {
	owner   string
	balance int64 // cents
}

func NewAccount(owner string, initial int64) *Account
func (a *Account) Owner() string
func (a *Account) Balance() int64
func (a *Account) SetOwner(name string) bool
func (a *Account) Deposit(amount int64) bool
func (a *Account) Withdraw(amount int64) bool
func (a *Account) Transfer(to *Account, amount int64) bool
```

## Examples

| Code | Result | Balance afterwards |
|---|---|---|
| `a := NewAccount("ann", 100)` | | `100` |
| `NewAccount("bob", -5)` | | `0` |
| `a.Deposit(50)` | `true` | `150` |
| `a.Deposit(0)` | `false` | unchanged |
| `a.Withdraw(30)` | `true` | `120` |
| `a.Withdraw(500)` | `false` | unchanged |
| `a.Transfer(b, 20)` | `true` | `a` −20, `b` +20 |
| `a.SetOwner("")` | `false` | owner unchanged |

## Rules
- `NewAccount` treats a negative `initial` as `0`.
- `Deposit` and `Withdraw` reject `amount <= 0`; `Withdraw` also rejects `amount > Balance()`. On `false`, nothing changes.
- `Transfer` returns `false` (changing nothing) when `to` is `nil`, `to` is the same account as `a`, or the withdrawal would fail.
- `SetOwner` rejects the empty string.
- The tests live in another package (`bank_test`), so they can only use the exported API, just like real callers.
