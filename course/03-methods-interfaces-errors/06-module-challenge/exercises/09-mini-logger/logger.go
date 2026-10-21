package logx

import (
	"errors"
	"io"
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

// String returns "DEBUG", "INFO", "WARN", "ERROR", or "Level(N)".
func (l Level) String() string {
	// TODO
	return ""
}

// ParseLevel converts a level name (case-insensitive, surrounding spaces
// ignored) back into a Level.
func ParseLevel(s string) (Level, error) {
	// TODO
	return 0, nil
}

// Logger records messages.
type Logger interface {
	Log(level Level, msg string) error
}

// TextLogger writes one line per message to an io.Writer.
type TextLogger struct {
	w      io.Writer
	prefix string
}

// NewTextLogger returns a TextLogger writing to w with the given prefix.
func NewTextLogger(w io.Writer, prefix string) *TextLogger {
	// TODO
	return nil
}

// Log writes "PREFIX" + level padded to 5 + " " + msg + "\n".
func (l *TextLogger) Log(level Level, msg string) error {
	// TODO
	return nil
}

// With returns a new TextLogger on the same writer whose prefix is this
// logger's prefix followed by prefix.
func (l *TextLogger) With(prefix string) *TextLogger {
	// TODO
	return nil
}

// MultiLogger sends every message to all of its loggers.
type MultiLogger []Logger

// Log sends the message to every non-nil logger and joins their errors.
func (m MultiLogger) Log(level Level, msg string) error {
	// TODO
	return nil
}

// LevelFilter passes on only messages at Min or above.
type LevelFilter struct {
	Min  Level
	Next Logger
}

// Log forwards the message to Next if level >= Min.
func (f LevelFilter) Log(level Level, msg string) error {
	// TODO
	return nil
}

// Entry is one message kept by a MemoryLogger.
type Entry struct {
	Level Level
	Msg   string
}

// String returns "LEVEL msg".
func (e Entry) String() string {
	// TODO
	return ""
}

// MemoryLogger keeps every message in memory (useful in tests).
type MemoryLogger struct {
	Entries []Entry
}

// Log appends the message to Entries.
func (m *MemoryLogger) Log(level Level, msg string) error {
	// TODO
	return nil
}

// Logf formats a message with fmt.Sprintf and logs it. A nil logger
// discards the message.
func Logf(l Logger, level Level, format string, args ...any) error {
	// TODO
	return nil
}
