package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ginWriter struct {
	gin.ResponseWriter
	code    int
	path    string
	method  string
	traceID string
	start   time.Time
}

func (w *ginWriter) Write(b []byte) (int, error) {
	go func() {
		Info(
			fmt.Sprintf("response info status=%d %s", w.code, w.path),
			slog.String("traceId", w.traceID),
			slog.Group(
				"response",
				"method", w.method,
				"body", string(b),
				"status", w.code,
				"latency", time.Since(w.start).String(),
			),
		)
	}()

	return w.ResponseWriter.Write(b)
}

func (w *ginWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

type echoWriter struct {
	http.ResponseWriter
	code    int
	path    string
	method  string
	traceID string
	start   time.Time
}

func (w *echoWriter) Write(b []byte) (int, error) {
	go func() {
		Info(
			fmt.Sprintf("response info status=%d %s", w.code, w.path),
			slog.String("traceId", w.traceID),
			slog.Group(
				"response",
				"method", w.method,
				"body", string(b),
				"status", w.code,
				"latency", time.Since(w.start).String(),
			),
		)
	}()

	return w.ResponseWriter.Write(b)
}

func (w *echoWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
