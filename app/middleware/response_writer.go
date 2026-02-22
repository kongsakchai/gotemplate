package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kongsakchai/gotemplate/app"
)

type echoResponseWriter struct {
	http.ResponseWriter
	ctx    context.Context
	status int
	meta   map[string]any
}

func (w *echoResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *echoResponseWriter) Write(b []byte) (int, error) {
	if w.status == http.StatusOK || w.status == http.StatusCreated {
		slog.InfoContext(w.ctx, fmt.Sprintf("response %d %s", w.status, w.meta["url"]),
			"body", string(b),
			"latency", time.Since(w.meta["now"].(time.Time)).String(),
			app.TraceIDKey, w.meta[app.TraceIDKey],
		)
	} else {
		slog.ErrorContext(w.ctx, fmt.Sprintf("response %d %s", w.status, w.meta["url"]),
			"body", string(b),
			"latency", time.Since(w.meta["now"].(time.Time)).String(),
			app.TraceIDKey, w.meta[app.TraceIDKey],
		)
	}

	return w.ResponseWriter.Write(b)
}
