package app

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type EchoApp struct {
	*echo.Echo

	server *http.Server
}

func NewEchoApp() *EchoApp {
	e := echo.New()
	e.Use(middleware.Recover())

	return &EchoApp{Echo: e}
}

func (app *EchoApp) Start(ctx context.Context, addr string) error {
	app.server = &http.Server{Addr: addr, Handler: app.Echo}
	return app.server.ListenAndServe()
}

func (app *EchoApp) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
