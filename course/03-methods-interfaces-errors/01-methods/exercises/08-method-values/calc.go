package calc

// Calc is an accumulator calculator. The zero value holds 0.
type Calc struct {
	acc int
}

// Add adds n to the accumulator.
func (c *Calc) Add(n int) {
	// TODO
}

// Sub subtracts n from the accumulator.
func (c *Calc) Sub(n int) {
	// TODO
}

// Mul multiplies the accumulator by n.
func (c *Calc) Mul(n int) {
	// TODO
}

// Result returns the accumulator.
func (c *Calc) Result() int {
	// TODO
	return 0
}

// Bound returns the method values c.Add, c.Sub and c.Mul, in that order.
func Bound(c *Calc) []func(int) {
	// TODO
	return nil
}

// Lookup returns the method expression for "add", "sub" or "mul", or nil
// for any other name.
func Lookup(name string) func(*Calc, int) {
	// TODO
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
	// TODO
	return 0, false
}

// RunAll applies steps to every calculator in cs and reports success. If any
// op is unknown, it returns false and changes nothing.
func RunAll(cs []*Calc, steps []Step) bool {
	// TODO
	return false
}
