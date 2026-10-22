package crawler

import "context"

// Fetcher fetches a page and returns the links found on it.
type Fetcher interface {
	Fetch(ctx context.Context, url string) (links []string, err error)
}

// Crawl fetches start and every page reachable from it in at most depth
// links, each exactly once, with at most limit fetches running at a time
// (limit < 1 means 1). It returns the successfully fetched URLs, sorted,
// and the joined fetch errors (sorted by URL), plus ctx.Err() if ctx ended
// the crawl early.
func Crawl(ctx context.Context, f Fetcher, start string, depth, limit int) ([]string, error) {
	// TODO
	return nil, nil
}
