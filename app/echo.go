package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type echoContext struct {
	logger *slog.Logger
	echo.Context
}

var echoContextPool = &sync.Pool{
	New: func() any {
		return &echoContext{}
	},
}

func newEchoContext(logger *slog.Logger, ctx echo.Context) *echoContext {
	c := echoContextPool.Get().(*echoContext)
	c.reset(ctx, logger)
	return c
}

func putEchoContext(ctx *echoContext) {
	echoContextPool.Put(ctx)
}

func (e *echoContext) reset(ctx echo.Context, logger *slog.Logger) {
	e.Context = ctx
	e.logger = logger.With("traceID", ctx.Get("traceID"))
}

func (e *echoContext) Query(key string) string {
	return e.Context.QueryParam(key)
}

func (e *echoContext) Param(key string) string {
	return e.Context.Param(key)
}

func (e *echoContext) Bind(obj any) error {
	return e.Context.Bind(obj)
}

func (e *echoContext) JSON(code int, obj any) error {
	return e.Context.JSON(code, obj)
}

func (e *echoContext) OK(obj any) error {
	return e.JSON(200, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) OKWithMessage(message string, obj any) error {
	return e.JSON(200, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
	})
}

func (e *echoContext) Created(obj any) error {
	return e.JSON(201, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) CreatedWithMessage(message string, obj any) error {
	return e.JSON(201, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
	})
}

func (e *echoContext) Error(err *Error) error {
	e.logger.Error(err.Error())
	return e.JSON(err.StatusCd, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) Ctx() context.Context {
	return e.Context.Request().Context()
}

func (e *echoContext) Get(key string) any {
	return e.Context.Get(key)
}

func (e *echoContext) Set(key string, value any) {
	e.Context.Set(key, value)
}

func (e *echoContext) Logger() *slog.Logger {
	return e.logger
}

func newEchoHandler(handler Handler, logger *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := newEchoContext(logger, ctx)
		defer putEchoContext(c)
		return handler(c)
	}
}

type echoRouter struct {
	*echo.Echo
	logger *slog.Logger
}

func NewEchoRoute(logger *slog.Logger) *echoRouter {
	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	return &echoRouter{
		Echo:   e,
		logger: logger,
	}
}

func (e *echoRouter) GET(path string, handler Handler, m ...echo.MiddlewareFunc) {
	e.Echo.GET(path, newEchoHandler(handler, e.logger), m...)
}

func (e *echoRouter) POST(path string, handler Handler, m ...echo.MiddlewareFunc) {
	e.Echo.POST(path, newEchoHandler(handler, e.logger), m...)
}

func (e *echoRouter) PUT(path string, handler Handler, m ...echo.MiddlewareFunc) {
	e.Echo.PUT(path, newEchoHandler(handler, e.logger), m...)
}

func (e *echoRouter) DELETE(path string, handler Handler, m ...echo.MiddlewareFunc) {
	e.Echo.DELETE(path, newEchoHandler(handler, e.logger), m...)
}

func (e *echoRouter) PATCH(path string, handler Handler, m ...echo.MiddlewareFunc) {
	e.Echo.PATCH(path, newEchoHandler(handler, e.logger), m...)
}

type echoGroup struct {
	EchoGroup *echo.Group
	logger    *slog.Logger
}

func (e *echoRouter) Group(prefix string, m ...echo.MiddlewareFunc) *echoGroup {
	grp := e.Echo.Group(prefix, m...)
	return &echoGroup{
		EchoGroup: grp,
		logger:    e.logger,
	}
}

func (g *echoGroup) GET(path string, handler Handler, m ...echo.MiddlewareFunc) {
	g.EchoGroup.GET(path, newEchoHandler(handler, g.logger), m...)
}

func (g *echoGroup) POST(path string, handler Handler, m ...echo.MiddlewareFunc) {
	g.EchoGroup.POST(path, newEchoHandler(handler, g.logger), m...)
}

func (g *echoGroup) PUT(path string, handler Handler, m ...echo.MiddlewareFunc) {
	g.EchoGroup.PUT(path, newEchoHandler(handler, g.logger), m...)
}

func (g *echoGroup) DELETE(path string, handler Handler, m ...echo.MiddlewareFunc) {
	g.EchoGroup.DELETE(path, newEchoHandler(handler, g.logger), m...)
}

func (g *echoGroup) PATCH(path string, handler Handler, m ...echo.MiddlewareFunc) {
	g.EchoGroup.PATCH(path, newEchoHandler(handler, g.logger), m...)
}

func (g *echoGroup) Group(prefix string, m ...echo.MiddlewareFunc) *echoGroup {
	grp := g.EchoGroup.Group(prefix, m...)
	return &echoGroup{
		EchoGroup: grp,
		logger:    g.logger,
	}
}
