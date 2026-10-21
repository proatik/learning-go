# Bank Transfers

Module 2's ledger, redesigned with types: accounts with encapsulated state and methods, a `Money` type that prints itself, a typed error carrying details, and a `Transfer` that **rolls back** with `defer` when its second step fails.

```go
var (
	ErrInvalidAmount = errors.New("bank: invalid amount")
	ErrAccountFrozen = errors.New("bank: account frozen")
	ErrNoAccount     = errors.New("bank: no such account")
	ErrAccountExists = errors.New("bank: account already exists")
	ErrSameAccount   = errors.New("bank: same account")
	ErrInvalidID     = errors.New("bank: invalid account id")
)

type Money int64 // cents
func (m Money) String() string

type InsufficientFundsError struct {
	ID      string
	Balance Money
	Amount  Money
}
func (e *InsufficientFundsError) Error() string

type Entry struct {
	Desc    string
	Amount  Money // positive = in, negative = out
	Balance Money // balance after the entry
}
func (e Entry) String() string

type Account struct { id, owner string; balance Money; frozen bool; entries []Entry }
func (a *Account) ID() string
func (a *Account) Owner() string
func (a *Account) Balance() Money
func (a *Account) Frozen() bool
func (a *Account) Freeze()
func (a *Account) Unfreeze()
func (a *Account) Entries() []Entry
func (a *Account) Deposit(amt Money) error
func (a *Account) Withdraw(amt Money) error
func (a *Account) String() string

type Bank struct{ accounts map[string]*Account }
func NewBank() *Bank
func (b *Bank) Open(id, owner string) (*Account, error)
func (b *Bank) Account(id string) (*Account, error)
func (b *Bank) Transfer(from, to string, amt Money) error
```

## Examples

| Code | Result |
|---|---|
| `Money(1234)`, `Money(-5)` | `12.34`, `-0.05` |
| `a.Withdraw(10_001)` with balance `100.00` | `*InsufficientFundsError` — `"bank: a1 has 100.00, needs 100.01"` |
| `a.Freeze(); a.Deposit(1)` | `"a1: bank: account frozen"` |
| `b.Transfer("b2", "a1", 1)` | `"transfer b2 -> a1: bank: b2 has 0.00, needs 0.01"` |
| `b2.Freeze(); b.Transfer("a1", "b2", 3000)` | error, and `a1` is **unchanged** |

Statement (`a.String()`) after a deposit and a transfer:

```text
a1 (Ada): 74.50
  deposit +100.00 = 100.00
  transfer to b2 -25.50 = 74.50
```

## Rules
- **Money.String**: `"%d.%02d"` of the magnitude with a leading `-` for negatives; must work for `math.MinInt64` (`"-92233720368547758.08"`).
- **Entry.String**: `"DESC +AMOUNT = BALANCE"` — a `+` for positive amounts; negative amounts already print their `-`.
- **Deposit / Withdraw** check in this order:
  1. `amt <= 0` → bare `ErrInvalidAmount`
  2. frozen → `fmt.Errorf("%s: %w", id, ErrAccountFrozen)`
  3. (Withdraw) `amt > balance` → `&InsufficientFundsError{id, balance, amt}`
- On success they update the balance and append an `Entry` with desc `"deposit"` / `"withdraw"`. A failure changes **nothing**.
- **Entries** returns a copy. **Account.String** is the header `"ID (OWNER): BALANCE"` followed by `"\n  " + entry` for each entry.
- **Open**: empty `id` → `ErrInvalidID`; duplicate → an error wrapping `ErrAccountExists`. **Account**: unknown ID → an error wrapping `ErrNoAccount`.
- **Transfer** checks, in order: both accounts exist (`ErrNoAccount`), `from != to` (`ErrSameAccount`), then debits `from` with desc `"transfer to TO"` and credits `to` with desc `"transfer from FROM"`, using the same checks as `Withdraw` / `Deposit`.
  - If the credit fails, a **deferred** rollback restores `from`'s balance **and** removes its debit entry.
  - Every error returned by Transfer is wrapped as `fmt.Errorf("transfer %s -> %s: %w", from, to, err)` — do that in a deferred function too, using the **named** result `err`.
