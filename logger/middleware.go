package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/kongsakchai/gotemplate/app"
)

func GinLogger() app.Handler {
	return func(c app.Context) error {
		traceID := c.Get("traceID").(string)
		ctx, ok := c.(*app.GinContext)
		if !ok {
			return c.Next()
		}

		body, err := ctx.GetRawData()
		if err != nil {
			return c.Error(app.InternalServerError(app.UnknownErrorCode, "failed to read request body", err))
		}

		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		go func() {

			InfoCtx(ctx,
				fmt.Sprintf("request info %s", path),
				slog.String("traceId", traceID),
				slog.Group(
					"request",
					"method", method,
					"body", string(body),
				))
		}()

		ctx.Writer = &ginWriter{
			ResponseWriter: ctx.Writer,
			path:           path,
			method:         method,
			traceID:        traceID,
			start:          time.Now(),
		}

		ctx.Request.Body.Close()
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		return c.Next()
	}
}

func EchoLogger() app.Handler {
	return func(c app.Context) error {
		traceID := c.Get("traceID").(string)
		ctx, ok := c.(*app.EchoContext)
		if !ok {
			return c.Next()
		}

		body, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return c.Error(app.InternalServerError(app.UnknownErrorCode, "failed to read request body", err))
		}

		method := ctx.Request().Method
		path := ctx.Request().URL.Path
		go func() {
			InfoCtx(ctx.Request().Context(),
				fmt.Sprintf("request info %s", path),
				slog.String("traceId", traceID),
				slog.Group(
					"request",
					"method", method,
					"body", string(body),
				))
		}()

		ctx.Response().Writer = &echoWriter{
			ResponseWriter: ctx.Response().Writer,
			path:           path,
			method:         method,
			traceID:        traceID,
			start:          time.Now(),
		}

		ctx.Request().Body.Close()
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

		return c.Next()
	}
}
