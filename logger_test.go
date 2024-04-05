package eslog

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// loggerInstance prepares the logger instance for testing.
func loggerInstance(config Config) (Logger, *slog.LevelVar) {
	lvl := new(slog.LevelVar)
	lvl.Set(config.Leveler())

	return New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: config.AddSource,
	},
	), lvl), lvl
}

// TestLogger tests the logging functionality.
func TestLogger(t *testing.T) {
	log, _ := loggerInstance(Config{Level: "trace"})
	expectedMsg := "test message"
	attr := slog.String("test", "some text")

	log.Trace(expectedMsg)
	log.Trace(expectedMsg, attr)

	log.Debug(expectedMsg)
	log.Debug(expectedMsg, attr)

	log.Info(expectedMsg)
	log.Info(expectedMsg, attr)

	log.Warn(expectedMsg)
	log.Warn(expectedMsg, attr)

	log.Error(expectedMsg)
	log.Error(expectedMsg, attr)

	log.Sprintf(expectedMsg)
	log.Print(expectedMsg)
	log.Printf(expectedMsg)
}

// TestLoggerReplaceAttr tests the replacement of attributes in log messages.
func TestLoggerReplaceAttr(t *testing.T) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	log := New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				return func(a slog.Attr) slog.Attr {
					source := a.Value.Any().(slog.Source)

					pwd, err := os.Getwd()
					if err != nil {
						return a
					}

					relPath, err := filepath.Rel(pwd, source.File)
					if err != nil {
						return a
					}

					basePath := filepath.Base(relPath)

					formattedPath := fmt.Sprintf("%s:%d", basePath, source.Line)

					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(formattedPath),
					}
				}(a)
			case slog.LevelKey:
				return func(a slog.Attr) slog.Attr {
					l := a.Value.Any().(slog.Level)
					a.Value = slog.StringValue(l.String())

					return a
				}(a)
			case slog.TimeKey:
				return func(a slog.Attr) slog.Attr {
					const timestampFormat = "2006-01-02 15:04:05.999999999 -0700 -07"

					t, err := time.Parse(timestampFormat, a.Value.String())
					if err != nil {
						return a
					}

					formattedTime := t.Format(time.DateTime)
					a.Value = slog.StringValue(formattedTime)

					return a
				}(a)
			default:
				return a
			}
		},
	},
	), lvl)

	msg := "test message"
	attr := slog.String("test", "some text")

	log.Trace(msg)
	log.Trace(msg, attr)
	log.Debug(msg)
	log.Debug(msg, attr)
	log.Info(msg)
	log.Info(msg, attr)
	log.Warn(msg)
	log.Warn(msg, attr)
	log.Error(msg)
	log.Error(msg, attr)
	log.Sprintf(msg)
	log.Print(msg)
	log.Printf(msg)
}

// TestLoggerSetLevel tests the setting of logging levels.
func TestLoggerSetLevel(t *testing.T) {
	cases := []struct {
		name  string
		level string
		lvl   slog.Level
	}{
		{
			name:  "trace level",
			level: "trace",
			lvl:   LevelTrace,
		},
		{
			name:  "info level",
			level: "info",
			lvl:   slog.LevelInfo,
		},
		{
			name:  "warn level",
			level: "warn",
			lvl:   slog.LevelWarn,
		},
		{
			name:  "error level",
			level: "error",
			lvl:   slog.LevelError,
		},
		{
			name:  "incorrect level",
			level: "bugagaga",
			lvl:   slog.Level(-16),
		},
		{
			name:  "fatal level",
			level: "fatal",
			lvl:   LevelFatal,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			log, lvl := loggerInstance(Config{Level: test.level})
			log.SetLevel(test.lvl)

			assert.Equal(t, test.lvl, lvl.Level())
		})
	}
}
