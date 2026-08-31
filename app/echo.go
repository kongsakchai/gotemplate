package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kongsakchai/gotemplate/pkg/config"
	"github.com/kongsakchai/gotemplate/pkg/logger"
	"github.com/kongsakchai/gotemplate/pkg/validator"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type EchoApp struct {
	*echo.Echo
}

func NewEchoApp(cfg config.Config) *EchoApp {
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = errorHandler

	e.Use(
		middleware.Recover(),
		middleware.CORS("*"),
		TagMiddleware(),
		RefIDMiddleware(cfg.Header.RefIDKey),
		LoggerMiddleware(cfg.Log.Enable),
	)

	return &EchoApp{Echo: e}
}

func (app *EchoApp) Start(ctx context.Context, addr string, gracefulTimeout time.Duration) error {
	for _, r := range app.Router().Routes() {
		if r.Method == echo.RouteNotFound {
			continue
		}
		slog.DebugContext(ctx, r.Method, "path", r.Path)
	}

	sc := echo.StartConfig{
		Address:         addr,
		GracefulTimeout: gracefulTimeout,
		HidePort:        true,
		HideBanner:      true,
	}
	return sc.Start(ctx, app)
}

func errorHandler(ctx *echo.Context, err error) {
	if appErr, ok := err.(Error); ok {
		ctx.Logger().LogAttrs(ctx.Request().Context(), slog.LevelError, "app error", logger.ErrorAttrs(appErr.Err)...)
		if err := Fail(ctx, appErr); err != nil {
			ctx.Logger().ErrorContext(ctx.Request().Context(), "error handler fail", "err", err.Error()) // rare case
		}
		return
	}

	defaultEchoErrorHandler(ctx, err)
}

// reference from echo.DefaultHTTPErrorHandler but convert to app.Error
func defaultEchoErrorHandler(ctx *echo.Context, err error) {
	ctx.Logger().LogAttrs(ctx.Request().Context(), slog.LevelError, "unhandle error", logger.ErrorAttrs(err)...)

	appErr := Error{Code: InternalErrorCode, HTTPCode: http.StatusInternalServerError}

	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		if tmp := sc.StatusCode(); tmp != 0 {
			appErr.Code = fmt.Sprintf("HTTP_%d", tmp)
			appErr.HTTPCode = tmp
		}
	}

	switch m := sc.(type) {
	case json.Marshaler: // this type knows how to format itself to JSON
		b, _ := m.MarshalJSON()
		appErr.Message = string(b)
	case *echo.HTTPError:
		appErr.Message = m.Message
	}

	if appErr.Message == "" {
		appErr.Message = http.StatusText(appErr.HTTPCode)
	}

	if err := Fail(ctx, appErr); err != nil {
		ctx.Logger().ErrorContext(ctx.Request().Context(), "error handler fail", "err", err.Error()) // rare case
	}
}
