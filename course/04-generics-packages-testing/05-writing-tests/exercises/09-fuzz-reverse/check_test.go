package rev

import (
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// Reverse returns s with its runes in reverse order. During the bug hunt it
// is swapped for a broken version, and your fuzz test's seeds must notice
// each one.
func Reverse(s string) string { return impl(s) }

var impl = reverse

func reverse(s string) string {
	r := []rune(s)
	slices.Reverse(r)
	return string(r)
}

var bugs = map[string]func(string) string{
	"reverses bytes, not runes": func(s string) string {
		b := []byte(s)
		slices.Reverse(b)
		return string(b)
	},
	"drops the last rune": func(s string) string {
		r := []rune(s)
		if len(r) == 0 {
			return ""
		}
		return reverse(string(r[:len(r)-1]))
	},
	"keeps the first rune in place": func(s string) string {
		r := []rune(s)
		if len(r) == 0 {
			return ""
		}
		return string(r[0]) + reverse(string(r[1:]))
	},
	"reverses each word separately": func(s string) string {
		words := strings.Split(s, " ")
		for i, w := range words {
			words[i] = reverse(w)
		}
		return strings.Join(words, " ")
	},
}

const bugEnv = "LEARNGO_BUG"

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// TestYourFuzzCatchesEveryBug re-runs your fuzz test's seed corpus once per
// bug: each run must fail.
func TestYourFuzzCatchesEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourFuzzCatchesEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=^Fuzz")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return strings.HasPrefix(n, "Fuzz") }) {
		t.Fatal("no Fuzz functions found: write one")
	}
	// Without -fuzz, a fuzz test runs once per seed, like a table test.
	if out, err := run("", "-test.run=^Fuzz"); err != nil {
		t.Fatalf("your fuzz test fails on the correct Reverse:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, "-test.run=^Fuzz"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
