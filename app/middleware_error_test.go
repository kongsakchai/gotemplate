package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestErrorMiddleware(t *testing.T) {
	t.Run("should pass through nil error", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		called := false

		middleware := ErrorMiddleware(func(err error) Error {
			called = true
			return InternalError(InternalErrorCode, InternalErrorMsg, err)
		})

		err := middleware(func(c *echo.Context) error {
			return nil
		})(ctx)

		assert.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("should pass through app.Error without calling handler", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		target := BadRequest("4000", "bad request", errors.New("test error"))
		called := false

		middleware := ErrorMiddleware(func(err error) Error {
			called = true
			return InternalError(InternalErrorCode, InternalErrorMsg, err)
		})

		err := middleware(func(c *echo.Context) error {
			return target
		})(ctx)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, target, appErr)
		assert.False(t, called)
	})

	t.Run("should map plain error via handler", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		expected := InternalError(InternalErrorCode, InternalErrorMsg, assert.AnError)
		var got error

		middleware := ErrorMiddleware(func(err error) Error {
			got = err
			return expected
		})

		err := middleware(func(c *echo.Context) error {
			return assert.AnError
		})(ctx)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, assert.AnError, got)
		assert.Equal(t, expected, appErr)
	})
}
