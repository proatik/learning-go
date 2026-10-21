package logx_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	logx "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/09-mini-logger"
)

// Every type must satisfy Logger (note which ones need a pointer).
var (
	_ logx.Logger = (*logx.TextLogger)(nil)
	_ logx.Logger = logx.MultiLogger(nil)
	_ logx.Logger = logx.LevelFilter{}
	_ logx.Logger = (*logx.MemoryLogger)(nil)
)

var errDisk = errors.New("disk full")

// brokenWriter fails every write.
type brokenWriter struct{}

func (brokenWriter) Write(p []byte) (int, error) { return 0, errDisk }

func newText(t *testing.T, sb *strings.Builder, prefix string) *logx.TextLogger {
	t.Helper()
	l := logx.NewTextLogger(sb, prefix)
	if l == nil {
		t.Fatalf("NewTextLogger(w, %q) = nil, want a logger", prefix)
	}
	return l
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		l    logx.Level
		want string
	}{
		{logx.Debug, "DEBUG"},
		{logx.Info, "INFO"},
		{logx.Warn, "WARN"},
		{logx.Error, "ERROR"},
		{logx.Level(4), "Level(4)"},
		{logx.Level(-2), "Level(-2)"},
	}
	for _, tt := range tests {
		if got := fmt.Sprint(tt.l); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", int(tt.l), got, tt.want)
		}
	}
	if !(logx.Debug < logx.Info && logx.Info < logx.Warn && logx.Warn < logx.Error) {
		t.Errorf("levels must be ordered Debug < Info < Warn < Error")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want logx.Level
	}{
		{"DEBUG", logx.Debug},
		{"info", logx.Info},
		{"  Warn ", logx.Warn},
		{"eRRoR", logx.Error},
	}
	for _, tt := range tests {
		got, err := logx.ParseLevel(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, %v, want %v, nil", tt.in, got, err, tt.want)
		}
	}
	for _, in := range []string{"", "warning", "fatal", "Level(1)", "1"} {
		_, err := logx.ParseLevel(in)
		if !errors.Is(err, logx.ErrUnknownLevel) {
			t.Errorf("ParseLevel(%q) error = %v, want %v", in, err, logx.ErrUnknownLevel)
			continue
		}
		if want := fmt.Sprintf("logx: unknown level: %q", in); err.Error() != want {
			t.Errorf("ParseLevel(%q) error = %q, want %q", in, err.Error(), want)
		}
	}
	// Round trip.
	for _, l := range []logx.Level{logx.Debug, logx.Info, logx.Warn, logx.Error} {
		if got, err := logx.ParseLevel(l.String()); got != l || err != nil {
			t.Errorf("ParseLevel(%q) = %v, %v, want %v, nil", l.String(), got, err, l)
		}
	}
}

func TestTextLogger(t *testing.T) {
	var sb strings.Builder
	l := newText(t, &sb, "app: ")
	lines := []struct {
		level logx.Level
		msg   string
	}{
		{logx.Info, "started"},
		{logx.Debug, "config loaded"},
		{logx.Warn, "disk 91% full"}, // % in a message is not a format verb
		{logx.Error, "crash: héllo"},
		{logx.Level(7), "odd"},
	}
	for _, ln := range lines {
		if err := l.Log(ln.level, ln.msg); err != nil {
			t.Fatalf("Log(%v, %q) error = %v, want nil", ln.level, ln.msg, err)
		}
	}
	want := "app: INFO  started\n" +
		"app: DEBUG config loaded\n" +
		"app: WARN  disk 91% full\n" +
		"app: ERROR crash: héllo\n" +
		"app: Level(7) odd\n"
	if sb.String() != want {
		t.Errorf("output =\n%s\nwant\n%s", sb.String(), want)
	}
}

func TestTextLoggerWith(t *testing.T) {
	var sb strings.Builder
	root := newText(t, &sb, "")
	db := root.With("db: ")
	if db == nil {
		t.Fatalf("With(%q) = nil, want a logger", "db: ")
	}
	tx := db.With("tx: ")
	if tx == nil {
		t.Fatalf("With(%q) = nil, want a logger", "tx: ")
	}
	root.Log(logx.Info, "a")
	db.Log(logx.Info, "b")
	tx.Log(logx.Error, "c")
	db.Log(logx.Warn, "d") // With must not change the parent
	want := "INFO  a\ndb: INFO  b\ndb: tx: ERROR c\ndb: WARN  d\n"
	if sb.String() != want {
		t.Errorf("output =\n%s\nwant\n%s", sb.String(), want)
	}
}

func TestTextLoggerWriteError(t *testing.T) {
	l := logx.NewTextLogger(brokenWriter{}, "")
	if l == nil {
		t.Fatalf("NewTextLogger(...) = nil, want a logger")
	}
	err := l.Log(logx.Info, "x")
	if !errors.Is(err, errDisk) {
		t.Fatalf("Log to a broken writer error = %v, want %v", err, errDisk)
	}
	if err.Error() != "logx: disk full" {
		t.Errorf("Log to a broken writer error = %q, want %q", err.Error(), "logx: disk full")
	}
}

func TestMemoryLoggerAndEntry(t *testing.T) {
	m := &logx.MemoryLogger{}
	m.Log(logx.Warn, "w")
	logx.Logf(m, logx.Info, "user %s logged in %d times", "ada", 3)
	want := []logx.Entry{{logx.Warn, "w"}, {logx.Info, "user ada logged in 3 times"}}
	if !slices.Equal(m.Entries, want) {
		t.Fatalf("Entries = %v, want %v", m.Entries, want)
	}
	if got := fmt.Sprint(m.Entries); got != "[WARN w INFO user ada logged in 3 times]" {
		t.Errorf("fmt.Sprint(Entries) = %q, want %q", got, "[WARN w INFO user ada logged in 3 times]")
	}
	if err := logx.Logf(nil, logx.Info, "dropped %d", 1); err != nil {
		t.Errorf("Logf(nil, ...) = %v, want nil", err)
	}
}

func TestLevelFilter(t *testing.T) {
	m := &logx.MemoryLogger{}
	f := logx.LevelFilter{Min: logx.Warn, Next: m}
	for _, l := range []logx.Level{logx.Debug, logx.Info, logx.Warn, logx.Error, logx.Info} {
		if err := f.Log(l, l.String()); err != nil {
			t.Fatalf("LevelFilter.Log(%v) error = %v, want nil", l, err)
		}
	}
	if got := fmt.Sprint(m.Entries); got != "[WARN WARN ERROR ERROR]" {
		t.Errorf("filtered entries = %s, want [WARN WARN ERROR ERROR]", got)
	}
	// A filter with no Next discards everything.
	if err := (logx.LevelFilter{Min: logx.Debug}).Log(logx.Error, "x"); err != nil {
		t.Errorf("LevelFilter{Next: nil}.Log = %v, want nil", err)
	}
	// Errors from Next pass through unchanged.
	bad := logx.LevelFilter{Min: logx.Info, Next: logx.NewTextLogger(brokenWriter{}, "")}
	if err := bad.Log(logx.Info, "x"); !errors.Is(err, errDisk) {
		t.Errorf("LevelFilter over a broken logger error = %v, want %v", err, errDisk)
	}
	if err := bad.Log(logx.Debug, "x"); err != nil {
		t.Errorf("LevelFilter: filtered message error = %v, want nil (never reaches Next)", err)
	}
}

func TestMultiLogger(t *testing.T) {
	var sb strings.Builder
	text := newText(t, &sb, "")
	m := &logx.MemoryLogger{}
	all := logx.MultiLogger{
		text,
		nil, // skipped
		logx.LevelFilter{Min: logx.Error, Next: m},
	}
	all.Log(logx.Info, "hello")
	all.Log(logx.Error, "oops")
	if sb.String() != "INFO  hello\nERROR oops\n" {
		t.Errorf("text output = %q, want %q", sb.String(), "INFO  hello\nERROR oops\n")
	}
	if got := fmt.Sprint(m.Entries); got != "[ERROR oops]" {
		t.Errorf("memory entries = %s, want [ERROR oops]", got)
	}
	if err := (logx.MultiLogger{}).Log(logx.Info, "x"); err != nil {
		t.Errorf("empty MultiLogger.Log = %v, want nil", err)
	}
	if err := logx.MultiLogger(nil).Log(logx.Info, "x"); err != nil {
		t.Errorf("nil MultiLogger.Log = %v, want nil", err)
	}
}

func TestMultiLoggerErrors(t *testing.T) {
	m := &logx.MemoryLogger{}
	all := logx.MultiLogger{
		logx.NewTextLogger(brokenWriter{}, "a"),
		m, // still receives the message
		logx.NewTextLogger(brokenWriter{}, "b"),
	}
	err := all.Log(logx.Warn, "x")
	if err == nil {
		t.Fatalf("MultiLogger.Log error = nil, want both write errors joined")
	}
	if err.Error() != "logx: disk full\nlogx: disk full" {
		t.Errorf("MultiLogger.Log error = %q, want %q", err.Error(), "logx: disk full\nlogx: disk full")
	}
	if !errors.Is(err, errDisk) {
		t.Errorf("errors.Is(err, errDisk) = false, want true")
	}
	if len(m.Entries) != 1 {
		t.Errorf("memory logger got %d entries, want 1 (one failing logger must not stop the rest)", len(m.Entries))
	}
	// Multi loggers nest: a MultiLogger is a Logger.
	nested := logx.MultiLogger{all, m}
	nested.Log(logx.Info, "y")
	if len(m.Entries) != 3 {
		t.Errorf("after a nested MultiLogger: memory logger has %d entries, want 3", len(m.Entries))
	}
}

func ExampleTextLogger() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	base := logx.NewTextLogger(w, "shop: ")
	log := logx.MultiLogger{
		logx.LevelFilter{Min: logx.Info, Next: base},
		logx.LevelFilter{Min: logx.Error, Next: base.With("ALERT: ")},
	}
	logx.Logf(log, logx.Debug, "cache warm")
	logx.Logf(log, logx.Info, "order %d paid", 42)
	logx.Logf(log, logx.Error, "payment %q failed", "p-7")
	// Output:
	// shop: INFO  order 42 paid
	// shop: ERROR payment "p-7" failed
	// shop: ALERT: ERROR payment "p-7" failed
}
