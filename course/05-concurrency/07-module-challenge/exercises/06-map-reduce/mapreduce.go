package mapreduce

// Partition returns the reducer, in [0, r), that owns key: the 32-bit FNV-1a
// hash of key modulo r.
func Partition(key string, r int) int {
	// TODO
	return 0
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
	// TODO
	return nil
}

// WordCount counts the lower-cased, whitespace-separated words in docs using
// MapReduce.
func WordCount(docs []string, m, r int) map[string]int {
	// TODO
	return nil
}
