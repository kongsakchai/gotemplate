package app

import (
	"context"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoApp struct {
	*echo.Echo
}

func NewEchoApp() *EchoApp {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())

	return &EchoApp{e}
}

func (app *EchoApp) Start(addr string) error {
	return app.Echo.Start(addr)
}

func (app *EchoApp) Shutdown(ctx context.Context) error {
	return app.Echo.Shutdown(ctx)
}

func NewMockContext(method, target, payload string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, target, strings.NewReader(payload))
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
