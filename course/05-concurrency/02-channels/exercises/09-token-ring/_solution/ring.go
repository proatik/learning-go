package ring

import "sync"

// Ring runs one goroutine per name, connected in a ring by channels, and
// passes a token around it for the given number of turns. Each player appends
// its name to the token. Ring returns the token once the turns are done and
// every player goroutine has exited.
func Ring(names []string, turns int) []string {
	n := len(names)
	if n == 0 {
		return nil
	}
	// One slot per link, so a lone player can pass the token to itself.
	links := make([]chan []string, n)
	for i := range links {
		links[i] = make(chan []string, 1)
	}
	result := make(chan []string)

	var wg sync.WaitGroup
	for i, name := range names {
		in, out := links[i], links[(i+1)%n]
		wg.Go(func() {
			for token := range in {
				if len(token) >= turns {
					result <- token
					continue
				}
				out <- append(token, name)
			}
		})
	}

	links[0] <- nil
	log := <-result
	// The token is out of the ring, so no player will send again: closing every
	// link ends each player's range loop.
	for _, l := range links {
		close(l)
	}
	wg.Wait()
	return log
}
