# Insufficient Funds

Report a failed withdrawal with an error that carries the numbers, so the caller can show "you are short by $50".

```go
var ErrInvalidAmount = errors.New("invalid amount")

type InsufficientFundsError struct {
	Need int // the amount requested
	Have int // the balance at that moment
}

func (e *InsufficientFundsError) Error() string // "insufficient funds: need 150, have 100"
func (e *InsufficientFundsError) Shortfall() int // Need - Have

type Account struct {
	Balance int
}

func (a *Account) Withdraw(amount int) error
func (a *Account) Transfer(to *Account, amount int) error
```

## Examples

| Start balance | Call | Error | Balance after |
|---|---|---|---|
| 100 | `Withdraw(30)` | `nil` | 70 |
| 100 | `Withdraw(100)` | `nil` | 0 |
| 100 | `Withdraw(150)` | `insufficient funds: need 150, have 100` | 100 |
| 100 | `Withdraw(0)` | `withdraw 0: invalid amount` | 100 |
| 100 | `Withdraw(-5)` | `withdraw -5: invalid amount` | 100 |
| a=100, b=0 | `a.Transfer(&b, 150)` | `transfer: insufficient funds: need 150, have 100` | a=100, b=0 |

## Rules
- `Withdraw`:
  - `amount <= 0`: return `fmt.Errorf("withdraw %d: %w", amount, ErrInvalidAmount)`.
  - `amount > Balance`: return `&InsufficientFundsError{Need: amount, Have: Balance}`.
  - Otherwise subtract `amount` and return `nil`.
  - On any error the balance must **not** change.
- `Transfer`: withdraw from `a`. If that fails, wrap the error as `transfer: <err>` with `%w` and leave both balances untouched. Otherwise add `amount` to `to`.
