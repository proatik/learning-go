package contracts_test

import (
	"testing"

	contracts "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/03-satisfy"
)

func TestTagIsNamer(t *testing.T) {
	if _, ok := any(contracts.Tag{}).(contracts.Namer); !ok {
		t.Fatalf("Tag{} does not implement Namer, want it to (check the method name)")
	}
	for _, text := range []string{"go", "", "héllo wörld"} {
		n := any(contracts.Tag{Text: text}).(contracts.Namer)
		if got := n.Name(); got != text {
			t.Errorf("Tag{%q}.Name() = %q, want %q", text, got, text)
		}
	}
}

func TestTagIsLabeler(t *testing.T) {
	if _, ok := any(contracts.Tag{}).(contracts.Labeler); !ok {
		t.Fatalf("Tag{} does not implement Labeler, want it to (use a value receiver)")
	}
	tests := []struct{ text, want string }{
		{"go", "#go"},
		{"", "#"},
		{"día", "#día"},
	}
	for _, tt := range tests {
		l := any(contracts.Tag{Text: tt.text}).(contracts.Labeler)
		if got := l.Label(); got != tt.want {
			t.Errorf("Tag{%q}.Label() = %q, want %q", tt.text, got, tt.want)
		}
	}
}

func TestTallyIsResettableCounter(t *testing.T) {
	c, ok := any(&contracts.Tally{}).(contracts.ResettableCounter)
	if !ok {
		t.Fatalf("*Tally does not implement ResettableCounter, want it to")
	}
	if got := c.Count(); got != 0 {
		t.Errorf("new Tally: Count() = %d, want 0", got)
	}
	for i := 1; i <= 3; i++ {
		c.Inc()
		if got := c.Count(); got != i {
			t.Errorf("after %d Inc: Count() = %d, want %d", i, got, i)
		}
	}
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Errorf("after Reset: Count() = %d, want 0", got)
	}
	c.Inc()
	if got := c.Count(); got != 1 {
		t.Errorf("after Reset and Inc: Count() = %d, want 1", got)
	}
}

func TestTallyCountersAreIndependent(t *testing.T) {
	a, okA := any(&contracts.Tally{}).(contracts.Counter)
	b, okB := any(&contracts.Tally{}).(contracts.Counter)
	if !okA || !okB {
		t.Fatalf("*Tally does not implement Counter, want it to")
	}
	a.Inc()
	a.Inc()
	b.Inc()
	if a.Count() != 2 || b.Count() != 1 {
		t.Errorf("Count() = %d and %d, want 2 and 1", a.Count(), b.Count())
	}
}

func TestCelsiusIsLabeler(t *testing.T) {
	if _, ok := any(contracts.Celsius(0)).(contracts.Labeler); !ok {
		t.Fatalf("Celsius does not implement Labeler, want it to")
	}
	tests := []struct {
		c    contracts.Celsius
		want string
	}{
		{21.5, "21.5°C"},
		{-40, "-40.0°C"},
		{0, "0.0°C"},
		{36.6, "36.6°C"},
		{100, "100.0°C"},
	}
	for _, tt := range tests {
		l := any(tt.c).(contracts.Labeler)
		if got := l.Label(); got != tt.want {
			t.Errorf("Celsius(%v).Label() = %q, want %q", float64(tt.c), got, tt.want)
		}
	}
}
