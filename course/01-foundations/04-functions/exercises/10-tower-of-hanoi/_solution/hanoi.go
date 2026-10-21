package hanoi

// Hanoi moves n disks from peg from to peg to, using via as the spare peg.
// It calls move(src, dst) for every single-disk move (unless move is nil)
// and returns the total number of moves.
func Hanoi(n int, from, to, via string, move func(from, to string)) int {
	if n <= 0 {
		return 0 // base case: nothing to move
	}
	// 1. Park the n-1 smaller disks on the spare peg.
	count := Hanoi(n-1, from, via, to, move)
	// 2. Move the largest disk to its destination.
	if move != nil {
		move(from, to)
	}
	count++
	// 3. Put the n-1 smaller disks back on top of it.
	count += Hanoi(n-1, via, to, from, move)
	return count
}
