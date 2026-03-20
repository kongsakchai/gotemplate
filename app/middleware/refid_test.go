package middleware

import (
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/require"
)

func TestRefID(t *testing.T) {
	t.Run("should return ctx with refID and tag", func(t *testing.T) {
		// arrange
		next := func(ctx *echo.Context) error {
			require.NotEmpty(t, ctx.Get(app.TraceID))
			require.Equal(t, "test-tag", ctx.Get(app.Tag))
			return nil
		}
		ctx := echotest.ContextConfig{
			Headers: http.Header{
				"X-Ref-ID": []string{"test-ref-id"},
			},
		}.ToContext(t)

		// act
		middleware := RefID("X-Ref-ID", map[string]string{"/": "test-tag"})
		err := middleware(next)(ctx)

		// assert
		require.NoError(t, err)
	})

	t.Run("should return ctx with refID without tag", func(t *testing.T) {
		// arrange
		next := func(ctx *echo.Context) error {
			require.NotEmpty(t, ctx.Get(app.TraceID))
			require.Equal(t, "test-tag", ctx.Get(app.Tag))
			return nil
		}
		ctx := echotest.ContextConfig{
			Headers: http.Header{
				"X-Ref-ID": []string{"test-ref-id"},
			},
		}.ToContext(t)
		ctx.Request().URL.Path = "/test/tag"

		// act
		middleware := RefID("X-Ref-ID", nil)
		err := middleware(next)(ctx)

		// assert
		require.NoError(t, err)
	})
}
