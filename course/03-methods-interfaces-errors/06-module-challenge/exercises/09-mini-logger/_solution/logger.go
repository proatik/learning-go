package logx

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrUnknownLevel is returned (wrapped) by ParseLevel.
var ErrUnknownLevel = errors.New("logx: unknown level")

// Level is a log severity.
type Level int

// The levels, from least to most severe.
const (
	Debug Level = iota
	Info
	Warn
	Error
)

var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}

// String returns "DEBUG", "INFO", "WARN", "ERROR", or "Level(N)".
func (l Level) String() string {
	if l < 0 || int(l) >= len(levelNames) {
		return fmt.Sprintf("Level(%d)", int(l))
	}
	return levelNames[l]
}

// ParseLevel converts a level name (case-insensitive, surrounding spaces
// ignored) back into a Level.
func ParseLevel(s string) (Level, error) {
	name := strings.ToUpper(strings.TrimSpace(s))
	for i, n := range levelNames {
		if n == name {
			return Level(i), nil
		}
	}
	return 0, fmt.Errorf("%w: %q", ErrUnknownLevel, s)
}

// Logger records messages.
type Logger interface {
	Log(level Level, msg string) error
}

// Compile-time checks.
var (
	_ Logger = (*TextLogger)(nil)
	_ Logger = MultiLogger(nil)
	_ Logger = LevelFilter{}
	_ Logger = (*MemoryLogger)(nil)
)

// TextLogger writes one line per message to an io.Writer.
type TextLogger struct {
	w      io.Writer
	prefix string
}

// NewTextLogger returns a TextLogger writing to w with the given prefix.
func NewTextLogger(w io.Writer, prefix string) *TextLogger {
	return &TextLogger{w: w, prefix: prefix}
}

// Log writes "PREFIX" + level padded to 5 + " " + msg + "\n".
func (l *TextLogger) Log(level Level, msg string) error {
	if _, err := fmt.Fprintf(l.w, "%s%-5s %s\n", l.prefix, level, msg); err != nil {
		return fmt.Errorf("logx: %w", err)
	}
	return nil
}

// With returns a new TextLogger on the same writer whose prefix is this
// logger's prefix followed by prefix.
func (l *TextLogger) With(prefix string) *TextLogger {
	return &TextLogger{w: l.w, prefix: l.prefix + prefix}
}

// MultiLogger sends every message to all of its loggers.
type MultiLogger []Logger

// Log sends the message to every non-nil logger and joins their errors.
func (m MultiLogger) Log(level Level, msg string) error {
	var errs []error
	for _, l := range m {
		if l == nil {
			continue
		}
		if err := l.Log(level, msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// LevelFilter passes on only messages at Min or above.
type LevelFilter struct {
	Min  Level
	Next Logger
}

// Log forwards the message to Next if level >= Min.
func (f LevelFilter) Log(level Level, msg string) error {
	if level < f.Min || f.Next == nil {
		return nil
	}
	return f.Next.Log(level, msg)
}

// Entry is one message kept by a MemoryLogger.
type Entry struct {
	Level Level
	Msg   string
}

// String returns "LEVEL msg".
func (e Entry) String() string { return e.Level.String() + " " + e.Msg }

// MemoryLogger keeps every message in memory (useful in tests).
type MemoryLogger struct {
	Entries []Entry
}

// Log appends the message to Entries.
func (m *MemoryLogger) Log(level Level, msg string) error {
	m.Entries = append(m.Entries, Entry{Level: level, Msg: msg})
	return nil
}

// Logf formats a message with fmt.Sprintf and logs it. A nil logger
// discards the message.
func Logf(l Logger, level Level, format string, args ...any) error {
	if l == nil {
		return nil
	}
	return l.Log(level, fmt.Sprintf(format, args...))
}
