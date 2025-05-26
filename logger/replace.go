package logger

import (
	"log/slog"
)

var replaceKey = map[string]string{
	slog.MessageKey: "message",
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if key, ok := replaceKey[a.Key]; ok {
		a.Key = key
	}

	return a
}
