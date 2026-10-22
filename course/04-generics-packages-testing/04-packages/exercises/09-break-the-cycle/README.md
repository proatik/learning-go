# Break the Cycle

`bank` imports `audit` (to log activity) and `audit` imports `bank` (to read the account ID). The starter doesn't build:

```text
package learngo/course/…/09-break-the-cycle
	imports learngo/course/…/09-break-the-cycle/audit from wire.go
	imports learngo/course/…/09-break-the-cycle/bank from audit.go
	imports learngo/course/…/09-break-the-cycle/audit from bank.go: import cycle not allowed
```

Break the cycle by moving what both packages share into the small package `event`, so the imports become `bank → event ← audit`.

```go
// package event  (event/event.go): imports nothing from this exercise
type Event struct {
	Account string
	Kind    string // "deposit", "withdraw" or "rejected"
	Amount  int
}

type Recorder interface {
	Record(e Event)
}

// package bank  (bank/bank.go): imports event, NOT audit
func New(id string, rec event.Recorder) *Account // rec may be nil: nothing is recorded
func (a *Account) ID() string
func (a *Account) Balance() int
func (a *Account) Deposit(n int) error
func (a *Account) Withdraw(n int) error

// package audit  (audit/audit.go): imports event, NOT bank
type Log struct{ /* unexported fields */ } // the zero value is an empty log; *Log is an event.Recorder

func (l *Log) Record(e event.Event)
func (l *Log) Lines() []string
func (l *Log) Total(kind string) int

// package banking  (wire.go): the only place that knows both
func Open(id string) (*bank.Account, *audit.Log)
```

## Examples

```go
acc, log := banking.Open("acc-1")
acc.Deposit(100)   // nil
acc.Withdraw(30)   // nil
acc.Withdraw(500)  // ErrFunds
acc.Deposit(0)     // ErrAmount (not recorded)
acc.Balance()      // 70
log.Lines()        // ["acc-1 deposit 100", "acc-1 withdraw 30", "acc-1 rejected 500"]
log.Total("withdraw") // 30
```

## Rules
- `Deposit` / `Withdraw`:
  - `n <= 0` → `ErrAmount`, and nothing is recorded.
  - A withdrawal larger than the balance → `ErrFunds`, recorded as `Kind: "rejected"`.
  - Otherwise, update the balance and record `"deposit"` or `"withdraw"`.
- `Record` appends the line `"<account> <kind> <amount>"`. `Lines` returns a copy. `Total(kind)` sums the amounts of that kind.
- The behaviour is the same as in the starter: only the **dependencies** change.
- A test checks the imports: `bank` must not import `audit`, and `audit` must not import `bank`.
