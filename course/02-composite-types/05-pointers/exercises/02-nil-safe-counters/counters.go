package counters

// Increment adds 1 to the int p points to. It does nothing if p is nil.
func Increment(p *int) {
	// TODO
}

// Reset sets the int p points to to 0. It does nothing if p is nil.
func Reset(p *int) {
	// TODO
}

// SafeDeref returns *p, or def if p is nil.
func SafeDeref(p *int, def int) int {
	// TODO
	return 0
}

// Sum adds the ints that ps point to. Nil entries count as 0.
func Sum(ps []*int) int {
	// TODO
	return 0
}
