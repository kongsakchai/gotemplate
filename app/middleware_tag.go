package app

import (
	"context"

	"github.com/labstack/echo/v5"
)

func TagMiddleware(tagValue string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			req := ctx.Request()

			ctx.Set(TagKey, tagValue)
			ctx.SetRequest(req.WithContext(context.WithValue(req.Context(), TagKey, tagValue)))

			ctx.SetLogger(ctx.Logger().With(
				TagKey, tagValue,
			))

			return next(ctx)
		}
	}
}
