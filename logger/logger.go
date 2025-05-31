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
	}
}

func init() {
	logLevel = slog.LevelInfo
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		SetLevel(level)
	}
}

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	}

	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = &textHandler{
			opts: handlerOptions{
				level:       logLevel,
				replaceAttr: replaceAttr,
				timeFormat:  "[2006/01/02 15:04:05]",
			},
			w: os.Stdout,
		}
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
