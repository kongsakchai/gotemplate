package app

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/labstack/echo/v5"
)

func LoggerMiddleware(enable bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			if !enable {
				return next(ctx)
			}

			req := ctx.Request()
			logger := ctx.Logger()

			responseWriter := ctx.Response()
			ctx.SetResponse(&echoResponseWriter{
				ResponseWriter: responseWriter,
				ctx:            ctx,
				reqTime:        time.Now(),
			})

			b, err := io.ReadAll(req.Body)
			if err != nil {
				logger.ErrorContext(req.Context(), "failed to read request body", "error", err)
				return err
			}

			logger.InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
				"method", req.Method,
				"body", string(b),
				"event", "api_request",
			)

			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(b))
			ctx.SetLogger(logger)

			return next(ctx)
		}
	}
}
