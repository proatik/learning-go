package grade

import (
	"maps"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// Grade turns a score from 0 to 100 into a letter grade. During the bug hunt
// it is swapped for a broken version, and your tests must notice each one.
func Grade(score int) string { return impl(score) }

var impl = grade

func grade(score int) string {
	switch {
	case score < 0 || score > 100:
		return "invalid"
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	}
	return "F"
}

var bugs = map[string]func(int) string{
	"90 is a B": func(s int) string {
		if s == 90 {
			return "B"
		}
		return grade(s)
	},
	"59 is a D": func(s int) string {
		if s == 59 {
			return "D"
		}
		return grade(s)
	},
	"no D band": func(s int) string {
		if g := grade(s); g != "D" {
			return g
		}
		return "F"
	},
	"accepts 101": func(s int) string {
		if s == 101 {
			return "A"
		}
		return grade(s)
	},
	"negative scores get an F": func(s int) string {
		if s < 0 {
			return "F"
		}
		return grade(s)
	},
}

const bugEnv = "LEARNGO_BUG"

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

const minSubtests = 6

var subtestLine = regexp.MustCompile(`(?m)^=== RUN\s+(Test\w*/\S+)`)

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must
// fail. It also counts your named subtests.
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
	out, err := run("", "-test.run=^Test", "-test.v")
	if err != nil {
		t.Fatalf("your tests fail on the correct Grade:\n%s", out)
	}
	names := map[string]bool{}
	for _, m := range subtestLine.FindAllStringSubmatch(out, -1) {
		names[m[1]] = true
	}
	if len(names) < minSubtests {
		t.Errorf("named subtests = %d, want at least %d (use t.Run)", len(names), minSubtests)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
