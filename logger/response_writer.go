package logger

import (
	"fmt"
	"log/slog"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	code int
	meta map[string]any
}

func (w *responseWriter) Write(b []byte) (int, error) {
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

	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
