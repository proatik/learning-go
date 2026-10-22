package gen

// Count returns a channel that receives 1, 2, …, n and is then closed.
// For n <= 0 the channel is closed without any values.
func Count(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

// Runes returns a channel that receives each rune of s in order and is then
// closed.
func Runes(s string) <-chan rune {
	out := make(chan rune)
	go func() {
		defer close(out)
		for _, r := range s {
			out <- r
		}
	}()
	return out
}
