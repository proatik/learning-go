package isomorphic

import "strings"

// IsIsomorphic reports whether the runes of s can be replaced one-to-one
// to produce t.
func IsIsomorphic(s, t string) bool {
	rs, rt := []rune(s), []rune(t)
	if len(rs) != len(rt) {
		return false
	}
	fwd := map[rune]rune{}  // rune of s → rune of t
	back := map[rune]rune{} // rune of t → rune of s
	for i := range rs {
		a, b := rs[i], rt[i]
		if m, ok := fwd[a]; ok && m != b {
			return false // a already maps to something else
		}
		if m, ok := back[b]; ok && m != a {
			return false // b is already the image of another rune
		}
		fwd[a] = b
		back[b] = a
	}
	return true
}

// WordPattern reports whether the words of s follow pattern: each rune of
// pattern maps to exactly one word, and each word to exactly one rune.
func WordPattern(pattern, s string) bool {
	letters, words := []rune(pattern), strings.Fields(s)
	if len(letters) != len(words) {
		return false
	}
	fwd := map[rune]string{}
	back := map[string]rune{}
	for i, r := range letters {
		w := words[i]
		if m, ok := fwd[r]; ok && m != w {
			return false
		}
		if m, ok := back[w]; ok && m != r {
			return false
		}
		fwd[r] = w
		back[w] = r
	}
	return true
}
