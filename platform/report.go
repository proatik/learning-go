package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// TestReport is a LeetCode-style summary of one `go test -json` run.
type TestReport struct {
	Verdict     string       `json:"verdict"` // accepted | wrong_answer | compile_error | runtime_error | timeout
	Passed      int          `json:"passed"`
	Total       int          `json:"total"`
	Tests       []*TestCase  `json:"tests"`
	BuildErrors []BuildError `json:"buildErrors,omitempty"`
	BuildOutput string       `json:"buildOutput,omitempty"` // build lines that aren't file:line:col errors
}

type BuildError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col"`
	Message string `json:"message"`
}

// TestCase is one top-level Test or Example function.
type TestCase struct {
	Name     string    `json:"name"`
	Display  string    `json:"display"`
	Passed   bool      `json:"passed"`
	Panicked bool      `json:"panicked"`
	Elapsed  float64   `json:"elapsed"`
	Failures []Failure `json:"failures,omitempty"`
	Log      string    `json:"log,omitempty"` // other output: prints, panics, example got/want
	done     bool
}

// Failure is one t.Errorf/t.Fatalf message, split into input/got/want when
// it follows the "F(x) = got, want want" convention used by the tests.
type Failure struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
	Input   string `json:"input,omitempty"`
	Got     string `json:"got,omitempty"`
	Want    string `json:"want,omitempty"`
}

type testEvent struct {
	Action     string
	Test       string
	Output     string
	OutputType string
	Elapsed    float64
}

var (
	reBuildErr = regexp.MustCompile(`(?:^|[\\/])([\w.-]+\.go):(\d+):(\d+): (.*)$`)
	reTestErr  = regexp.MustCompile(`^\s*([\w.-]+\.go):(\d+): (.*)$`)
	reEqWant   = regexp.MustCompile(`^(.+?) = (.*), want (.*)$`)
	reGotWant  = regexp.MustCompile(`^(?:(.*?):?\s+)?got (.*?),? want (.*)$`)
)

// parseTestJSON turns `go test -json` output into a report plus the plain
// text that `go test -v` would have printed.
func parseTestJSON(out []byte, timedOut bool) (*TestReport, string) {
	rep := &TestReport{}
	var raw strings.Builder
	cases := map[string]*TestCase{}
	buildFailed := false
	var last *Failure

	sc := bufio.NewScanner(bytes.NewReader(out))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		var ev testEvent
		if len(line) == 0 || line[0] != '{' || json.Unmarshal(line, &ev) != nil {
			raw.Write(line) // e.g. toolchain download messages
			raw.WriteByte('\n')
			continue
		}
		raw.WriteString(ev.Output)

		switch ev.Action {
		case "build-output":
			text := strings.TrimRight(ev.Output, "\n")
			if strings.HasPrefix(text, "#") || text == "" {
				continue
			}
			if m := reBuildErr.FindStringSubmatch(text); m != nil {
				ln, _ := strconv.Atoi(m[2])
				col, _ := strconv.Atoi(m[3])
				rep.BuildErrors = append(rep.BuildErrors, BuildError{File: m[1], Line: ln, Col: col, Message: m[4]})
			} else {
				rep.BuildOutput += text + "\n"
			}
			continue
		case "build-fail":
			buildFailed = true
			continue
		}
		if ev.Test == "" {
			if strings.Contains(ev.Output, "panic:") {
				rep.Verdict = "runtime_error"
			}
			continue
		}

		top, _, _ := strings.Cut(ev.Test, "/")
		tc := cases[top]
		if tc == nil {
			tc = &TestCase{Name: top, Display: displayName(top)}
			cases[top] = tc
			rep.Tests = append(rep.Tests, tc)
		}
		switch ev.Action {
		case "pass", "fail", "skip":
			if ev.Test == top {
				tc.done = true
				tc.Passed = ev.Action != "fail"
				tc.Elapsed = ev.Elapsed
			}
		case "output":
			text := strings.TrimRight(ev.Output, "\n")
			switch {
			case ev.OutputType == "frame" || isFrame(text):
				// === RUN / --- FAIL lines carry no information for the learner
			case ev.OutputType == "error-continue" && last != nil:
				last.Message += "\n" + strings.TrimSpace(text)
				splitFailure(last)
			case ev.OutputType == "error" || (ev.OutputType == "" && reTestErr.MatchString(text)):
				if m := reTestErr.FindStringSubmatch(text); m != nil {
					ln, _ := strconv.Atoi(m[2])
					tc.Failures = append(tc.Failures, Failure{File: m[1], Line: ln, Message: m[3]})
				} else {
					tc.Failures = append(tc.Failures, Failure{Message: strings.TrimSpace(text)})
				}
				last = &tc.Failures[len(tc.Failures)-1]
				splitFailure(last)
			default:
				if strings.HasPrefix(strings.TrimSpace(text), "panic:") {
					tc.Panicked = true
				}
				tc.Log += ev.Output
			}
		}
	}

	for _, tc := range rep.Tests {
		if !tc.done { // a panic or timeout stopped the binary mid-test
			tc.Passed = false
		}
		if tc.Passed {
			rep.Passed++
		} else if strings.HasPrefix(tc.Name, "Example") && len(tc.Failures) == 0 {
			if got, want, ok := exampleDiff(tc.Log); ok {
				tc.Failures = append(tc.Failures, Failure{Message: "output mismatch", Got: got, Want: want})
				tc.Log = ""
			}
		}
		if tc.Panicked {
			rep.Verdict = "runtime_error"
		}
	}
	rep.Total = len(rep.Tests)

	switch {
	case timedOut:
		rep.Verdict = "timeout"
	case buildFailed:
		rep.Verdict = "compile_error"
	case rep.Verdict == "runtime_error":
	case rep.Total > 0 && rep.Passed == rep.Total:
		rep.Verdict = "accepted"
	default:
		rep.Verdict = "wrong_answer"
	}
	return rep, raw.String()
}

func isFrame(s string) bool {
	t := strings.TrimSpace(s)
	for _, p := range []string{"=== RUN", "=== PAUSE", "=== CONT", "=== NAME", "--- FAIL", "--- PASS", "--- SKIP"} {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}

// splitFailure fills Input/Got/Want from messages like
// `Greet("Rob") = "", want "Hello, Rob!"` or `Sum: got 3, want 4`.
func splitFailure(f *Failure) {
	f.Input, f.Got, f.Want = "", "", ""
	if strings.Contains(f.Message, "\n") {
		return
	}
	if m := reEqWant.FindStringSubmatch(f.Message); m != nil {
		f.Input, f.Got, f.Want = m[1], m[2], m[3]
		return
	}
	if m := reGotWant.FindStringSubmatch(f.Message); m != nil {
		f.Input, f.Got, f.Want = strings.TrimSpace(m[1]), m[2], m[3]
	}
}

func exampleDiff(log string) (got, want string, ok bool) {
	g := strings.Index(log, "got:\n")
	w := strings.LastIndex(log, "want:\n")
	if g < 0 || w < g {
		return "", "", false
	}
	return strings.TrimSuffix(log[g+5:w], "\n"), strings.TrimSuffix(log[w+6:], "\n"), true
}

// displayName turns TestCircleArea into "CircleArea" and
// ExamplePrintProfile_inactive into "PrintProfile (inactive) · output".
func displayName(name string) string {
	switch {
	case name == "Example":
		return "Program output"
	case strings.HasPrefix(name, "Example"):
		base, suffix, _ := strings.Cut(strings.TrimPrefix(name, "Example"), "_")
		base = strings.TrimPrefix(base, "_")
		if suffix != "" {
			base += " (" + suffix + ")"
		}
		return strings.TrimSpace(base) + " · output"
	case strings.HasPrefix(name, "Test"):
		return strings.TrimPrefix(name, "Test")
	}
	return name
}
