package middleware

import (
	"log/slog"
	"net/http"
)

type echoResponseWriter struct {
	http.ResponseWriter
	status int
	meta   map[string]any
}

func (w *echoResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *echoResponseWriter) Write(b []byte) (int, error) {
	slog.Info("response",
		slog.Int("status", w.status),
		slog.Int("length", len(b)),
		slog.Any("method", w.meta["method"]),
		slog.Any("url", w.meta["url"]),
		slog.Any("traceID", w.meta["traceID"]),
	)

	return w.ResponseWriter.Write(b)
}
