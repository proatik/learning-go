# Deadlock-Free Transfers

Move money between two accounts that each have their own mutex. Locking `from` then `to` deadlocks as soon as one goroutine sends A→B while another sends B→A. Fix it with a **lock order**: always lock the account with the smaller ID first.

```go
var (
	ErrInvalidAmount     = errors.New("ledger: amount must be positive")
	ErrSameAccount       = errors.New("ledger: cannot transfer to the same account")
	ErrInsufficientFunds = errors.New("ledger: insufficient funds")
)

type Account struct {
	ID      int
	mu      sync.Mutex // guards balance
	balance int
}

func NewAccount(id, balance int) *Account
func (a *Account) Balance() int

func Transfer(from, to *Account, amount int) error
func Total(accounts ...*Account) int
```

## Examples

| Code (`a` = ID 1 with 100, `b` = ID 2 with 50) | Result |
|---|---|
| `Transfer(a, b, 30)` | `nil`; balances `70`, `80` |
| `Transfer(b, a, 81)` | `ErrInsufficientFunds`; nothing changes |
| `Transfer(a, b, 0)` | `ErrInvalidAmount` |
| `Transfer(a, a, 10)` | `ErrSameAccount` |
| `Total(a, b)` | `150` |
| `Total()` | `0` |
| 40 goroutines × 1000 transfers in both directions | finishes (no deadlock); `Total` is always the starting sum |

## Rules
- Errors match the sentinels with `errors.Is` (wrapping is fine). Checks, in order: `amount <= 0` → `ErrInvalidAmount`; `from == to` → `ErrSameAccount` (locking one mutex twice would deadlock); then, under the locks, `from`'s balance `< amount` → `ErrInsufficientFunds`. On any error nothing changes.
- `Transfer` holds **both** locks while it checks and moves the money, so no one can observe the money in neither account or in both.
- Lock order: the account with the **smaller `ID`** first, whatever the direction. The tests lock one account themselves and check that `Transfer` waits for it **without** holding the other.
- `Total` locks **every** account, in ascending `ID` order (sort a **copy** of the slice; don't reorder the caller's), sums the balances, then unlocks them all. Summing `Balance()` one account at a time could count a transfer's money twice or not at all.
- Accounts passed to `Total` are distinct; IDs are unique.
