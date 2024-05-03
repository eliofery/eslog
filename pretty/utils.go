package pretty

import (
	"encoding/json"
	"fmt"
	"github.com/eliofery/eslog"
	"github.com/fatih/color"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// NameColor represents a combination of a name and a color function for formatting.
type NameColor struct {
	name  string
	color func(format string, a ...any) string
}

// replaceAttrFn is a function type for modifying logging attributes.
type replaceAttrFn func(groups []string, attr slog.Attr) slog.Attr

// LevelConf is a mapping of logging levels to their respective names and colors for formatting.
var LevelConf = map[slog.Leveler]NameColor{
	eslog.LevelTrace: {"TRACE", color.WhiteString},
	slog.LevelDebug:  {"DEBUG", color.HiWhiteString},
	slog.LevelInfo:   {"INFO", color.HiGreenString},
	slog.LevelWarn:   {"WARN", color.HiYellowString},
	slog.LevelError:  {"ERROR", color.HiMagentaString},
	eslog.LevelFatal: {"FATAL", color.HiRedString},
}

// dateTime formats the date and time of a logging record.
func (p *pretty) dateTime(r slog.Record, replaceAttrFn replaceAttrFn) string {
	timeStr := r.Time.Format(time.DateTime)
	if p.JSON {
		timeStr = r.Time.Format(time.RFC3339Nano)
	}

	if replaceAttrFn != nil {
		if attr := replaceAttrFn(nil, slog.Time(slog.TimeKey, r.Time.Round(0))); attr.Key != "" {
			timeStr = attr.Value.String()
		}
	}

	if p.JSON {
		return timeStr
	}

	return color.WhiteString(timeStr)
}

// level formats the logging level of a logging record.
func (p *pretty) level(r slog.Record, replaceAttrFn replaceAttrFn) string {
	level, ok := LevelConf[r.Level]
	if !ok {
		level.name = r.Level.String()
	}

	if replaceAttrFn != nil {
		if attr := replaceAttrFn(nil, slog.Any(slog.LevelKey, r.Level)); attr.Key != "" {
			level.name = attr.Value.String()
		}
	}

	if p.JSON {
		return level.name
	}

	return level.color(level.name)
}

// source formats the source of a logging record.
func (p *pretty) source(r slog.Record, replaceAttrFn replaceAttrFn) string {
	var pathSource string

	if !p.SlogOptions.AddSource {
		return pathSource
	}

	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	if f.File == "" {
		return pathSource
	}

	var src slog.Source
	src.Function = f.Function
	src.File = f.File
	src.Line = f.Line

	pwd, err := os.Getwd()
	if err != nil {
		return pathSource
	}

	relPath, err := filepath.Rel(pwd, src.File)
	if err != nil {
		return pathSource
	}

	basePath := filepath.Base(relPath)
	pathSource = fmt.Sprintf("%s:%d", basePath, src.Line)

	if replaceAttrFn != nil {
		if attr := replaceAttrFn(nil, slog.Any(slog.SourceKey, src)); attr.Key != "" {
			pathSource = attr.Value.String()
		}
	}

	return pathSource
}

// message extracts the message from a logging record.
func (p *pretty) message(r slog.Record) string {
	if p.JSON {
		return r.Message
	}

	return color.HiWhiteString(r.Message)
}

// attrs formats the attributes of a logging record.
func (p *pretty) attrs(r slog.Record) string {
	attrs := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()

		return true
	})

	var jsonAttrs []byte
	if p.JSON {
		jsonAttrs, _ = json.Marshal(attrs)
	} else {
		jsonAttrs, _ = json.MarshalIndent(attrs, "", "  ")
	}

	return string(jsonAttrs)
}
