package mapreduce

import (
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestPartition(t *testing.T) {
	tests := []struct {
		key  string
		r    int
		want int
	}{
		{"go", 4, 3},
		{"go", 3, 1},
		{"go", 7, 2},
		{"", 4, 1},
		{"the", 4, 0},
		{"the", 7, 3},
		{"cat", 3, 2},
		{"Gopher", 4, 0},
		{"ünï", 7, 3},
		{"anything", 1, 0},
		{"anything", 0, 0},
		{"anything", -3, 0},
	}
	for _, tt := range tests {
		if got := Partition(tt.key, tt.r); got != tt.want {
			t.Errorf("Partition(%q, %d) = %d, want %d", tt.key, tt.r, got, tt.want)
		}
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		docs []string
		want map[string]int
	}{
		{[]string{"the cat", "The dog the end"}, map[string]int{"the": 3, "cat": 1, "dog": 1, "end": 1}},
		{nil, map[string]int{}},
		{[]string{"", "   ", "\t\n"}, map[string]int{}},
		{[]string{"Go go GO", "gopher\tGo\n"}, map[string]int{"go": 4, "gopher": 1}},
		{[]string{"Ünï ÜNÏ", "café, café"}, map[string]int{"ünï": 2, "café,": 1, "café": 1}},
	}
	for _, tt := range tests {
		for _, mr := range [][2]int{{1, 1}, {2, 3}, {8, 8}, {0, -1}} {
			got := WordCount(tt.docs, mr[0], mr[1])
			if got == nil {
				t.Fatalf("WordCount(%q, %d, %d) = nil, want %v", tt.docs, mr[0], mr[1], tt.want)
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("WordCount(%q, %d, %d) = %v, want %v", tt.docs, mr[0], mr[1], got, tt.want)
			}
		}
	}
}

func TestWordCountLarge(t *testing.T) {
	words := strings.Fields("alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu")
	docs := make([]string, 300)
	want := map[string]int{}
	for i := range docs {
		var b strings.Builder
		for j := range 40 {
			w := words[(i*7+j*j)%len(words)]
			if j%3 == 0 {
				w = strings.ToUpper(w)
			}
			b.WriteString(w + " ")
			want[strings.ToLower(w)]++
		}
		docs[i] = b.String()
	}
	done := make(chan map[string]int, 1)
	go func() { done <- WordCount(docs, 6, 4) }()
	select {
	case got := <-done:
		if !maps.Equal(got, want) {
			t.Errorf("WordCount(300 docs, 6, 4) = %v, want %v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("WordCount(300 docs, 6, 4) did not return within 2s (deadlock?)")
	}
}

func TestEveryDocMappedOnce(t *testing.T) {
	docs := make([]string, 50)
	for i := range docs {
		docs[i] = fmt.Sprintf("doc%d", i)
	}
	var mu sync.Mutex
	seen := map[string]int{}
	MapReduce(docs, 5, 3,
		func(doc string, emit func(string, int)) {
			mu.Lock()
			seen[doc]++
			mu.Unlock()
		},
		func(string, []int) int { return 0 },
	)
	for _, d := range docs {
		if seen[d] != 1 {
			t.Errorf("mapFn(%q) called %d times, want 1", d, seen[d])
		}
	}
}

func TestReduceOncePerKeyWithAllValues(t *testing.T) {
	// doc "k=v" emits (k, v) three times: v, v*10 and v*100.
	var docs []string
	for i := range 40 {
		docs = append(docs, fmt.Sprintf("k%d=%d", i%8, i))
	}
	var mu sync.Mutex
	calls := map[string]int{}
	gotValues := map[string][]int{}
	got := MapReduce(docs, 4, 3,
		func(doc string, emit func(string, int)) {
			var k string
			var v int
			fmt.Sscanf(strings.Replace(doc, "=", " ", 1), "%s %d", &k, &v)
			emit(k, v)
			emit(k, v*10)
			emit(k, v*100)
		},
		func(key string, values []int) int {
			mu.Lock()
			calls[key]++
			gotValues[key] = slices.Sorted(slices.Values(values))
			mu.Unlock()
			return len(values)
		},
	)
	for k := range 8 {
		key := fmt.Sprintf("k%d", k)
		var want []int
		for i := k; i < 40; i += 8 {
			want = append(want, i, i*10, i*100)
		}
		slices.Sort(want)
		if calls[key] != 1 {
			t.Errorf("reduceFn(%q, ...) called %d times, want 1", key, calls[key])
		}
		if !slices.Equal(gotValues[key], want) {
			t.Errorf("reduceFn(%q) values = %v, want %v (any order)", key, gotValues[key], want)
		}
		if got[key] != 15 {
			t.Errorf("MapReduce result[%q] = %d, want 15 (len(values))", key, got[key])
		}
	}
	if len(got) != 8 {
		t.Errorf("MapReduce returned %d keys, want 8", len(got))
	}
}

// tracker records how many calls run at the same moment.
type tracker struct{ running, max atomic.Int64 }

func (tr *tracker) enter() {
	n := tr.running.Add(1)
	for {
		m := tr.max.Load()
		if n <= m || tr.max.CompareAndSwap(m, n) {
			return
		}
	}
}

func (tr *tracker) leave() { tr.running.Add(-1) }

func TestMappersRunConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		docs := make([]string, 12)
		for i := range docs {
			docs[i] = fmt.Sprintf("w%d", i)
		}
		var tr tracker
		start := time.Now()
		got := MapReduce(docs, 4, 2,
			func(doc string, emit func(string, int)) {
				tr.enter()
				time.Sleep(time.Second)
				tr.leave()
				emit(doc, 1)
			},
			func(_ string, values []int) int { return len(values) },
		)
		if len(got) != 12 {
			t.Fatalf("MapReduce(12 docs, 4, 2) returned %d keys, want 12", len(got))
		}
		if m := tr.max.Load(); m != 4 {
			t.Errorf("MapReduce(12 docs, 4, 2): max mapFn calls at once = %d, want 4", m)
		}
		if el := time.Since(start); el != 3*time.Second {
			t.Errorf("MapReduce(12 docs, 4 mappers), 1s per doc: took %v, want 3s", el)
		}
	})
}

func TestReducersRunConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var docs []string
		for c := 'a'; c <= 'z'; c++ {
			docs = append(docs, string(c))
		}
		// keys per reducer, computed independently of Partition
		perReducer := map[uint32]int{}
		for _, k := range docs {
			h := fnv.New32a()
			h.Write([]byte(k))
			perReducer[h.Sum32()%3]++
		}
		busiest := slices.Max(slices.Collect(maps.Values(perReducer)))

		var tr tracker
		start := time.Now()
		got := MapReduce(docs, 2, 3,
			func(doc string, emit func(string, int)) { emit(doc, 1) },
			func(_ string, values []int) int {
				tr.enter()
				time.Sleep(time.Second)
				tr.leave()
				return len(values)
			},
		)
		if len(got) != 26 {
			t.Fatalf("MapReduce(26 keys, 2, 3) returned %d keys, want 26", len(got))
		}
		if m := tr.max.Load(); m != 3 {
			t.Errorf("MapReduce(26 keys, 2, 3): max reduceFn calls at once = %d, want 3", m)
		}
		if el, want := time.Since(start), time.Duration(busiest)*time.Second; el != want {
			t.Errorf("MapReduce(26 keys, 3 reducers), 1s per key: took %v, want %v (the busiest reducer owns %d keys)", el, want, busiest)
		}
	})
}
