package app

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RefIDMiddleware(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			refID := ctx.Request().Header.Get(key)
			req := ctx.Request()
			if refID == "" {
				ctx.Logger().DebugContext(req.Context(), "no refID", slog.String("key", key))
				refID = uuid.NewString()
			}

			ctx.Set(TraceIDKey, refID)
			ctx.SetRequest(req.WithContext(context.WithValue(req.Context(), TraceIDKey, refID)))
			ctx.SetLogger(ctx.Logger().With(TraceIDKey, refID))

			return next(ctx)
		}
	}
}
