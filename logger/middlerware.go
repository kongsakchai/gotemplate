package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kongsakchai/gotemplate/app"
)

type httpContext interface {
	Request() *http.Request
	Writer() http.ResponseWriter
	SetWriter(http.ResponseWriter)
}

func LoggerRequest() app.Handler {
	return func(ctx app.Context) error {
		ctx.Set("request_id", uuid.NewString())

		httpCtx, ok := ctx.(httpContext)
		if !ok {
			return ctx.Next(ctx)
		}

		req := httpCtx.Request()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return ctx.InternalServer(app.Error{
				Code:    app.InternalServerCode,
				Message: err.Error(),
			})
		}
		req.Body.Close()

		go func() {
			InfoCtx(
				req.Context(),
				fmt.Sprintf("request info %s", req.URL.Path),
				slog.Group(
					"request",
					"id", ctx.Get("request_id"),
					"method", req.Method,
					"body", string(body),
				),
			)
		}()

		req.Body = io.NopCloser(bytes.NewBuffer(body))
		return ctx.Next(ctx)
	}
}

func LoggerResponse() app.Handler {
	return func(ctx app.Context) error {
		httpCtx, ok := ctx.(httpContext)
		if !ok {
			return ctx.Next(ctx)
		}

		req := httpCtx.Request()
		meta := map[string]string{
			"request_id": ctx.Get("request_id").(string),
			"method":     req.Method,
			"path":       req.URL.Path,
		}

		httpCtx.SetWriter(&responseWriter{
			ResponseWriter: httpCtx.Writer(),
			meta:           meta,
		})

		return ctx.Next(ctx)
	}
}
