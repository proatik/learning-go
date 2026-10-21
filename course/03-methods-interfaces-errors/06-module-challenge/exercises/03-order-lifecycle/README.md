# Order Lifecycle

Model an order as a **state machine**: methods enforce the allowed transitions, and every illegal move returns a typed error that also matches a sentinel.

```text
created ──► paid ──► shipped ──► delivered
   │          │
   └──────────┴──► cancelled
```

```go
type State int

const (
	Created State = iota
	Paid
	Shipped
	Delivered
	Cancelled
)

func (s State) String() string

var ErrInvalidTransition = errors.New("orders: invalid transition")

type TransitionError struct {
	From, To State
}

func (e *TransitionError) Error() string
func (e *TransitionError) Is(target error) bool

type Order struct {
	ID      string
	state   State
	history []State
}

func CanTransition(from, to State) bool
func NewOrder(id string) *Order
func (o *Order) State() State
func (o *Order) History() []State
func (o *Order) Transition(to State) error
func (o *Order) Pay() error
func (o *Order) Ship() error
func (o *Order) Deliver() error
func (o *Order) Cancel() error
func (o *Order) String() string
func Replay(id string, steps ...State) (*Order, error)
```

## Examples

| Code | Result |
|---|---|
| `Shipped.String()`, `State(7).String()` | `"shipped"`, `"State(7)"` |
| `o := NewOrder("A1"); o.Pay(); o.Ship()` | `nil`, `nil` |
| `o.Cancel()` | `*TransitionError{Shipped, Cancelled}` |
| `err.Error()` | `"orders: cannot go from shipped to cancelled"` |
| `errors.Is(err, ErrInvalidTransition)` | `true` |
| `o.String()` | `"A1: created -> paid -> shipped"` |
| `Replay("R", Paid, Delivered)` | order `R: created -> paid`, error `"step 1: orders: cannot go from paid to delivered"` |

## Rules
- **String** returns `created`, `paid`, `shipped`, `delivered`, `cancelled`; any other value is `"State(N)"`.
- The **only** allowed transitions are the arrows above. `Delivered` and `Cancelled` are final; moving to the **same** state is not allowed; unknown states are never allowed.
- **NewOrder** starts in `Created` with history `[Created]`.
- **Transition** on success sets the state and appends it to the history. On failure it returns `&TransitionError{From: current, To: to}` and changes **nothing**.
- On success return a literal `nil` — not a nil `*TransitionError` stored in an `error`.
- **Pay / Ship / Deliver / Cancel** are shortcuts for `Transition(Paid)` etc.
- **TransitionError.Is** returns `true` exactly when `target == ErrInvalidTransition`.
- **History** returns a copy — changing it must not change the order.
- **Order.String**: `"ID: "` followed by the history names joined with `" -> "`.
- **Replay** creates an order and applies `steps` in order. At the first failing step `i` it stops and returns the order **so far** and `fmt.Errorf("step %d: %w", i, err)`; otherwise the order and `nil`.
