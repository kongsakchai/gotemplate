package app

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func LoggerMiddleware(enable bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			traceID, _ := ctx.Get(TraceIDKey).(string)
			tag, _ := ctx.Get(TagKey).(string)
			req := ctx.Request()

			if !enable {
				return next(ctx)
			}

			logger := ctx.Logger().With(
				slog.String(TraceIDKey, traceID),
				slog.String(TagKey, tag),
			)

			responseWriter := ctx.Response()
			ctx.SetResponse(&echoResponseWriter{
				ResponseWriter: responseWriter,
				logger:         logger,
				ctx:            req.Context(),
				url:            req.URL.String(),
				now:            time.Now(),
			})

			b, err := io.ReadAll(req.Body)
			if err != nil {
				logger.ErrorContext(req.Context(), "failed to read request body",
					"error", err,
					TraceIDKey, traceID,
					TagKey, tag,
				)
				return err
			}

			logger.InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
				"method", req.Method,
				"body", string(b),
			)

			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(b))
			ctx.SetLogger(logger)

			return next(ctx)
		}
	}
}
