package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RefID(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			refID := ctx.Request().Header.Get(key)
			req := ctx.Request()
			if refID == "" {
				slog.WarnContext(req.Context(), "no refID", slog.String("key", refID))
				refID = uuid.NewString()
			}

			ctx.SetRequest(req.WithContext(newRefIDContext(req.Context(), key, refID)))
			return next(ctx)
		}
	}
}

func newRefIDContext(ctx context.Context, key, refID string) context.Context {
	return context.WithValue(ctx, key, refID)
}
