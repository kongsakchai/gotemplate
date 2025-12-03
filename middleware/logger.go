package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func Logger(key string, req, res bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			traceID, _ := ctx.Get(key).(string)

			if req {
				b, err := io.ReadAll(ctx.Request().Body)
				if err != nil {
					slog.Error("failed to read request body", "error", err)
					return err
				}

				slog.Info(fmt.Sprintf("request %s", ctx.Request().URL),
					"method", ctx.Request().Method,
					"body", string(b),
					"trace_id", traceID,
				)

				ctx.Request().Body.Close()
				ctx.Request().Body = io.NopCloser(bytes.NewBuffer(b))
			}

			if res {
				ctx.Response().Writer = &echoResponseWriter{
					ResponseWriter: ctx.Response().Writer,
					meta: map[string]any{
						"trace_id": traceID,
						"url":      ctx.Request().URL.String(),
					},
				}
			}

			return next(ctx)
		}
	}
}
