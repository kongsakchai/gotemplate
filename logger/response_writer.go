package logger

import (
	"fmt"
	"log/slog"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	code int
	meta map[string]string
}

func (w *responseWriter) Write(b []byte) (int, error) {
	go func() {
		Info(
			fmt.Sprintf("response info %s", w.meta["path"]),
			slog.Group(
				"response",
				"id", w.meta["request_id"],
				"method", w.meta["method"],
				"body", string(b),
				"status", w.code,
			),
		)
	}()

	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
