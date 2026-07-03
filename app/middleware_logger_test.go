package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

type errReader struct {
	*strings.Reader
}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("should skip logging when disabled", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		middleware := LoggerMiddleware(false)
		var called bool
		handler := middleware(func(ctx *echo.Context) error {
			called = true
			return nil
		})

		err := handler(ctx)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("should log request when enabled", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		middleware := LoggerMiddleware(true)
		var called bool
		handler := middleware(func(ctx *echo.Context) error {
			called = true
			return nil
		})

		err := handler(ctx)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("should handle body read error when enabled", func(t *testing.T) {
		body := strings.NewReader("test")
		req := httptest.NewRequest(http.MethodPost, "/test", &errReader{body})
		req.ContentLength = 4
		ctx, _ := echotest.ContextConfig{
			Request: req,
		}.ToContextRecorder(t)

		middleware := LoggerMiddleware(true)
		var called bool
		handler := middleware(func(ctx *echo.Context) error {
			called = true
			return nil
		})

		err := handler(ctx)
		assert.Error(t, err)
		assert.False(t, called)
	})
}
