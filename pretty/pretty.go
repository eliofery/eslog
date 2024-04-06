// Package pretty is a handler for slog that formats messages beautifully.
//
// See: https://betterstack.com/community/guides/logging/logging-in-go/
package pretty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
)

// jsonMessage used for outputting logging messages in JSON format.
type jsonMessage struct {
	Time    string `json:"time,omitempty"`
	Level   string `json:"level,omitempty"`
	Source  string `json:"source,omitempty"`
	Message string `json:"msg,omitempty"`
	Attrs   string `json:"attrs,omitempty"`
}

// HandlerOptions contains options for configuring a logging handler.
type HandlerOptions struct {
	SlogOptions *slog.HandlerOptions
	JSON        bool
}

// Pretty is an interface for handling logging records with pretty formatting.
type Pretty interface {
	// Handle formats and handles a logging record.
	Handle(ctx context.Context, r slog.Record) error
}

// pretty represents a structure that implements the [Pretty] interface.
type pretty struct {
	slog.Handler
	*log.Logger
	*HandlerOptions
}

// NewHandler creates a new logging handler [pkg/log/slog.Handler] with the specified output writer and options.
func NewHandler(out io.Writer, opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	return &pretty{
		Handler:        slog.NewTextHandler(out, opts.SlogOptions),
		Logger:         log.New(out, "", 0),
		HandlerOptions: opts,
	}
}

// Handle formats and handles a logging record.
func (p *pretty) Handle(_ context.Context, record slog.Record) error {
	fn := p.SlogOptions.ReplaceAttr

	time := p.dateTime(record, fn)
	level := p.level(record, fn)
	source := p.source(record, fn)
	message := p.message(record)
	attrs := p.attrs(record)

	if p.JSON {
		jsonResult, _ := json.Marshal(jsonMessage{
			Time:    time,
			Level:   level,
			Source:  source,
			Message: message,
			Attrs:   attrs,
		})

		p.Logger.Println(string(jsonResult))

		return nil
	}

	var result string
	if source != "" {
		result = fmt.Sprintf("%s %s %s %s", time, level, source, message)
	} else {
		result = fmt.Sprintf("%s %s %s", time, level, message)
	}

	// Skip empty attrs if attributes contain only brackets "{}" or "".
	if len(attrs) > 2 {
		result = fmt.Sprintf("%s %s", result, attrs)
	}

	p.Logger.Println(result)

	return nil
}
