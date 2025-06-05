package logger

import (
	"log/slog"
	"os"
)

var logLevel slog.Level

func SetLevel(level string) {
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo // Default to Info level if an unknown level is provided
	}
}

func init() {
	SetLevel(os.Getenv("LOG_LEVEL"))
}

func New() *slog.Logger {
	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = NewTextHandler(os.Stdout, &HandlerOptions{
			Level:       logLevel,
			ReplaceAttr: replaceAttr,
			TimeFormat:  "[2006/01/02 15:04:05]",
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       logLevel,
			ReplaceAttr: replaceAttr,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
