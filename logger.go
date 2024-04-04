package eslog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

// Names for common levels.
// Complements standard slog levels.
// See: https://pkg.go.dev/golang.org/x/exp/slog#Level
const (
	LevelTrace slog.Level = -8
	LevelFatal slog.Level = 12
)

// Leveler is an interface for setting the logging level.
type Leveler interface {
	// SetLevel sets the logging level.
	SetLevel(level slog.Level)
}

// Loggerer is an interface for logging messages at different levels.
// Complements standard slog methods.
type Loggerer interface {
	// Trace logs a message at trace level.
	Trace(msg string, args ...any)

	// Debug logs a message at trace level.
	Debug(msg string, args ...any)

	// Info logs a message at trace level.
	Info(msg string, args ...any)

	// Warn logs a message at trace level.
	Warn(msg string, args ...any)

	// Error logs a message at trace level.
	Error(msg string, args ...any)

	// Fatal logs a message at trace level.
	Fatal(msg string, args ...any)
}

// Printerer is an interface for printing messages.
type Printerer interface {
	// Sprintf formats according to a format specifier and returns the resulting string.
	Sprintf(msg string, args ...any) string

	// Fatalf formats according to a format specifier and prints the resulting string,
	// then exits the program with a non-zero exit status.
	Fatalf(msg string, args ...any)

	// Print prints the message.
	Print(msg string, args ...any)

	// Printf formats according to a format specifier and prints the resulting string.
	Printf(msg string, args ...any)
}

// Logger is an interface that combines functionality for setting logging level.
type Logger interface {
	// Leveler is an interface for setting the logging level.
	Leveler

	// Loggerer is an interface for logging messages at different levels.
	// Complements standard slog methods.
	Loggerer

	// Printerer is an interface for printing messages.
	Printerer
}

// logger represents a structure that implements the [Logger] interface.
type logger struct {
	logger *slog.Logger
	level  *slog.LevelVar
}

// New creates a new instance of [Logger].
// See: https://github.com/golang/go/issues/59145#issuecomment-1481920720
func New(handler slog.Handler, lvl *slog.LevelVar) Logger {
	log := slog.New(handler)
	slog.SetDefault(log)

	return &logger{
		logger: log,
		level:  lvl,
	}
}

// SetLevel sets the logging level.
func (l *logger) SetLevel(level slog.Level) {
	l.level.Set(level)
}

// log message at the specified logging level.
// See: https://github.com/golang/go/issues/59145#issuecomment-1481920720
func (l *logger) log(level slog.Level, msg string, args ...any) {
	if !l.logger.Enabled(context.Background(), level) {
		return
	}

	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)

	_ = l.logger.Handler().Handle(context.Background(), r)
}

// Trace logs a message at trace level.
func (l *logger) Trace(msg string, args ...any) {
	l.log(LevelTrace, msg, args...)
}

// Debug logs a message at trace level.
func (l *logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

// Info logs a message at trace level.
func (l *logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

// Warn logs a message at trace level.
func (l *logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

// Error logs a message at trace level.
func (l *logger) Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, args...)
}

// Fatal logs a message at fatal level.
func (l *logger) Fatal(msg string, args ...any) {
	l.log(LevelFatal, msg, args...)
	os.Exit(1)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func (l *logger) Sprintf(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

// Fatalf formats according to a format specifier and prints the resulting string,
// then exits the program with a non-zero exit status.
func (l *logger) Fatalf(msg string, args ...any) {
	l.Fatal(l.Sprintf(l.removeLineBreak(msg), args...))
}

// Print prints the message.
func (l *logger) Print(msg string, args ...any) {
	l.log(l.level.Level(), l.removeLineBreak(msg), args...)
}

// Printf formats according to a format specifier and prints the resulting string.
func (l *logger) Printf(msg string, args ...any) {
	l.log(l.level.Level(), l.Sprintf(l.removeLineBreak(msg), args...))
}

// removeLineBreak removes line breaks from the given message and replaces them with spaces.
func (l *logger) removeLineBreak(msg string) string {
	return strings.Replace(msg, "\n", " ", -1)
}
