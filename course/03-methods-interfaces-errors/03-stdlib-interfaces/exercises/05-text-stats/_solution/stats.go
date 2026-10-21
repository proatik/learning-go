package textstats

import (
	"bufio"
	"io"
	"unicode/utf8"
)

// count returns the number of tokens produced by split.
func count(r io.Reader, split bufio.SplitFunc) (int, error) {
	sc := bufio.NewScanner(r)
	sc.Split(split)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

// LineCount returns the number of lines in r.
func LineCount(r io.Reader) (int, error) {
	return count(r, bufio.ScanLines)
}

// WordCount returns the number of white-space separated words in r.
func WordCount(r io.Reader) (int, error) {
	return count(r, bufio.ScanWords)
}

// Longest returns the first line of r with the most runes.
func Longest(r io.Reader) (string, error) {
	sc := bufio.NewScanner(r)
	best, bestLen := "", -1
	for sc.Scan() {
		line := sc.Text()
		if n := utf8.RuneCountInString(line); n > bestLen {
			best, bestLen = line, n
		}
	}
	return best, sc.Err()
}
