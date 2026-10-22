package gen

// Count returns a channel that receives 1, 2, …, n and is then closed.
// For n <= 0 the channel is closed without any values.
func Count(n int) <-chan int {
	// TODO
	return nil
}

// Runes returns a channel that receives each rune of s in order and is then
// closed.
func Runes(s string) <-chan rune {
	// TODO
	return nil
}
