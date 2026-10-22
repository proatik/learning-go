package rle

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file holds the code under test (Encode and Decode) and the checker.
// During the bug hunt one planted bug is switched on at a time, and the
// tests, fuzz tests and examples in rle_test.go must fail for every one.

// ErrMalformed is returned by Decode for input Encode could not produce.
var ErrMalformed = errors.New("rle: malformed input")

// maxCount is the largest run Decode accepts, so hostile input can't make it
// allocate gigabytes.
const maxCount = 1_000_000

// Encode run-length encodes s as "<count>:<rune>" pairs: "aaab" → "3:a1:b".
// It is defined for valid UTF-8 only.
func Encode(s string) string { return encode(s, active) }

// Decode reverses Encode. It returns ErrMalformed for any other input and
// never panics.
func Decode(s string) (string, error) { return decode(s, active) }

var bugs = []string{
	"Encode drops the last run",
	"Encode works on bytes, not runes",
	"Encode merges letters that differ only in case",
	"Decode reads one-digit counts",
	"Decode cannot decode the ':' character",
	"Decode panics on a count with no character",
	"Decode accepts a zero count",
}

const bugEnv = "LEARNGO_BUG"

var active = os.Getenv(bugEnv)

func encode(s, bug string) string {
	var runes []rune
	if bug == "Encode works on bytes, not runes" {
		for i := range len(s) {
			runes = append(runes, rune(s[i]))
		}
	} else {
		runes = []rune(s)
	}
	same := func(a, b rune) bool { return a == b }
	if bug == "Encode merges letters that differ only in case" {
		same = func(a, b rune) bool { return unicode.ToLower(a) == unicode.ToLower(b) }
	}
	var b strings.Builder
	for i := 0; i < len(runes); {
		j := i
		for j < len(runes) && same(runes[j], runes[i]) {
			j++
		}
		if j == len(runes) && bug == "Encode drops the last run" {
			break
		}
		fmt.Fprintf(&b, "%d:%c", j-i, runes[i])
		i = j
	}
	return b.String()
}

func decode(s, bug string) (string, error) {
	var b strings.Builder
	for len(s) > 0 {
		colon := strings.IndexByte(s, ':')
		if bug == "Decode reads one-digit counts" {
			colon = 1
			if len(s) < 2 || s[1] != ':' {
				return "", ErrMalformed
			}
		}
		if colon <= 0 {
			return "", ErrMalformed
		}
		n, err := strconv.Atoi(s[:colon])
		minCount := 1
		if bug == "Decode accepts a zero count" {
			minCount = 0
		}
		if err != nil || n < minCount || n > maxCount {
			return "", ErrMalformed
		}
		rest := s[colon+1:]
		var r rune
		var size int
		if bug == "Decode panics on a count with no character" {
			r = []rune(rest)[0]
			size = utf8.RuneLen(r)
		} else {
			if rest == "" {
				return "", ErrMalformed
			}
			r, size = utf8.DecodeRuneInString(rest)
		}
		if r == utf8.RuneError && size <= 1 {
			return "", ErrMalformed
		}
		if r == ':' && bug == "Decode cannot decode the ':' character" {
			return "", ErrMalformed
		}
		b.WriteString(strings.Repeat(string(r), n))
		s = rest[size:]
	}
	return b.String(), nil
}

// TestYourTestsCatchEveryBug checks that rle_test.go has a fuzz test, an
// example and a benchmark, that they pass, and that they catch every bug.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if active != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=.")
	names := strings.Fields(list)
	for _, kind := range []string{"Fuzz", "Example", "Benchmark"} {
		if !slices.ContainsFunc(names, func(n string) bool { return strings.HasPrefix(n, kind) }) {
			t.Errorf("%s functions = none, want at least one", kind)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
	const yours = "-test.run=^(Test|Fuzz|Example)"
	if out, err := run("", yours); err != nil {
		t.Fatalf("your tests fail on the correct Encode/Decode:\n%s", out)
	}
	if out, err := run("", "-test.run=^$", "-test.bench=.", "-test.benchtime=1x"); err != nil {
		t.Fatalf("your benchmark fails:\n%s", out)
	}
	for _, name := range bugs {
		if _, err := run(name, yours); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
