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
			req := ctx.Request()

			if enable {
				ctx.Response().Writer = &echoResponseWriter{
					ResponseWriter: ctx.Response().Writer,
					ctx:            req.Context(),
					url:            req.URL.String(),
					now:            time.Now(),
					traceID:        traceID,
				}
			}

			if enable {
				if req.Header.Get(echo.HeaderContentType) != echo.MIMEApplicationJSON {
					slog.InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
						"method", req.Method,
						"body", req.Header.Get(echo.HeaderContentType),
						app.TraceIDKey, traceID,
					)
					return next(ctx)
				}

				b, err := io.ReadAll(req.Body)
				if err != nil {
					slog.ErrorContext(req.Context(), "failed to read request body", "error", err)
					return err
				}

				slog.InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
					"method", req.Method,
					"body", string(b),
					app.TraceIDKey, traceID,
				)

				req.Body.Close()
				req.Body = io.NopCloser(bytes.NewBuffer(b))
			}

			return next(ctx)
		}
	}
}
