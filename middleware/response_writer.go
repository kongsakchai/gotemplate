package middleware

import (
	"fmt"
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
	if w.status == http.StatusOK || w.status == http.StatusCreated {
		slog.Info(fmt.Sprintf("response %d %s", w.status, w.meta["url"]),
			"body", string(b),
			"traceID", w.meta["traceID"],
		)
	} else {
		slog.Error(fmt.Sprintf("response %d %s", w.status, w.meta["url"]),
			"body", string(b),
			"traceID", w.meta["traceID"],
		)
	}

	return w.ResponseWriter.Write(b)
}
