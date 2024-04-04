package eslog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConfigLeveler tests the Leveler method of the eslog.Config struct.
func TestConfigLeveler(t *testing.T) {
	cases := []struct {
		name     string
		cfg      Config
		expected string
	}{
		{
			name:     "default values",
			cfg:      Config{},
			expected: "INFO",
		},
		{
			name: "filled values",
			cfg: Config{
				Level: "warn",
			},
			expected: "WARN",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.cfg.Leveler().String())
		})
	}
}
