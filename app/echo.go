package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type echoContext struct {
	echo.Context
}

func (e *echoContext) Param(key string) string {
	return e.Context.Param(key)
}

func (e *echoContext) Bind(obj any) error {
	return e.Context.Bind(obj)
}

func (e *echoContext) OK(obj any) error {
	return e.Context.JSON(200, response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) Created(obj any) error {
	return e.Context.JSON(201, response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) NotFound(err Error) error {
	return e.Context.JSON(404, response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) InternalServer(err Error) error {
	return e.Context.JSON(500, response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) BadRequest(err Error) error {
	return e.Context.JSON(400, response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

type echoRouter struct {
	*echo.Echo
}

func NewEcho(logLevel string) *echoRouter {
	e := echo.New()

	switch logLevel {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "warn":
		e.Logger.SetLevel(log.WARN)
	case "error":
		e.Logger.SetLevel(log.ERROR)
	}

	return &echoRouter{e}
}

func newEchoHandler(handler func(c Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(&echoContext{c})
	}
}

func (e *echoRouter) GET(path string, handler func(c Context) error) {
	e.Echo.GET(path, newEchoHandler(handler))
}

func (e *echoRouter) POST(path string, handler func(c Context) error) {
	e.Echo.POST(path, newEchoHandler(handler))
}

func (e *echoRouter) PUT(path string, handler func(c Context) error) {
	e.Echo.PUT(path, newEchoHandler(handler))
}

func (e *echoRouter) DELETE(path string, handler func(c Context) error) {
	e.Echo.DELETE(path, newEchoHandler(handler))
}

func (e *echoRouter) PATCH(path string, handler func(c Context) error) {
	e.Echo.PATCH(path, newEchoHandler(handler))
}
