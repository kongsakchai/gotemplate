package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type echoResponseWriter struct {
	http.ResponseWriter
	logger *slog.Logger
	ctx    context.Context
	status int
	url    string
	now    time.Time
}

func (w *echoResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *echoResponseWriter) Write(b []byte) (int, error) {
	body := string(b)
	if w.Header().Get(echo.HeaderContentType) != echo.MIMEApplicationJSON {
		body = w.Header().Get(echo.HeaderContentType)
	}

	if w.status == http.StatusOK || w.status == http.StatusCreated {
		w.logger.InfoContext(w.ctx, fmt.Sprintf("response %d %s", w.status, w.url),
			"body", body,
			"latency", time.Since(w.now).String(),
		)
	} else {
		w.logger.ErrorContext(w.ctx, fmt.Sprintf("response %d %s", w.status, w.url),
			"body", body,
			"latency", time.Since(w.now).String(),
		)
	}

	return w.ResponseWriter.Write(b)
}
