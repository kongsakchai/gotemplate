package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v4"
)

func RefID(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			refID := ctx.Request().Header.Get(key)
			req := ctx.Request()
			if refID == "" {
				slog.DebugContext(req.Context(), "no refID", slog.String("key", key))
				refID = uuid.NewString()
			}

			ctx.Set(app.RefIDKey, refID)
			newCtx := context.WithValue(req.Context(), key, refID)
			ctx.SetRequest(req.WithContext(newCtx))
			return next(ctx)
		}
	}
}
