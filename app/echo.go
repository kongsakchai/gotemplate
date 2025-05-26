package app

import "github.com/labstack/echo/v4"

type EchoApp struct {
	*echo.Echo
}

func NewEchoApp() *EchoApp {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &EchoApp{e}
}
