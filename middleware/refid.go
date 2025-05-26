package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func EchoRefID(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			refID := c.Request().Header.Get(key)
			req := c.Request()
			if refID != "" {
				slog.WarnContext(req.Context(), "no refID", slog.String("key", refID))
				refID = uuid.NewString()
			}

			c.SetRequest(req.WithContext(newRefIDContext(req.Context(), key, refID)))
			return next(c)
		}
	}
}

func newRefIDContext(ctx context.Context, key, refID string) context.Context {
	return context.WithValue(ctx, key, refID)
}
