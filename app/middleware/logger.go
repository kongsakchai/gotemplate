package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v4"
)

func Logger(enable bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			traceID, _ := ctx.Get(app.RefIDKey).(string)
			req := ctx.Request().Context()

			if enable {
				b, err := io.ReadAll(ctx.Request().Body)
				if err != nil {
					slog.ErrorContext(req, "failed to read request body", "error", err)
					return err
				}

				slog.InfoContext(req, fmt.Sprintf("request %s", ctx.Request().URL),
					"method", ctx.Request().Method,
					"body", string(b),
					app.TraceIDKey, traceID,
				)

				ctx.Request().Body.Close()
				ctx.Request().Body = io.NopCloser(bytes.NewBuffer(b))
			}

			if enable {
				ctx.Response().Writer = &echoResponseWriter{
					ctx:            ctx.Request().Context(),
					ResponseWriter: ctx.Response().Writer,
					meta: map[string]any{
						"url":          ctx.Request().URL.String(),
						"now":          time.Now(),
						app.TraceIDKey: traceID,
					},
				}
			}

			return next(ctx)
		}
	}
}
