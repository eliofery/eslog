/*
Package eslog is an extension of the [pkg/log/slog.Logger] logging package.
This package formats logging messages in a human-readable form,
contains two additional logging levels Trace and Fatal.

Includes [Logger] that combines functionality for setting logging level.

// Demonstrates how to configure and use the logger with
// various logging levels and message types.
// It sets up the logger with custom configuration, including logging level and output format.
// Then, it overrides the logging level for demonstration purposes.
// Finally, it logs messages at different levels using the configured logger.

	// Configures the logger with the specified settings.
	config := eslog.Config{
		Level:     slog.LevelInfo,
		AddSource: true,
		JSON:      false,
	}

	// Sets the logging level.
	lvl := new(slog.LevelVar)
	lvl.Set(config.Leveler())

	// Creates a new instance of eslog with the configured settings.
	logger := eslog.New(pretty.NewHandler(os.Stdout, &pretty.HandlerOptions{
		SlogOptions: &slog.HandlerOptions{
			Level:     lvl,
			AddSource: config.AddSource,
		},
		JSON: config.JSON,
	}), lvl)

	// Overrides the logging level.
	logger.SetLevel(LevelTrace)
	logger.SetLevel(slog.LevelDebug)
	logger.SetLevel(slog.LevelInfo)
	logger.SetLevel(slog.LevelWarn)
	logger.SetLevel(slog.LevelError)
	logger.SetLevel(LevelFatal)

	// Logs messages at different levels.
	logger.Trace("Trace example", slog.Any("message", "trace message"))
	logger.Debug("Debug example", slog.Any("message", "debug message"))
	logger.Info("Info example", slog.Any("message", "info message"))
	logger.Warn("Warn example", slog.Any("message", "warn message"))
	logger.Error("Error example", slog.Any("message", "error message"))
	logger.Fatal("Fatal example", slog.Any("message", "fatal message"))
*/
package eslog
