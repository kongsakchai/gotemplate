package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
)

func Logger(enable bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			traceID, _ := ctx.Get(app.TraceID).(string)
			tag, _ := ctx.Get(app.Tag).(string)
			req := ctx.Request()

			if !enable {
				return next(ctx)
			}

			responseWriter := ctx.Response()
			ctx.SetResponse(&echoResponseWriter{
				ResponseWriter: responseWriter,
				ctx:            req.Context(),
				url:            req.URL.String(),
				now:            time.Now(),
				traceID:        traceID,
				tag:            tag,
			})

			switch req.Header.Get(echo.HeaderContentType) {
			case echo.MIMEApplicationJSON, echo.MIMETextPlain:
				break
			default:
				ctx.Logger().InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
					"method", req.Method,
					"body", req.Header.Get(echo.HeaderContentType),
					app.TraceID, traceID,
					app.Tag, tag,
				)
				return next(ctx)
			}

			b, err := io.ReadAll(req.Body)
			if err != nil {
				ctx.Logger().ErrorContext(req.Context(), "failed to read request body",
					"error", err,
					app.TraceID, traceID,
					app.Tag, tag,
				)
				return err
			}

			ctx.Logger().InfoContext(req.Context(), fmt.Sprintf("request %s", req.URL),
				"method", req.Method,
				"body", string(b),
				app.TraceID, traceID,
				app.Tag, tag,
			)

			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(b))

			return next(ctx)
		}
	}
}
