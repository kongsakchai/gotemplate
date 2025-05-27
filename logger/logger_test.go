package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetLogLevel(defaultLevel slog.Level) {
	logLevel = defaultLevel
}

func TestSetLevel(t *testing.T) {
	testcases := []struct {
		title string
		level string
		want  slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			defaultLevel := logLevel
			defer resetLogLevel(defaultLevel)

			SetLevel(tc.level)
			assert.Equal(t, tc.want, logLevel)
		})
	}
}

func resetLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

func TestNew(t *testing.T) {
	t.Run("should create a new logger with default settings", func(t *testing.T) {
		defaultLogger := slog.Default()
		defer resetLogger(defaultLogger)

		logger := New()
		assert.NotNil(t, logger)
	})

	t.Run("should create a new logger with text format", func(t *testing.T) {
		defaultLogger := slog.Default()
		defer resetLogger(defaultLogger)

		t.Setenv("LOG_FORMAT", "text")
		logger := New()
		assert.NotNil(t, logger)
	})

	t.Run("should create a new logger with JSON format", func(t *testing.T) {
		defaultLogger := slog.Default()
		defer resetLogger(defaultLogger)

		t.Setenv("LOG_FORMAT", "json")
		logger := New()
		assert.NotNil(t, logger)
	})
}
