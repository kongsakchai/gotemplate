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
		title  string
		level  string
		enable string
		want   slog.Level
	}{
		{"debug", "debug", "true", slog.LevelDebug},
		{"info", "info", "true", slog.LevelInfo},
		{"warn", "warn", "true", slog.LevelWarn},
		{"error", "error", "true", slog.LevelError},
		{"unknown", "unknown", "true", slog.LevelInfo}, // Default to Info level for unknown levels
		{"disable", "warn", "false", 99},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			defaultLevel := logLevel
			defer resetLogLevel(defaultLevel)

			SetLevel(tc.level, tc.enable)
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
		logger.Info("Test log message", slog.String("key", "value"))
	})

	t.Run("should create a new logger with text format", func(t *testing.T) {
		defaultLogger := slog.Default()
		defer resetLogger(defaultLogger)

		t.Setenv("LOG_FORMAT", "text")
		logger := New()
		assert.NotNil(t, logger)
		logger.Info("Test log message", slog.String("key", "value"))
	})

	t.Run("should create a new logger with JSON format", func(t *testing.T) {
		defaultLogger := slog.Default()
		defer resetLogger(defaultLogger)

		t.Setenv("LOG_FORMAT", "json")
		logger := New()
		assert.NotNil(t, logger)
		logger.Info("Test log message", slog.String("key", "value"))
	})
}

func TestNewReplaceFuncGroup(t *testing.T) {
	t.Run("should return new attr when replace func return true", func(t *testing.T) {
		// arrang
		req := []slog.Attr{
			slog.String("other", "value"),
			slog.String("msg", "value"),
		}
		expected := []slog.Attr{
			slog.String("other", "value"),
			slog.String("message", "value"),
		}
		fn := newReplaceFuncGroup(GCPKeyReplacer)

		//act
		for i, r := range req {
			res := fn(nil, r)

			//assert
			assert.Equal(t, expected[i], res)
		}
	})
}
