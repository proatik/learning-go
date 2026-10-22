package results

// Result is the outcome of one call: a value or an error.
type Result struct {
	Value int
	Err   error
}

// Run calls f for each input in order in a single goroutine and sends each
// outcome. It stops after the first error. The channel has room for every
// result, so the goroutine never blocks, and it is closed at the end.
func Run(inputs []int, f func(int) (int, error)) <-chan Result {
	out := make(chan Result, len(inputs))
	go func() {
		defer close(out)
		for _, n := range inputs {
			v, err := f(n)
			out <- Result{Value: v, Err: err}
			if err != nil {
				return
			}
		}
	}()
	return out
}

// Collect returns the values received from ch until it is closed. At the
// first error it stops and returns the values before it and that error. A nil
// ch yields nil, nil.
func Collect(ch <-chan Result) ([]int, error) {
	if ch == nil {
		return nil, nil
	}
	var vals []int
	for r := range ch {
		if r.Err != nil {
			return vals, r.Err
		}
		vals = append(vals, r.Value)
	}
	return vals, nil
}
