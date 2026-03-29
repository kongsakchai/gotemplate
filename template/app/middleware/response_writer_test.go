package middleware

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func setupResponseWriterLogger(w io.Writer) *slog.Logger {
	defaultLogger := slog.Default()
	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)

	return defaultLogger
}

func resetResponseWriterLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

func TestEchoResponseWriter_Write(t *testing.T) {
	t.Run("should log info and keep json body", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupResponseWriterLogger(buf)
		defer resetResponseWriterLogger(defaultLogger)

		rec := httptest.NewRecorder()
		w := &echoResponseWriter{
			ResponseWriter: rec,
			ctx:            context.Background(),
			url:            "/users",
			traceID:        "trace-1",
		}

		w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte(`{"name":"john"}`))

		assert.NoError(t, err)
		assert.Equal(t, len(`{"name":"john"}`), n)
		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "response 200 /users")
		assert.Contains(t, buf.String(), "body=\"{\\\"name\\\":\\\"john\\\"}\"")
		assert.Equal(t, `{"name":"john"}`, rec.Body.String())
	})

	t.Run("should log content-type for non json response", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupResponseWriterLogger(buf)
		defer resetResponseWriterLogger(defaultLogger)

		rec := httptest.NewRecorder()
		w := &echoResponseWriter{
			ResponseWriter: rec,
			ctx:            context.Background(),
			url:            "/ping",
			traceID:        "trace-2",
		}

		w.Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte("Pong!"))

		assert.NoError(t, err)
		assert.Equal(t, len("Pong!"), n)
		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "response 200 /ping")
		assert.Contains(t, buf.String(), "body=\"text/plain; charset=UTF-8\"")
		assert.Equal(t, "Pong!", rec.Body.String())
	})

	t.Run("should log error for non success status", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupResponseWriterLogger(buf)
		defer resetResponseWriterLogger(defaultLogger)

		rec := httptest.NewRecorder()
		w := &echoResponseWriter{
			ResponseWriter: rec,
			ctx:            context.Background(),
			url:            "/error",
			traceID:        "trace-3",
		}

		w.Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		n, err := w.Write([]byte("Internal Server Error"))

		assert.NoError(t, err)
		assert.Equal(t, len("Internal Server Error"), n)
		assert.Contains(t, buf.String(), "level=ERROR")
		assert.Contains(t, buf.String(), "response 500 /error")
		assert.Contains(t, buf.String(), "body=\"text/plain; charset=UTF-8\"")
		assert.Equal(t, "Internal Server Error", rec.Body.String())
	})
}
