package logger

import (
	"log/slog"
	"os"

	"github.com/kongsakchai/paint"
)

var logLevel slog.Level

func SetLevel(level string, enable string) {
	if enable != "true" {
		logLevel = 99
		return
	}

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
	SetLevel(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_ENABLE"))
}

type ReplaceFunc func(groups []string, a slog.Attr) (slog.Attr, bool)

func New(replaceAttrs ...ReplaceFunc) *slog.Logger {
	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = paint.NewTextHandler(os.Stdout, &paint.HandlerOptions{
			Level:       logLevel,
			ReplaceAttr: newReplaceFuncGroup(replaceAttrs...),
			TimeFormat:  "[2006/01/02 15:04:05]",
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       logLevel,
			ReplaceAttr: newReplaceFuncGroup(replaceAttrs...),
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func newReplaceFuncGroup(replaceAttrs ...ReplaceFunc) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, replace := range replaceAttrs {
			if v, ok := replace(groups, a); ok {
				return v
			}
		}
		return a
	}
}
