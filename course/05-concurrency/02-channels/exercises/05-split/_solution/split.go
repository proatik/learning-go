package split

// Split returns two channels: evens receives the even values from in and odds
// the odd ones, each in input order. Both are closed after in is closed.
func Split(in <-chan int) (evens, odds <-chan int) {
	ev := make(chan int)
	od := make(chan int)
	go func() {
		defer close(ev)
		defer close(od)
		for v := range in {
			if v%2 == 0 {
				ev <- v
			} else {
				od <- v
			}
		}
	}()
	return ev, od
}
