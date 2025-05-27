package app

import (
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
