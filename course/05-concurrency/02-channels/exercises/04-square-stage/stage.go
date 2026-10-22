package stage

// Generate returns a channel that receives each of nums in order and is then
// closed.
func Generate(nums ...int) <-chan int {
	// TODO
	return nil
}

// Square returns a channel that receives v*v for every v received from in.
// It is closed after in is closed.
func Square(in <-chan int) <-chan int {
	// TODO
	return nil
}
