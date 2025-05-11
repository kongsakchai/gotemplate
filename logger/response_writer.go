package logger

import (
	"fmt"
	"log/slog"
)

type responseWriter struct {
	code int
	meta map[string]any
}

func (w *responseWriter) Write(b []byte) {
	go func() {
		Info(
			fmt.Sprintf("response info status=%d %s", w.code, w.meta["path"]),
			slog.String("traceId", w.meta["traceId"].(string)),
			slog.Group(
				"response",
				"method", w.meta["method"],
				"body", string(b),
				"status", w.code,
				"latency", w.meta["latency"],
			),
		)
	}()
}

func (w *responseWriter) WriteHeader(code int) {
	w.code = code
}
