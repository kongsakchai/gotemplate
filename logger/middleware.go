package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.NewString()
		c.Set("traceID", traceID)

		body, err := c.GetRawData()
		if err != nil {
			c.String(500, "failed to read request body")
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path
		go func() {

			InfoCtx(c,
				fmt.Sprintf("request info %s", path),
				slog.String("traceId", traceID),
				slog.Group(
					"request",
					"method", method,
					"body", string(body),
				))
		}()

		c.Writer = &ginWriter{
			ResponseWriter: c.Writer,
			path:           path,
			method:         method,
			traceID:        traceID,
			start:          time.Now(),
		}

		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}

func EchoLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceID := uuid.NewString()
			c.Set("traceID", traceID)

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.String(500, "failed to read request body")
			}

			method := c.Request().Method
			path := c.Request().URL.Path
			go func() {
				InfoCtx(c.Request().Context(),
					fmt.Sprintf("request info %s", path),
					slog.String("traceId", traceID),
					slog.Group(
						"request",
						"method", method,
						"body", string(body),
					))
			}()

			c.Response().Writer = &echoWriter{
				ResponseWriter: c.Response().Writer,
				path:           path,
				method:         method,
				traceID:        traceID,
				start:          time.Now(),
			}

			c.Request().Body.Close()
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
			return next(c)
		}
	}
}
