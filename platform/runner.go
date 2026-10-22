package main

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

const testTimeout = 2 * time.Minute

type TestResult struct {
	Passed   bool        `json:"passed"`
	Output   string      `json:"output"` // plain `go test -v` text
	Duration string      `json:"duration"`
	Command  string      `json:"command"`
	Report   *TestReport `json:"report"`
}

// testCommand is the command shown to the learner (plain text output).
func testCommand(ex *Exercise) []string {
	return []string{"go", "test", "-count=1", "-v", "./" + ex.RelDir + "/..."}
}

func joinCmd(args []string) string { return strings.Join(args, " ") }

func runTests(ctx context.Context, root string, ex *Exercise) TestResult {
	ctx, cancel := context.WithTimeout(ctx, testTimeout)
	defer cancel()

	// Run with -json so the result can be shown case by case.
	args := testCommand(ex)
	cmd := exec.CommandContext(ctx, "go", "test", "-count=1", "-json", "./"+ex.RelDir+"/...") // /... includes nested packages
	cmd.Dir = root

	start := time.Now()
	out, err := cmd.CombinedOutput()
	timedOut := errors.Is(ctx.Err(), context.DeadlineExceeded)
	report, text := parseTestJSON(out, timedOut)
	res := TestResult{
		Passed:   err == nil && report.Verdict == "accepted",
		Output:   text,
		Duration: time.Since(start).Round(time.Millisecond).String(),
		Command:  strings.Join(args, " "),
		Report:   report,
	}
	if timedOut {
		res.Output += "\n\n[platform] tests timed out after " + testTimeout.String() + " (infinite loop?)"
	}
	return res
}
