package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetString("traceID")
		body, err := ctx.GetRawData()
		if err != nil {
			ctx.Next()
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

		ctx.Next()
	}
}

func EchoLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			traceID := ctx.Request().Header.Get("traceID")
			body, err := io.ReadAll(ctx.Request().Body)
			if err != nil {
				return next(ctx)
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

			return next(ctx)
		}
	}
}
