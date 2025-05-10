package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kongsakchai/gotemplate/app"
)

type httpContext interface {
	Request() *http.Request
	Writer() http.ResponseWriter
	SetWriter(http.ResponseWriter)
}

func LoggerRequest() app.Middleware {
	return func(next app.Handler) app.Handler {
		return func(ctx app.Context) error {
			traceID := uuid.NewString()
			startTime := time.Now()
			ctx.Set("traceId", traceID)
			ctx.Set("startTime", startTime)

			httpCtx, ok := ctx.(httpContext)
			if !ok {
				return next(ctx)
			}

			req := httpCtx.Request()
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return ctx.InternalServer(app.Error{
					Code:    app.InternalServerCode,
					Message: err.Error(),
				})
			}

			go func() {
				ctx.Logger().InfoContext(
					req.Context(),
					fmt.Sprintf("request info %s", req.URL.Path),
					slog.Group(
						"request",
						"method", req.Method,
						"body", string(body),
					),
				)
			}()

			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			return next(ctx)
		}
	}
}

func LoggerResponse() app.Middleware {
	return func(next app.Handler) app.Handler {
		return func(ctx app.Context) error {
			httpCtx, ok := ctx.(httpContext)
			if !ok {
				return next(ctx)
			}

			traceID := ctx.Get("traceId").(string)
			startTime := ctx.Get("startTime").(time.Time)

			req := httpCtx.Request()
			meta := map[string]any{
				"method":  req.Method,
				"path":    req.URL.Path,
				"traceId": traceID,
				"latency": time.Since(startTime).String(),
			}

			httpCtx.SetWriter(&responseWriter{
				ResponseWriter: httpCtx.Writer(),
				meta:           meta,
			})

			return next(ctx)
		}
	}
}
