package logger

import (
	"log/slog"
)

const (
	MessageKey = "message"
)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		return slog.Attr{
			Key:   MessageKey,
			Value: a.Value,
		}
	}

	return a
}
