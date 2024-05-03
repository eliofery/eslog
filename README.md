# eslog : 🌈 slog.Handler that writes prettier logs

Package `eslog` implements a [slog.Handler](https://pkg.go.dev/log/slog#Handler)
that writes prettier logs. Two modes are supported: **pretty** and **json**.

Supported logging levels:
- Trace
- Debug
- Info
- Warn
- Error
- Fatal

## Pretty

Logging formatting in **Pretty** mode. In this mode, it's easy to navigate to the file position where the log was called.

![Prettier example](https://raw.githubusercontent.com/eliofery/eslog/assets/screen.png)

## JSON

Logging formatting in **JSON** mode.

![JSON example](https://raw.githubusercontent.com/eliofery/eslog/assets/screen2.png)

## Installation

```bash
go get github.com/eliofery/eslog
```

## Usage

```go
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
```

### Customize Attributes

`ReplaceAttr` can be used to alter or drop attributes. See [`slog.HandlerOptions`](https://pkg.go.dev/log/slog#HandlerOptions) for details.

```go
eslog.New(pretty.NewHandler(os.Stdout, &pretty.HandlerOptions{
    SlogOptions: &slog.HandlerOptions{
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
}), lvl)
```
## Support for the standard package slog

You can still use methods from the standard package slog.

```go
// NewTextHandler
eslog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}), lvl)

// NewJSONHandler
eslog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}), lvl)
```

## Inspiration

I was inspired to write this package by the [slog](https://pkg.go.dev/log/slog) and [tint](https://github.com/lmittmann/tint).
