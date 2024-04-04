package eslog

import (
	"log/slog"
)

// Default logging level.
const levelDefault = "info"

// Level names for common levels.
var Level = map[string]slog.Level{
	"trace": LevelTrace,
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
	"fatal": LevelFatal,
}

// Config contains configuration options for logging.
type Config struct {
	Level     string `yaml:"level" env-default:"info"`
	AddSource bool   `yaml:"add-source"`
	JSON      bool   `yaml:"json"`
}

// Leveler returns the logging level specified in the configuration.
func (c *Config) Leveler() slog.Level {
	level, ok := Level[c.Level]
	if !ok {
		return Level[levelDefault]
	}

	return level
}
