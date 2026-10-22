package checkout

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// This file holds the code under test (Total) and the checker. During the
// bug hunt one planted bug is switched on at a time, and your tests in
// checkout_test.go must fail for every one of them.

// Errors returned by Total.
var (
	ErrEmptyCart      = errors.New("checkout: empty cart")
	ErrBadLine        = errors.New("checkout: bad line")
	ErrUnknownCode    = errors.New("checkout: unknown code")
	ErrCodeNotAllowed = errors.New("checkout: code not allowed")
)

// Line is one cart line: Qty units at Cents each.
type Line struct {
	SKU   string
	Cents int64
	Qty   int
}

// Total returns the cart total in cents. See the README for the rules.
func Total(lines []Line, code string) (int64, error) { return total(lines, code, active) }

var bugs = []string{
	"empty cart is accepted",
	"zero quantity is accepted",
	"free items are rejected",
	"bulk discount starts at 11 units",
	"bulk discount rounds up",
	"bulk discount covers the whole cart",
	"SAVE5 needs more than 2500",
	"SAVE5 looks at the price before bulk discount",
	"HALF discounts the cheapest unit",
	"HALF halves the whole line",
	"codes are case-sensitive",
	"unknown codes are ignored",
}

const bugEnv = "LEARNGO_BUG"

var active = os.Getenv(bugEnv)

// total is the real implementation; bug switches on one planted bug.
func total(lines []Line, code, bug string) (int64, error) {
	if len(lines) == 0 && bug != "empty cart is accepted" {
		return 0, ErrEmptyCart
	}
	var before, sum, maxUnit int64
	var maxQty int
	bulkCart := false
	for i, l := range lines {
		minQty := 1
		if bug == "zero quantity is accepted" {
			minQty = 0
		}
		badCents := l.Cents < 0
		if bug == "free items are rejected" {
			badCents = l.Cents <= 0
		}
		if l.Qty < minQty || badCents {
			return 0, fmt.Errorf("%w: line %d", ErrBadLine, i)
		}
		line := l.Cents * int64(l.Qty)
		before += line
		bulkAt := 10
		if bug == "bulk discount starts at 11 units" {
			bulkAt = 11
		}
		if l.Qty >= bulkAt {
			bulkCart = true
			if bug == "bulk discount rounds up" {
				line = (line*90 + 99) / 100
			} else if bug != "bulk discount covers the whole cart" {
				line = line * 90 / 100
			}
		}
		sum += line
		pick := l.Cents > maxUnit
		if bug == "HALF discounts the cheapest unit" {
			pick = i == 0 || l.Cents < maxUnit
		}
		if pick {
			maxUnit, maxQty = l.Cents, l.Qty
		}
	}
	if bug == "bulk discount covers the whole cart" && bulkCart {
		sum = sum * 90 / 100
	}

	if bug != "codes are case-sensitive" {
		code = strings.ToUpper(code)
	}
	switch code {
	case "":
	case "SAVE5":
		base := sum
		if bug == "SAVE5 looks at the price before bulk discount" {
			base = before
		}
		if base < 2500 || (bug == "SAVE5 needs more than 2500" && base == 2500) {
			return 0, fmt.Errorf("%w: SAVE5 needs 2500", ErrCodeNotAllowed)
		}
		sum -= 500
	case "HALF":
		off := maxUnit / 2
		if bug == "HALF halves the whole line" {
			off = maxUnit * int64(maxQty) / 2
		}
		sum -= off
	default:
		if bug != "unknown codes are ignored" {
			return 0, fmt.Errorf("%w: %q", ErrUnknownCode, code)
		}
	}
	return sum, nil
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must
// fail.
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
	list, _ := run("", "-test.list=^Test")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	if out, err := run("", "-test.run=^Test"); err != nil {
		t.Fatalf("your tests fail on the correct Total:\n%s", out)
	}
	for _, name := range bugs {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
