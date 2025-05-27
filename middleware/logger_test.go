package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func handler(c echo.Context) error {
	return c.String(http.StatusOK, "Pong!")
}

func handlerError(c echo.Context) error {
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}

func setupLogger(w io.Writer) *slog.Logger {
	defaultLogger := slog.Default()
	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)

	return defaultLogger
}

func resetLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

type mockReader struct {
	err error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("should log request", func(t *testing.T) {
		e := echo.New()
		e.Use(RefID("X-Request-ID"))
		e.Use(Logger("X-Request-ID", true, false))
		e.GET("/", handler)

		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupLogger(buf)
		defer resetLogger(defaultLogger)

		req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte("Ping!")))
		req.Header.Set("X-Request-ID", "12345")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "traceID=12345")
		assert.Contains(t, buf.String(), "body=Ping!")
	})

	t.Run("should return error when request body cannot be read", func(t *testing.T) {
		e := echo.New()
		e.Use(RefID("X-Request-ID"))
		e.Use(Logger("X-Request-ID", true, false))
		e.GET("/", handler)

		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupLogger(buf)
		defer resetLogger(defaultLogger)

		req := httptest.NewRequest(http.MethodGet, "/", &mockReader{err: io.ErrUnexpectedEOF})
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, buf.String(), "level=ERROR")
		assert.Contains(t, buf.String(), "failed to read request body")
	})

	t.Run("should log response", func(t *testing.T) {
		e := echo.New()
		e.Use(RefID("X-Request-ID"))
		e.Use(Logger("X-Request-ID", false, true))
		e.GET("/", handler)

		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupLogger(buf)
		defer resetLogger(defaultLogger)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "body=Pong!")
	})

	t.Run("should log error response", func(t *testing.T) {
		e := echo.New()
		e.Use(RefID("X-Request-ID"))
		e.Use(Logger("X-Request-ID", false, true))
		e.GET("/error", handlerError)

		buf := bytes.NewBuffer([]byte{})
		defaultLogger := setupLogger(buf)
		defer resetLogger(defaultLogger)

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, buf.String(), "level=ERROR")
		assert.Contains(t, buf.String(), "body=\"Internal Server Error\"")
	})
}
