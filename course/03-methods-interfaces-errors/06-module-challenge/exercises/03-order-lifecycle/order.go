package orders

import "errors"

// State is a step in an order's lifecycle.
type State int

// The lifecycle states.
const (
	Created State = iota
	Paid
	Shipped
	Delivered
	Cancelled
)

// String returns the lowercase state name, or "State(N)" for unknown values.
func (s State) String() string {
	// TODO
	return ""
}

// ErrInvalidTransition matches every *TransitionError via errors.Is.
var ErrInvalidTransition = errors.New("orders: invalid transition")

// TransitionError reports a move the lifecycle does not allow.
type TransitionError struct {
	From, To State
}

// Error returns "orders: cannot go from FROM to TO".
func (e *TransitionError) Error() string {
	// TODO
	return ""
}

// Is makes errors.Is(err, ErrInvalidTransition) true for a TransitionError.
func (e *TransitionError) Is(target error) bool {
	// TODO
	return false
}

// CanTransition reports whether the lifecycle allows from → to.
func CanTransition(from, to State) bool {
	// TODO
	return false
}

// Order is one customer order. Create it with NewOrder.
type Order struct {
	ID      string
	state   State
	history []State
}

// NewOrder returns an order in the Created state.
func NewOrder(id string) *Order {
	// TODO
	return nil
}

// State returns the current state.
func (o *Order) State() State {
	// TODO
	return 0
}

// History returns every state the order has been in, oldest first. The
// result is a copy.
func (o *Order) History() []State {
	// TODO
	return nil
}

// Transition moves the order to the given state, or returns a
// *TransitionError and leaves the order unchanged.
func (o *Order) Transition(to State) error {
	// TODO
	return nil
}

// Pay moves the order to Paid.
func (o *Order) Pay() error {
	// TODO
	return nil
}

// Ship moves the order to Shipped.
func (o *Order) Ship() error {
	// TODO
	return nil
}

// Deliver moves the order to Delivered.
func (o *Order) Deliver() error {
	// TODO
	return nil
}

// Cancel moves the order to Cancelled.
func (o *Order) Cancel() error {
	// TODO
	return nil
}

// Replay creates an order and applies steps in order. It stops at the first
// failing step and returns the order so far together with the error wrapped
// as "step N: ...".
func Replay(id string, steps ...State) (*Order, error) {
	// TODO
	return nil, nil
}

// String returns "ID: s1 -> s2 -> ...".
func (o *Order) String() string {
	// TODO
	return ""
}
