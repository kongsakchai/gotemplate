package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestTagMiddleware(t *testing.T) {
	t.Run("should not set tag when route has no name", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		called := false

		middleware := TagMiddleware()
		err := middleware(func(c *echo.Context) error {
			called = true
			return nil
		})(ctx)

		assert.NoError(t, err)
		assert.True(t, called)
		assert.Nil(t, ctx.Get(TagKey))
		assert.Nil(t, ctx.Request().Context().Value(TagKey))
	})

	t.Run("should set tag from route name", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
			RouteInfo: &echo.RouteInfo{
				Name:       "test-route",
				Method:     http.MethodGet,
				Path:       "/test",
				Parameters: []string{},
			},
		}.ToContextRecorder(t)

		called := false

		middleware := TagMiddleware()
		err := middleware(func(c *echo.Context) error {
			called = true
			return nil
		})(ctx)

		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, "test-route", ctx.Get(TagKey))
		assert.Equal(t, "test-route", ctx.Request().Context().Value(TagKey))
	})
}
