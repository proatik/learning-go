package wordcount

import (
	"strings"
	"sync"
)

// words splits a document into words. Use it for every document; the tests
// replace it with a slow version to check documents are counted
// concurrently.
var words = strings.Fields

// Count returns how often each lowercased word appears across docs,
// counting each document in its own goroutine.
func Count(docs []string) map[string]int {
	local := make([]map[string]int, len(docs))
	var wg sync.WaitGroup
	for i, doc := range docs {
		wg.Go(func() {
			m := make(map[string]int)
			for _, w := range words(doc) {
				m[strings.ToLower(w)]++
			}
			local[i] = m
		})
	}
	wg.Wait()

	total := make(map[string]int)
	for _, m := range local {
		for w, n := range m {
			total[w] += n
		}
	}
	return total
}
