package mapreduce

import (
	"hash/fnv"
	"strings"
	"sync"
)

// pair is one emitted key/value.
type pair struct {
	key   string
	value int
}

// Partition returns the reducer, in [0, r), that owns key: the 32-bit FNV-1a
// hash of key modulo r.
func Partition(key string, r int) int {
	r = max(r, 1)
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() % uint32(r))
}

// MapReduce runs mapFn over docs with m mapper goroutines, routes each
// emitted pair to one of r reducer goroutines by Partition, and returns
// reduceFn's result for every key.
func MapReduce(
	docs []string,
	m, r int,
	mapFn func(doc string, emit func(key string, value int)),
	reduceFn func(key string, values []int) int,
) map[string]int {
	m, r = max(m, 1), max(r, 1)

	toReducer := make([]chan pair, r)
	for i := range toReducer {
		toReducer[i] = make(chan pair, 64)
	}
	emit := func(key string, value int) {
		toReducer[Partition(key, r)] <- pair{key, value}
	}

	// Reducers start first: mappers block on their channels otherwise.
	partial := make([]map[string]int, r)
	var reducers sync.WaitGroup
	for i := range r {
		reducers.Go(func() {
			groups := map[string][]int{}
			for p := range toReducer[i] {
				groups[p.key] = append(groups[p.key], p.value)
			}
			out := make(map[string]int, len(groups))
			for k, vs := range groups {
				out[k] = reduceFn(k, vs)
			}
			partial[i] = out
		})
	}

	jobs := make(chan string)
	var mappers sync.WaitGroup
	for range m {
		mappers.Go(func() {
			for doc := range jobs {
				mapFn(doc, emit)
			}
		})
	}
	for _, doc := range docs {
		jobs <- doc
	}
	close(jobs)
	mappers.Wait()
	for _, ch := range toReducer {
		close(ch)
	}
	reducers.Wait()

	result := map[string]int{}
	for _, part := range partial {
		for k, v := range part {
			result[k] = v
		}
	}
	return result
}

// WordCount counts the lower-cased, whitespace-separated words in docs using
// MapReduce.
func WordCount(docs []string, m, r int) map[string]int {
	return MapReduce(docs, m, r,
		func(doc string, emit func(string, int)) {
			for _, w := range strings.Fields(strings.ToLower(doc)) {
				emit(w, 1)
			}
		},
		func(_ string, values []int) int {
			sum := 0
			for _, v := range values {
				sum += v
			}
			return sum
		},
	)
}
