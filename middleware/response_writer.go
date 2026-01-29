package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
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
			"latency", time.Since(w.meta["now"].(time.Time)).String(),
		)
	} else {
		slog.Error(fmt.Sprintf("response %d %s", w.status, w.meta["url"]),
			"body", string(b),
			"traceID", w.meta["traceID"],
			"latency", time.Since(w.meta["now"].(time.Time)).String(),
		)
	}

	return w.ResponseWriter.Write(b)
}
