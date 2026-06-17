package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestRefIDMiddleware(t *testing.T) {
	t.Run("should generate new refID when header not present", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/test", nil),
		}.ToContextRecorder(t)

		middleware := RefIDMiddleware("X-Ref-ID", nil)
		handler := middleware(func(ctx *echo.Context) error {
			traceID, _ := ctx.Get(TraceIDKey).(string)
			assert.NotEmpty(t, traceID)
			tag, _ := ctx.Get(TagKey).(string)
			assert.Equal(t, "api-v1-test", tag)
			return nil
		})

		err := handler(ctx)
		assert.NoError(t, err)
	})

	t.Run("should use refID from header when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
		req.Header.Set("X-Ref-ID", "custom-ref-id")
		ctx, _ := echotest.ContextConfig{
			Request: req,
		}.ToContextRecorder(t)

		middleware := RefIDMiddleware("X-Ref-ID", nil)
		handler := middleware(func(ctx *echo.Context) error {
			traceID, _ := ctx.Get(TraceIDKey).(string)
			assert.Equal(t, "custom-ref-id", traceID)
			return nil
		})

		err := handler(ctx)
		assert.NoError(t, err)
	})

	t.Run("should use tag from tags map when path matches", func(t *testing.T) {
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/members", nil),
		}.ToContextRecorder(t)

		tags := map[string]string{
			"/api/v1/members": "member-tag",
		}

		middleware := RefIDMiddleware("X-Ref-ID", tags)
		handler := middleware(func(ctx *echo.Context) error {
			tag, _ := ctx.Get(TagKey).(string)
			assert.Equal(t, "member-tag", tag)
			return nil
		})

		err := handler(ctx)
		assert.NoError(t, err)
	})
}
