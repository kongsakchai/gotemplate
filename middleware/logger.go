package middleware

import (
	"log/slog"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

func Logger(key string, req, res bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if req {
				body, err := httputil.DumpRequest(c.Request(), true)
				if err != nil {
					return err
				}

				slog.Info("request",
					slog.String("traceID", c.Request().Context().Value(key).(string)),
					slog.String("method", c.Request().Method),
					slog.String("url", c.Request().URL.String()),
					slog.String("body", string(body)),
				)
			}

			if res {
				c.Response().Writer = &echoResponseWriter{
					ResponseWriter: c.Response().Writer,
					meta: map[string]any{
						"traceID": c.Request().Context().Value(key),
						"method":  c.Request().Method,
						"url":     c.Request().URL.String(),
					},
				}
			}

			return next(c)
		}
	}
}
