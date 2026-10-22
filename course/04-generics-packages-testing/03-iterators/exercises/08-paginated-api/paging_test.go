package paging

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"testing"
)

// collect2 ranges over seq and returns each pair formatted as "k:v",
// breaking after max pairs when max > 0 (and always after 10000). A panic
// during the loop is returned instead of crashing.
func collect2[K, V any](seq iter.Seq2[K, V], max int) (got []string, panicked any) {
	defer func() { panicked = recover() }()
	for k, v := range seq {
		got = append(got, fmt.Sprintf("%v:%v", k, v))
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

var errTimeout = errors.New("timeout")

// api simulates a paginated endpoint and records every page requested.
type api struct {
	pages   [][]string
	failAt  int // page number that fails with errTimeout; 0 means never
	fetched []int
}

func (a *api) fetch(page int) ([]string, bool, error) {
	a.fetched = append(a.fetched, page)
	if len(a.fetched) > len(a.pages)+5 {
		panic(fmt.Sprintf("runaway: fetched pages %v (stop when more is false)", a.fetched))
	}
	if page == a.failAt {
		return []string{"ignored"}, true, errTimeout
	}
	if page < 1 || page > len(a.pages) {
		return nil, false, nil
	}
	return a.pages[page-1], page < len(a.pages), nil
}

func TestItems(t *testing.T) {
	tests := []struct {
		name    string
		pages   [][]string
		failAt  int
		want    []string
		fetched []int
	}{
		{"3 pages", [][]string{{"a", "b"}, {}, {"c"}}, 0, []string{"a:<nil>", "b:<nil>", "c:<nil>"}, []int{1, 2, 3}},
		{"1 page", [][]string{{"x"}}, 0, []string{"x:<nil>"}, []int{1}},
		{"empty only page", [][]string{{}}, 0, nil, []int{1}},
		{"page 2 fails", [][]string{{"a", "b"}, {"c"}, {"d"}}, 2, []string{"a:<nil>", "b:<nil>", ":timeout"}, []int{1, 2}},
		{"page 1 fails", [][]string{{"a"}}, 1, []string{":timeout"}, []int{1}},
	}
	for _, tt := range tests {
		a := &api{pages: tt.pages, failAt: tt.failAt}
		seq := Items(a.fetch)
		if seq == nil {
			t.Fatalf("Items(%s) = nil, want a non-nil iterator", tt.name)
		}
		got, p := collect2(seq, -1)
		if p != nil {
			t.Fatalf("range over Items(%s) panicked: %v", tt.name, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Items(%s) = %q, want %q", tt.name, got, tt.want)
		}
		if !slices.Equal(a.fetched, tt.fetched) {
			t.Errorf("Items(%s) fetched pages = %v, want %v", tt.name, a.fetched, tt.fetched)
		}
	}
}

func TestItemsErrorValue(t *testing.T) {
	a := &api{pages: [][]string{{"a"}, {"b"}}, failAt: 2}
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("range over Items(page 2 fails) panicked: %v", p)
		}
	}()
	var last error
	for _, err := range Items(a.fetch) {
		last = err
	}
	if !errors.Is(last, errTimeout) {
		t.Errorf("last error from Items(page 2 fails) = %v, want %v", last, errTimeout)
	}
}

func TestItemsIsLazy(t *testing.T) {
	a := &api{pages: [][]string{{"a", "b"}, {"c", "d"}, {"e"}}}
	seq := Items(a.fetch)
	if len(a.fetched) != 0 {
		t.Errorf("Items(fetch) before ranging: fetched pages = %v, want []", a.fetched)
	}
	for _, n := range []int{1, 2, 3} {
		a.fetched = nil
		got, p := collect2(seq, n)
		if p != nil {
			t.Fatalf("break after %d items panicked: %v (stop when yield returns false)", n, p)
		}
		if len(got) != n {
			t.Fatalf("break after %d items of [[a b] [c d] [e]]: items = %q, want %d items", n, got, n)
		}
		want := []int{1}
		if n == 3 {
			want = []int{1, 2}
		}
		if !slices.Equal(a.fetched, want) {
			t.Errorf("break after %d items of [[a b] [c d] [e]]: fetched pages = %v, want %v", n, a.fetched, want)
		}
	}
}

func TestItemsInts(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("CollectAll(Items(3 int pages)) panicked: %v", p)
		}
	}()
	fetch := func(page int) ([]int, bool, error) {
		if page > 10 {
			panic(fmt.Sprintf("runaway: fetched page %d of 3 (stop when more is false)", page))
		}
		return []int{page * 10, page*10 + 1}, page < 3, nil
	}
	got, err := CollectAll(Items(fetch))
	if want := []int{10, 11, 20, 21, 30, 31}; !slices.Equal(got, want) || err != nil {
		t.Errorf("CollectAll(Items(3 int pages)) = %v, %v, want %v, <nil>", got, err, want)
	}
}

func TestCollectAll(t *testing.T) {
	a := &api{pages: [][]string{{"a", "b"}, {}, {"c"}}}
	got, err := CollectAll(Items(a.fetch))
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) || err != nil {
		t.Errorf("CollectAll(Items(3 pages)) = %q, %v, want %q, <nil>", got, err, want)
	}

	a = &api{pages: [][]string{{"a", "b"}, {"c"}}, failAt: 2}
	got, err = CollectAll(Items(a.fetch))
	if want := []string{"a", "b"}; !slices.Equal(got, want) || !errors.Is(err, errTimeout) {
		t.Errorf("CollectAll(Items(page 2 fails)) = %q, %v, want %q, %v", got, err, want, errTimeout)
	}

	empty := func(yield func(string, error) bool) {}
	got, err = CollectAll(empty)
	if len(got) != 0 || err != nil {
		t.Errorf("CollectAll(empty) = %q, %v, want [], <nil>", got, err)
	}
}

func TestCollectAllStopsAtError(t *testing.T) {
	after := 0
	seq := func(yield func(int, error) bool) {
		if !yield(1, nil) || !yield(0, errTimeout) {
			return
		}
		after++
		yield(2, nil)
	}
	got, err := CollectAll(seq)
	if !slices.Equal(got, []int{1}) || !errors.Is(err, errTimeout) {
		t.Errorf("CollectAll((1, nil), (0, timeout), (2, nil)) = %v, %v, want [1], %v", got, err, errTimeout)
	}
	if after != 0 {
		t.Errorf("CollectAll kept reading after the first error: extra values pulled = %d, want 0", after)
	}
}
