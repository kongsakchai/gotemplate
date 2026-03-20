package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
)

func RefID(key string, tags map[string]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			refID := ctx.Request().Header.Get(key)
			req := ctx.Request()
			if refID == "" {
				ctx.Logger().DebugContext(req.Context(), "no refID", slog.String("key", key))
				refID = uuid.NewString()
			}
			ctx.Set(app.TraceID, refID)

			if tag, ok := tags[req.URL.Path]; ok {
				ctx.Set(app.Tag, tag)
			} else {
				ctx.Set(app.Tag, strings.ReplaceAll(req.URL.Path, "/", "-")[1:])
			}

			reqCtx := context.WithValue(req.Context(), key, refID)
			ctx.SetRequest(req.WithContext(context.WithValue(reqCtx, app.TraceID, refID)))

			ctx.SetLogger(ctx.Logger().With(
				app.TraceID, refID,
				app.Tag, ctx.Get(app.Tag),
			))

			return next(ctx)
		}
	}
}
