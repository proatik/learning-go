package ring

// Ring runs one goroutine per name, connected in a ring by channels, and
// passes a token around it for the given number of turns. Each player appends
// its name to the token. Ring returns the token once the turns are done and
// every player goroutine has exited.
func Ring(names []string, turns int) []string {
	// TODO
	return nil
}
