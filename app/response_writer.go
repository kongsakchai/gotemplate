package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type echoResponseWriter struct {
	http.ResponseWriter
	ctx     *echo.Context
	status  int
	reqTime time.Time
}

func (w *echoResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *echoResponseWriter) Write(b []byte) (int, error) {
	body := string(b)
	ctx := w.ctx.Request().Context()
	url := w.ctx.Request().URL.String()

	if w.status == http.StatusOK || w.status == http.StatusCreated {
		w.ctx.Logger().InfoContext(ctx, fmt.Sprintf("api response %d %s", w.status, url),
			"body", body,
			"latency", time.Since(w.reqTime).String(),
			"event", "api_response",
		)
	} else {
		w.ctx.Logger().ErrorContext(ctx, fmt.Sprintf("api response %d %s", w.status, url),
			"body", body,
			"latency", time.Since(w.reqTime).String(),
			"event", "api_response",
		)
	}

	return w.ResponseWriter.Write(b)
}
