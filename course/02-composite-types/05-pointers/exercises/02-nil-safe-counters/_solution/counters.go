package counters

// Increment adds 1 to the int p points to. It does nothing if p is nil.
func Increment(p *int) {
	if p == nil {
		return
	}
	*p++ // same as (*p)++: increments the variable, not the pointer
}

// Reset sets the int p points to to 0. It does nothing if p is nil.
func Reset(p *int) {
	if p != nil {
		*p = 0
	}
}

// SafeDeref returns *p, or def if p is nil.
func SafeDeref(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}

// Sum adds the ints that ps point to. Nil entries count as 0.
func Sum(ps []*int) int {
	total := 0
	for _, p := range ps {
		total += SafeDeref(p, 0)
	}
	return total
}
