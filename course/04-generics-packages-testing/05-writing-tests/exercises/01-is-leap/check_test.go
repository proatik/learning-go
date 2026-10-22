package leap

import (
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// IsLeap reports whether year is a leap year in the Gregorian calendar.
// During the bug hunt it is swapped for a broken version, and your tests must
// notice each one.
func IsLeap(year int) bool { return impl(year) }

var impl = func(y int) bool { return y%4 == 0 && (y%100 != 0 || y%400 == 0) }

var bugs = map[string]func(int) bool{
	"only the 4-year rule":      func(y int) bool { return y%4 == 0 },
	"forgets the 400-year rule": func(y int) bool { return y%4 == 0 && y%100 != 0 },
	"every century is leap":     func(y int) bool { return y%4 == 0 && y%100 != 0 || y%100 == 0 },
	"off by one year": func(y int) bool {
		y--
		return y%4 == 0 && (y%100 != 0 || y%400 == 0)
	},
}

const bugEnv = "LEARNGO_BUG"

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must fail.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=^Test")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	if out, err := run("", "-test.run=^Test"); err != nil {
		t.Fatalf("your tests fail on the correct IsLeap:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
