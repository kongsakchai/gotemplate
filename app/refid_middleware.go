package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RefIDMiddleware(key string, tags map[string]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			refID := ctx.Request().Header.Get(key)
			req := ctx.Request()
			if refID == "" {
				ctx.Logger().DebugContext(req.Context(), "no refID", slog.String("key", key))
				refID = uuid.NewString()
			}
			ctx.Set(TraceIDKey, refID)

			if tag, ok := tags[req.URL.Path]; ok {
				ctx.Set(TagKey, tag)
			} else {
				ctx.Set(TagKey, strings.ReplaceAll(req.URL.Path, "/", "-")[1:])
			}

			reqCtx := context.WithValue(req.Context(), key, refID)
			ctx.SetRequest(req.WithContext(context.WithValue(reqCtx, TraceIDKey, refID)))

			ctx.SetLogger(ctx.Logger().With(
				TraceIDKey, refID,
				TagKey, ctx.Get(TagKey),
			))

			return next(ctx)
		}
	}
}
