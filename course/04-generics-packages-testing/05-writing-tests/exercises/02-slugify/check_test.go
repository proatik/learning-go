package slug

import (
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// Slugify turns a title into a URL slug. During the bug hunt one of the named
// bugs below is switched on, and your tests must notice each one.
func Slugify(s string) string {
	var b strings.Builder
	dash := false // a separator is waiting to be written
	for _, r := range s {
		switch {
		case 'a' <= r && r <= 'z':
		case 'A' <= r && r <= 'Z':
			if bug != "keeps uppercase" {
				r += 'a' - 'A'
			}
		case '0' <= r && r <= '9' && bug != "drops digits":
		default:
			if bug == "does not collapse runs" && b.Len() > 0 {
				b.WriteByte('-')
				continue
			}
			dash = dash || b.Len() > 0 || bug == "keeps a leading dash"
			continue
		}
		if dash {
			b.WriteByte('-')
			dash = false
		}
		b.WriteRune(r)
	}
	if dash && bug == "keeps a trailing dash" {
		b.WriteByte('-')
	}
	return b.String()
}

var bugs = []string{
	"keeps uppercase",
	"drops digits",
	"does not collapse runs",
	"keeps a leading dash",
	"keeps a trailing dash",
}

const bugEnv = "LEARNGO_BUG"

var bug = os.Getenv(bugEnv)

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must fail.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if bug != "" {
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
		t.Fatalf("your tests fail on the correct Slugify:\n%s", out)
	}
	for _, name := range bugs {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
