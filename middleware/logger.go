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
		return func(c echo.Context) error {
			if req {
				b, err := io.ReadAll(c.Request().Body)
				if err != nil {
					return err
				}

				slog.Info(fmt.Sprintf("request %s", c.Request().URL),
					"method", c.Request().Method,
					"body", string(b),
					"traceID", c.Request().Context().Value(key).(string),
				)

				c.Request().Body.Close()
				c.Request().Body = io.NopCloser(bytes.NewBuffer(b))
			}

			if res {
				c.Response().Writer = &echoResponseWriter{
					ResponseWriter: c.Response().Writer,
					meta: map[string]any{
						"traceID": c.Request().Context().Value(key),
						"url":     c.Request().URL.String(),
					},
				}
			}

			return next(c)
		}
	}
}
