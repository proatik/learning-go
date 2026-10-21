# Bank Ledger

Model a small bank: accounts stored by ID, money moved between them, and every successful movement recorded in a history. Money is an `int64` number of **cents** — never use floats for money.

```go
type Account struct {
	ID      string
	Owner   string
	Balance int64 // cents
}

// Tx is one money movement. From is "" for a deposit, To is "" for a withdrawal.
type Tx struct {
	From   string
	To     string
	Amount int64
}

type Ledger struct {
	Accounts map[string]*Account
	History  []Tx
}

var (
	ErrAccountExists = errors.New("bank: account already exists")
	ErrNoAccount     = errors.New("bank: no such account")
	ErrSameAccount   = errors.New("bank: transfer to the same account")
	ErrInvalidAmount = errors.New("bank: invalid amount")
	ErrInsufficient  = errors.New("bank: insufficient funds")
)

func NewLedger() *Ledger
func Open(l *Ledger, id, owner string, initial int64) error
func Deposit(l *Ledger, id string, amount int64) error
func Withdraw(l *Ledger, id string, amount int64) error
func Transfer(l *Ledger, from, to string, amount int64) error
func Balance(l *Ledger, id string) (int64, error)
func Statement(l *Ledger, id string) ([]Tx, error)
func Total(l *Ledger) int64
func FormatCents(c int64) string
```

## Example

```go
l := NewLedger()
Open(l, "a1", "Ada", 10_000)       // $100.00, recorded as Tx{"", "a1", 10000}
Open(l, "b2", "Bob", 0)            // nothing recorded
Transfer(l, "a1", "b2", 2_550)     // nil
Transfer(l, "b2", "a1", 9_999)     // ErrInsufficient — nothing changes
Balance(l, "a1")                   // 7450, nil
FormatCents(7450)                  // "74.50"
Statement(l, "b2")                 // [{a1 b2 2550}], nil
```

## Rules
- `NewLedger` returns a ledger with a non-nil `Accounts` map and an empty history. `l` is never nil in the other functions.
- **Open**: an empty `id` → any error; `initial < 0` → `ErrInvalidAmount`; `id` already used → `ErrAccountExists`. If `initial > 0`, record a deposit `Tx{"", id, initial}`.
- **Deposit / Withdraw**: check in this order — unknown account → `ErrNoAccount`; `amount <= 0` → `ErrInvalidAmount`; (Withdraw) `amount > balance` → `ErrInsufficient`. Record `Tx{"", id, amount}` / `Tx{id, "", amount}`.
- **Transfer**: check in this order — either account unknown → `ErrNoAccount`; `from == to` → `ErrSameAccount`; `amount <= 0` → `ErrInvalidAmount`; `amount > from's balance` → `ErrInsufficient`. Record `Tx{from, to, amount}`.
- A failed operation changes **nothing**: no balance changes and no history entry.
- **Balance**: unknown → `0, ErrNoAccount`.
- **Statement**: every history entry where `From` or `To` is `id`, in order. Return a **new slice** — changing it must not affect `l.History`. Unknown → `nil, ErrNoAccount`.
- **Total**: the sum of all balances.
- **FormatCents**: `1234` → `"12.34"`, `5` → `"0.05"`, `0` → `"0.00"`, `-150` → `"-1.50"`. Must work for every `int64`, including `math.MinInt64` (`"-92233720368547758.08"`).
- Return the sentinel errors (wrapping with `%w` is fine); tests use `errors.Is`.
