package orders

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

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

var stateNames = [...]string{"created", "paid", "shipped", "delivered", "cancelled"}

// String returns the lowercase state name, or "State(N)" for unknown values.
func (s State) String() string {
	if s < 0 || int(s) >= len(stateNames) {
		return fmt.Sprintf("State(%d)", int(s))
	}
	return stateNames[s]
}

// ErrInvalidTransition matches every *TransitionError via errors.Is.
var ErrInvalidTransition = errors.New("orders: invalid transition")

// TransitionError reports a move the lifecycle does not allow.
type TransitionError struct {
	From, To State
}

// Error returns "orders: cannot go from FROM to TO".
func (e *TransitionError) Error() string {
	return fmt.Sprintf("orders: cannot go from %v to %v", e.From, e.To)
}

// Is makes errors.Is(err, ErrInvalidTransition) true for a TransitionError.
func (e *TransitionError) Is(target error) bool {
	return target == ErrInvalidTransition
}

// allowed lists, for each state, the states it may move to.
var allowed = map[State][]State{
	Created: {Paid, Cancelled},
	Paid:    {Shipped, Cancelled},
	Shipped: {Delivered},
}

// CanTransition reports whether the lifecycle allows from → to.
func CanTransition(from, to State) bool {
	return slices.Contains(allowed[from], to)
}

// Order is one customer order. Create it with NewOrder.
type Order struct {
	ID      string
	state   State
	history []State
}

// NewOrder returns an order in the Created state.
func NewOrder(id string) *Order {
	return &Order{ID: id, state: Created, history: []State{Created}}
}

// State returns the current state.
func (o *Order) State() State { return o.state }

// History returns every state the order has been in, oldest first. The
// result is a copy.
func (o *Order) History() []State { return slices.Clone(o.history) }

// Transition moves the order to the given state, or returns a
// *TransitionError and leaves the order unchanged.
func (o *Order) Transition(to State) error {
	if !CanTransition(o.state, to) {
		return &TransitionError{From: o.state, To: to}
	}
	o.state = to
	o.history = append(o.history, to)
	return nil
}

// Pay moves the order to Paid.
func (o *Order) Pay() error { return o.Transition(Paid) }

// Ship moves the order to Shipped.
func (o *Order) Ship() error { return o.Transition(Shipped) }

// Deliver moves the order to Delivered.
func (o *Order) Deliver() error { return o.Transition(Delivered) }

// Cancel moves the order to Cancelled.
func (o *Order) Cancel() error { return o.Transition(Cancelled) }

// Replay creates an order and applies steps in order. It stops at the first
// failing step and returns the order so far together with the error wrapped
// as "step N: ...".
func Replay(id string, steps ...State) (*Order, error) {
	o := NewOrder(id)
	for i, to := range steps {
		if err := o.Transition(to); err != nil {
			return o, fmt.Errorf("step %d: %w", i, err)
		}
	}
	return o, nil
}

// String returns "ID: s1 -> s2 -> ...".
func (o *Order) String() string {
	names := make([]string, len(o.history))
	for i, s := range o.history {
		names[i] = s.String()
	}
	return o.ID + ": " + strings.Join(names, " -> ")
}
