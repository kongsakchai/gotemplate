package app

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kongsakchai/gotemplate/template/pkg/config"
	"github.com/kongsakchai/gotemplate/template/pkg/errs"
	"github.com/kongsakchai/gotemplate/template/pkg/validator"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type EchoApp struct {
	*echo.Echo
}

func NewEchoApp(cfg config.Config) *EchoApp {
	e := echo.New()
	e.Validator = validator.NewReqValidator()
	e.HTTPErrorHandler = errorHandler

	e.Use(
		middleware.Recover(),
		middleware.CORS("*"),
		RefIDMiddleware(cfg.Header.RefIDKey, cfg.Log.Tags),
		LoggerMiddleware(cfg.Log.Enable),
	)

	return &EchoApp{Echo: e}
}

func (app *EchoApp) Start(ctx context.Context, addr string, gracefulTimeout time.Duration) error {
	for _, r := range app.Router().Routes() {
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
		ctx.Logger().LogAttrs(ctx.Request().Context(), slog.LevelError, "app error", errs.SlogAttr(appErr.Err)...)
		if err := Fail(ctx, appErr); err != nil {
			ctx.Logger().ErrorContext(ctx.Request().Context(), "error handler fail", "err", err.Error()) // rare case
		}
		return
	}

	defaultEchoErrorHandler(ctx, err)
}

// reference from echo.DefaultHTTPErrorHandler but convert to app.Error
func defaultEchoErrorHandler(ctx *echo.Context, err error) {
	ctx.Logger().LogAttrs(ctx.Request().Context(), slog.LevelError, "unhandle error", errs.SlogAttr(err)...)

	appErr := Error{HTTPCode: http.StatusInternalServerError}
	var sc echo.HTTPStatusCoder

	if errors.As(err, &sc) {
		if tmp := sc.StatusCode(); tmp != 0 {
			appErr.HTTPCode = tmp
		}
	}

	switch m := sc.(type) {
	case json.Marshaler: // this type knows how to format itself to JSON
		b, _ := m.MarshalJSON()
		appErr.Message = string(b)
	case *echo.HTTPError:
		appErr.Message = m.Message
		if appErr.Message == "" {
			appErr.Message = http.StatusText(appErr.HTTPCode)
		}
	default:
		appErr.Message = http.StatusText(appErr.HTTPCode)
	}

	if err := Fail(ctx, appErr); err != nil {
		ctx.Logger().ErrorContext(ctx.Request().Context(), "error handler fail", "err", err.Error()) // rare case
	}
}
