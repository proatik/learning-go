package wordcount

import "strings"

// words splits a document into words. Use it for every document; the tests
// replace it with a slow version to check documents are counted
// concurrently.
var words = strings.Fields

// Count returns how often each lowercased word appears across docs,
// counting each document in its own goroutine.
func Count(docs []string) map[string]int {
	// TODO
	return nil
}
