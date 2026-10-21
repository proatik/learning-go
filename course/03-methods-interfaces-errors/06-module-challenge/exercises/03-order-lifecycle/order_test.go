package orders_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	orders "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/03-order-lifecycle"
)

func newOrder(t *testing.T, id string) *orders.Order {
	t.Helper()
	o := orders.NewOrder(id)
	if o == nil {
		t.Fatalf("NewOrder(%q) = nil, want an order", id)
	}
	return o
}

func TestStateString(t *testing.T) {
	tests := []struct {
		s    orders.State
		want string
	}{
		{orders.Created, "created"},
		{orders.Paid, "paid"},
		{orders.Shipped, "shipped"},
		{orders.Delivered, "delivered"},
		{orders.Cancelled, "cancelled"},
		{orders.State(5), "State(5)"},
		{orders.State(-1), "State(-1)"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", int(tt.s), got, tt.want)
		}
		if got := fmt.Sprint(tt.s); got != tt.want {
			t.Errorf("fmt.Sprint(State(%d)) = %q, want %q", int(tt.s), got, tt.want)
		}
	}
}

func TestCanTransition(t *testing.T) {
	all := []orders.State{orders.Created, orders.Paid, orders.Shipped, orders.Delivered, orders.Cancelled, orders.State(9)}
	allowed := map[[2]orders.State]bool{
		{orders.Created, orders.Paid}:      true,
		{orders.Created, orders.Cancelled}: true,
		{orders.Paid, orders.Shipped}:      true,
		{orders.Paid, orders.Cancelled}:    true,
		{orders.Shipped, orders.Delivered}: true,
	}
	for _, from := range all {
		for _, to := range all {
			want := allowed[[2]orders.State{from, to}]
			if got := orders.CanTransition(from, to); got != want {
				t.Errorf("CanTransition(%v, %v) = %v, want %v", from, to, got, want)
			}
		}
	}
}

func TestHappyPath(t *testing.T) {
	o := newOrder(t, "A1")
	if o.ID != "A1" || o.State() != orders.Created {
		t.Fatalf("NewOrder(%q): ID, State() = %q, %v, want %q, %v", "A1", o.ID, o.State(), "A1", orders.Created)
	}
	steps := []struct {
		name string
		do   func() error
		want orders.State
	}{
		{"Pay", o.Pay, orders.Paid}, // method values bound to o
		{"Ship", o.Ship, orders.Shipped},
		{"Deliver", o.Deliver, orders.Delivered},
	}
	for _, s := range steps {
		if err := s.do(); err != nil {
			t.Fatalf("%s() error = %v, want nil", s.name, err)
		}
		if got := o.State(); got != s.want {
			t.Fatalf("after %s(): State() = %v, want %v", s.name, got, s.want)
		}
	}
	want := []orders.State{orders.Created, orders.Paid, orders.Shipped, orders.Delivered}
	if got := o.History(); !slices.Equal(got, want) {
		t.Errorf("History() = %v, want %v", got, want)
	}
	if got := o.String(); got != "A1: created -> paid -> shipped -> delivered" {
		t.Errorf("String() = %q, want %q", got, "A1: created -> paid -> shipped -> delivered")
	}
}

func TestInvalidTransitions(t *testing.T) {
	tests := []struct {
		setup []orders.State // applied first, must succeed
		to    orders.State
	}{
		{nil, orders.Shipped},
		{nil, orders.Delivered},
		{nil, orders.Created}, // staying in the same state is not a transition
		{[]orders.State{orders.Paid}, orders.Paid},
		{[]orders.State{orders.Paid}, orders.Delivered},
		{[]orders.State{orders.Paid, orders.Shipped}, orders.Cancelled}, // too late to cancel
		{[]orders.State{orders.Paid, orders.Shipped, orders.Delivered}, orders.Cancelled},
		{[]orders.State{orders.Cancelled}, orders.Paid},
		{[]orders.State{orders.Paid, orders.Cancelled}, orders.Shipped},
		{nil, orders.State(42)},
	}
	for _, tt := range tests {
		o := newOrder(t, "X")
		for _, s := range tt.setup {
			if err := o.Transition(s); err != nil {
				t.Fatalf("Transition(%v) error = %v, want nil", s, err)
			}
		}
		from := o.State()
		before := o.History()
		err := o.Transition(tt.to)
		if err == nil {
			t.Errorf("Transition(%v) from %v error = nil, want a *TransitionError", tt.to, from)
			continue
		}
		var te *orders.TransitionError
		if !errors.As(err, &te) {
			t.Errorf("Transition(%v) from %v error = %v (%T), want a *TransitionError", tt.to, from, err, err)
			continue
		}
		if te.From != from || te.To != tt.to {
			t.Errorf("Transition(%v) error = %+v, want From=%v To=%v", tt.to, *te, from, tt.to)
		}
		if !errors.Is(err, orders.ErrInvalidTransition) {
			t.Errorf("errors.Is(%v, ErrInvalidTransition) = false, want true", err)
		}
		if want := fmt.Sprintf("orders: cannot go from %v to %v", from, tt.to); err.Error() != want {
			t.Errorf("Transition(%v) error = %q, want %q", tt.to, err.Error(), want)
		}
		if o.State() != from || !slices.Equal(o.History(), before) {
			t.Errorf("after failed Transition(%v): State, History = %v, %v, want unchanged %v, %v", tt.to, o.State(), o.History(), from, before)
		}
	}
}

func TestShortcutMethods(t *testing.T) {
	o := newOrder(t, "B2")
	if err := o.Ship(); !errors.Is(err, orders.ErrInvalidTransition) {
		t.Errorf("Ship() on a created order error = %v, want ErrInvalidTransition", err)
	}
	if err := o.Cancel(); err != nil {
		t.Fatalf("Cancel() error = %v, want nil", err)
	}
	if err := o.Pay(); !errors.Is(err, orders.ErrInvalidTransition) {
		t.Errorf("Pay() on a cancelled order error = %v, want ErrInvalidTransition", err)
	}
	if err := o.Deliver(); !errors.Is(err, orders.ErrInvalidTransition) {
		t.Errorf("Deliver() on a cancelled order error = %v, want ErrInvalidTransition", err)
	}
	if got := fmt.Sprint(o); got != "B2: created -> cancelled" {
		t.Errorf("fmt.Sprint(order) = %q, want %q", got, "B2: created -> cancelled")
	}
	// A different sentinel must not match.
	if err := o.Pay(); errors.Is(err, errors.New("orders: invalid transition")) {
		t.Errorf("errors.Is matched an unrelated error with the same text")
	}
}

func TestHistoryIsACopy(t *testing.T) {
	o := newOrder(t, "C3")
	if err := o.Pay(); err != nil {
		t.Fatalf("Pay() error = %v, want nil", err)
	}
	h := o.History()
	if len(h) != 2 {
		t.Fatalf("History() = %v, want [created paid]", h)
	}
	h[0] = orders.Delivered
	if got := o.History(); got[0] != orders.Created {
		t.Errorf("changing the History() result changed the order: %v", got)
	}
}

func TestReplay(t *testing.T) {
	o, err := orders.Replay("R1", orders.Paid, orders.Shipped)
	if err != nil || o == nil {
		t.Fatalf("Replay(R1, paid, shipped) = %v, %v, want an order, nil", o, err)
	}
	if got := o.String(); got != "R1: created -> paid -> shipped" {
		t.Errorf("Replay(R1, paid, shipped).String() = %q, want %q", got, "R1: created -> paid -> shipped")
	}

	o, err = orders.Replay("R2", orders.Paid, orders.Delivered, orders.Shipped)
	if o == nil {
		t.Fatalf("Replay(R2, paid, delivered, shipped) order = nil, want the order so far")
	}
	if got := o.String(); got != "R2: created -> paid" {
		t.Errorf("Replay(R2, ...).String() = %q, want %q (stop at the failing step)", got, "R2: created -> paid")
	}
	if want := "step 1: orders: cannot go from paid to delivered"; err == nil || err.Error() != want {
		t.Errorf("Replay(R2, ...) error = %v, want %q", err, want)
	}
	te, ok := errors.AsType[*orders.TransitionError](err)
	if !ok || te.From != orders.Paid || te.To != orders.Delivered {
		t.Errorf("errors.AsType[*TransitionError](%v) = %v, %v, want {paid delivered}, true", err, te, ok)
	}

	o, err = orders.Replay("R3")
	if err != nil || o == nil || o.State() != orders.Created {
		t.Errorf("Replay(R3) = %v, %v, want a created order, nil", o, err)
	}
}
