package crawler

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
)

// Fetcher fetches a page and returns the links found on it.
type Fetcher interface {
	Fetch(ctx context.Context, url string) (links []string, err error)
}

// failure is a fetch error together with the URL it belongs to.
type failure struct {
	url string
	err error
}

// Crawl fetches start and every page reachable from it in at most depth
// links, each exactly once, with at most limit fetches running at a time
// (limit < 1 means 1). It returns the successfully fetched URLs, sorted,
// and the joined fetch errors (sorted by URL), plus ctx.Err() if ctx ended
// the crawl early.
func Crawl(ctx context.Context, f Fetcher, start string, depth, limit int) ([]string, error) {
	if depth < 0 {
		return nil, nil
	}
	sem := make(chan struct{}, max(limit, 1))
	var (
		mu       sync.Mutex
		seen     = map[string]bool{start: true}
		visited  []string
		failures []failure
	)

	// Level by level (breadth first): every URL is first reached by its
	// shortest path, so a slow page can't hide a link that is within depth.
	level := []string{start}
	for d := 0; d <= depth && len(level) > 0 && ctx.Err() == nil; d++ {
		var next []string
		var wg sync.WaitGroup
	launch:
		for _, url := range level {
			if ctx.Err() != nil {
				break
			}
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				break launch
			}
			wg.Go(func() {
				defer func() { <-sem }()
				links, err := f.Fetch(ctx, url)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					failures = append(failures, failure{url, fmt.Errorf("fetch %s: %w", url, err)})
					return
				}
				visited = append(visited, url)
				for _, l := range links {
					if !seen[l] {
						seen[l] = true
						next = append(next, l)
					}
				}
			})
		}
		wg.Wait()
		level = next
	}

	slices.Sort(visited)
	slices.SortFunc(failures, func(a, b failure) int { return cmp.Compare(a.url, b.url) })
	var errs []error
	for _, fl := range failures {
		errs = append(errs, fl.err)
	}
	if err := ctx.Err(); err != nil {
		errs = append(errs, err)
	}
	return visited, errors.Join(errs...)
}
