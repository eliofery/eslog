package pretty

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

// TestNewHandler tests the creation of new request handlers.
func TestNewHandler(t *testing.T) {
	prettyHandler := NewHandler(os.Stdout, nil)
	assert.Implements(t, (*slog.Handler)(nil), prettyHandler)
}

// TestHandle tests the handling of HTTP requests.
func TestHandle(t *testing.T) {
	cases := []struct {
		name   string
		json   bool
		source bool
	}{
		{
			name: "with json",
			json: true,
		},
		{
			name: "without json",
			json: false,
		},
		{
			name:   "with source",
			source: true,
		},
		{
			name:   "without source",
			source: false,
		},
		{
			name:   "with",
			json:   true,
			source: true,
		},
		{
			name:   "without",
			json:   false,
			source: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var pc uintptr

			r := slog.NewRecord(time.Now(), slog.LevelDebug, "test message", pc)
			prettyHandler := NewHandler(os.Stdout, &HandlerOptions{
				SlogOptions: &slog.HandlerOptions{AddSource: test.source},
				JSON:        test.json,
			})

			err := prettyHandler.Handle(context.Background(), r)
			assert.NoError(t, err)
		})
	}
}
