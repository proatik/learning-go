package session

import (
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

// Session expires a fixed TTL after it was started. It reads the time only
// through now, so tests can control the clock. During the bug hunt one of the
// named bugs below is switched on, and your tests must notice each one.
type Session struct {
	start time.Time
	ttl   time.Duration
	now   func() time.Time
}

// New starts a session that expires ttl from now(). A nil now means time.Now.
func New(ttl time.Duration, now func() time.Time) *Session {
	if now == nil {
		now = time.Now
	}
	return &Session{start: now(), ttl: ttl, now: now}
}

// Expired reports whether the TTL has fully passed: at exactly start+ttl the
// session is expired.
func (s *Session) Expired() bool {
	deadline := s.start.Add(s.ttl)
	switch bug {
	case "expires one tick late":
		return s.now().After(deadline)
	case "uses the real clock":
		return !time.Now().Before(deadline)
	}
	return !s.now().Before(deadline)
}

// Remaining returns the time left before expiry, never less than zero.
func (s *Session) Remaining() time.Duration {
	left := s.start.Add(s.ttl).Sub(s.now())
	if left < 0 && bug != "remaining goes negative" {
		return 0
	}
	return left
}

// Touch restarts the TTL from now().
func (s *Session) Touch() {
	switch bug {
	case "touch does nothing":
	case "touch extends from the old deadline":
		s.start = s.start.Add(s.ttl)
	default:
		s.start = s.now()
	}
}

var bugs = []string{
	"expires one tick late",
	"uses the real clock",
	"remaining goes negative",
	"touch does nothing",
	"touch extends from the old deadline",
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
		t.Fatalf("your tests fail on the correct Session:\n%s", out)
	}
	for _, name := range bugs {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
