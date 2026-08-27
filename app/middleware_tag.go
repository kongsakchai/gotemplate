package app

import (
	"context"

	"github.com/labstack/echo/v5"
)

func TagMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			req := ctx.Request()

			name := ctx.RouteInfo().Name
			ctx.Set(TagKey, name)
			ctx.SetRequest(req.WithContext(context.WithValue(req.Context(), TagKey, name)))
			ctx.SetLogger(ctx.Logger().With(
				TagKey, name,
			))

			return next(ctx)
		}
	}
}
