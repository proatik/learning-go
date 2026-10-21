package calc

import "testing"

func TestMethods(t *testing.T) {
	var c Calc
	if got := c.Result(); got != 0 {
		t.Fatalf("zero Calc Result() = %d, want 0", got)
	}
	c.Add(7)
	c.Sub(2)
	c.Mul(3)
	if got := c.Result(); got != 15 {
		t.Errorf("Add(7), Sub(2), Mul(3): Result() = %d, want 15", got)
	}
}

func TestBound(t *testing.T) {
	var c Calc
	ops := Bound(&c)
	if len(ops) != 3 {
		t.Fatalf("Bound(&c) returned %d funcs, want 3", len(ops))
	}
	for i, op := range ops {
		if op == nil {
			t.Fatalf("Bound(&c)[%d] = nil, want a method value", i)
		}
	}
	ops[0](5) // Add
	if got := c.Result(); got != 5 {
		t.Errorf("Bound(&c)[0](5): Result() = %d, want 5", got)
	}
	ops[2](3) // Mul
	if got := c.Result(); got != 15 {
		t.Errorf("then Bound(&c)[2](3): Result() = %d, want 15", got)
	}
	ops[1](20) // Sub
	if got := c.Result(); got != -5 {
		t.Errorf("then Bound(&c)[1](20): Result() = %d, want -5", got)
	}

	var other Calc
	otherOps := Bound(&other)
	if len(otherOps) == 3 && otherOps[0] != nil {
		otherOps[0](100)
	}
	if c.Result() != -5 || other.Result() != 100 {
		t.Errorf("Bound values must be bound to their own Calc: got (%d, %d), want (-5, 100)", c.Result(), other.Result())
	}
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name  string
		start int
		arg   int
		want  int
	}{
		{"add", 10, 4, 14},
		{"sub", 10, 4, 6},
		{"mul", 10, 4, 40},
		{"add", 0, -3, -3},
	}
	for _, tt := range tests {
		f := Lookup(tt.name)
		if f == nil {
			t.Errorf("Lookup(%q) = nil, want a method expression", tt.name)
			continue
		}
		c := Calc{acc: tt.start}
		f(&c, tt.arg)
		if got := c.Result(); got != tt.want {
			t.Errorf("Lookup(%q)(Calc{%d}, %d) = %d, want %d", tt.name, tt.start, tt.arg, got, tt.want)
		}
	}
	for _, bad := range []string{"", "div", "Add", "ADD", " add"} {
		if Lookup(bad) != nil {
			t.Errorf("Lookup(%q) = non-nil, want nil", bad)
		}
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		steps []Step
		want  int
		ok    bool
	}{
		{[]Step{{"add", 2}, {"mul", 10}, {"sub", 1}}, 19, true},
		{nil, 0, true},
		{[]Step{{"mul", 5}}, 0, true},
		{[]Step{{"sub", 3}, {"mul", -2}}, 6, true},
		{[]Step{{"add", 1}, {"pow", 2}}, 0, false},
		{[]Step{{"", 0}}, 0, false},
	}
	for _, tt := range tests {
		got, ok := Run(tt.steps)
		if got != tt.want || ok != tt.ok {
			t.Errorf("Run(%+v) = (%d, %v), want (%d, %v)", tt.steps, got, ok, tt.want, tt.ok)
		}
	}
}

func TestRunAll(t *testing.T) {
	a, b := &Calc{acc: 1}, &Calc{acc: 10}
	steps := []Step{{"add", 1}, {"mul", 3}}
	if !RunAll([]*Calc{a, b}, steps) {
		t.Errorf("RunAll(%+v) = false, want true", steps)
	}
	if a.Result() != 6 || b.Result() != 33 {
		t.Errorf("RunAll(%+v) on [1 10] = [%d %d], want [6 33]", steps, a.Result(), b.Result())
	}

	bad := []Step{{"add", 1}, {"oops", 0}}
	if RunAll([]*Calc{a, b}, bad) {
		t.Errorf("RunAll(%+v) = true, want false", bad)
	}
	if a.Result() != 6 || b.Result() != 33 {
		t.Errorf("failed RunAll(%+v) changed calculators to [%d %d], want [6 33]", bad, a.Result(), b.Result())
	}

	if !RunAll(nil, steps) {
		t.Errorf("RunAll(nil, %+v) = false, want true", steps)
	}
}
