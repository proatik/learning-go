# Bank Account

A bank account that many goroutines can use at once. The tricky part is `Withdraw`: the "is there enough money?" check and the subtraction must happen as **one** step.

```go
var ErrInvalidAmount = errors.New("bank: amount must be positive")
var ErrInsufficientFunds = errors.New("bank: insufficient funds")

type Account struct {
	mu      sync.Mutex // guards balance
	balance int
}

func (a *Account) Deposit(amount int) error
func (a *Account) Withdraw(amount int) error
func (a *Account) Balance() int
```

## Examples

| Code | Result |
|---|---|
| `var a Account; a.Deposit(100)` | `nil`; `Balance()` = `100` |
| `a.Withdraw(30)` | `nil`; `Balance()` = `70` |
| `a.Withdraw(71)` | `ErrInsufficientFunds`; `Balance()` still `70` |
| `a.Deposit(0)` / `a.Withdraw(-5)` | `ErrInvalidAmount`, nothing changes |
| balance 500, 1000 goroutines each `Withdraw(1)` | exactly 500 succeed, 500 get `ErrInsufficientFunds`, `Balance()` = `0` |

## Rules
- `amount <= 0` → an error matching `ErrInvalidAmount` (`errors.Is`), for both `Deposit` and `Withdraw`.
- `Withdraw` of more than the balance → an error matching `ErrInsufficientFunds`; the balance is unchanged. Withdrawing the whole balance is fine.
- You may wrap the sentinels with context (`fmt.Errorf("…: %w", …)`).
- The balance must **never** go below 0, even with many concurrent withdrawals: check and subtract inside the **same** `Lock` / `Unlock`.
- `Balance` locks too.
