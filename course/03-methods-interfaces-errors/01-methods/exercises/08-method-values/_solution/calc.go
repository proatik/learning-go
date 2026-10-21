package calc

// Calc is an accumulator calculator. The zero value holds 0.
type Calc struct {
	acc int
}

// Add adds n to the accumulator.
func (c *Calc) Add(n int) {
	c.acc += n
}

// Sub subtracts n from the accumulator.
func (c *Calc) Sub(n int) {
	c.acc -= n
}

// Mul multiplies the accumulator by n.
func (c *Calc) Mul(n int) {
	c.acc *= n
}

// Result returns the accumulator.
func (c *Calc) Result() int {
	return c.acc
}

// Bound returns the method values c.Add, c.Sub and c.Mul, in that order.
func Bound(c *Calc) []func(int) {
	return []func(int){c.Add, c.Sub, c.Mul}
}

// Lookup returns the method expression for "add", "sub" or "mul", or nil
// for any other name.
func Lookup(name string) func(*Calc, int) {
	switch name {
	case "add":
		return (*Calc).Add
	case "sub":
		return (*Calc).Sub
	case "mul":
		return (*Calc).Mul
	}
	return nil
}

// Step is one named operation with its argument.
type Step struct {
	Op  string
	Arg int
}

// Run applies steps to a fresh Calc and returns its result; ok is false if a
// step has an unknown op.
func Run(steps []Step) (int, bool) {
	var c Calc
	if !RunAll([]*Calc{&c}, steps) {
		return 0, false
	}
	return c.Result(), true
}

// RunAll applies steps to every calculator in cs and reports success. If any
// op is unknown, it returns false and changes nothing.
func RunAll(cs []*Calc, steps []Step) bool {
	ops := make([]func(*Calc, int), len(steps))
	for i, s := range steps {
		ops[i] = Lookup(s.Op)
		if ops[i] == nil {
			return false
		}
	}
	for _, c := range cs {
		for i, op := range ops {
			op(c, steps[i].Arg)
		}
	}
	return true
}
