package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type echoContext struct {
	logger *slog.Logger
	next   echo.HandlerFunc
	echo.Context
}

func newEchoContext(logger *slog.Logger, ctx echo.Context, next echo.HandlerFunc) *echoContext {
	if reuse, ok := ctx.(*echoContext); ok {
		reuse.next = next
		return reuse
	}

	return &echoContext{
		next:    next,
		Context: ctx,
		logger:  logger,
	}
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

func (e *echoContext) OK(obj any) error {
	return e.Context.JSON(200, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) Created(obj any) error {
	return e.Context.JSON(201, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *echoContext) NotFound(err Error) error {
	e.logger.Error(err.Error())
	return e.Context.JSON(404, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) InternalServer(err Error) error {
	e.logger.Error(err.Error())
	return e.Context.JSON(500, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) BadRequest(err Error) error {
	e.logger.Error(err.Error())
	return e.Context.JSON(400, Response{
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

func (e *echoContext) Next(ctx Context) error {
	if e.next == nil {
		return nil
	}

	return e.next(ctx.(*echoContext))
}

func (e *echoContext) Request() *http.Request {
	return e.Context.Request()
}

func (e *echoContext) Writer() http.ResponseWriter {
	return e.Context.Response().Writer
}

func (e *echoContext) SetWriter(w http.ResponseWriter) {
	e.Context.Response().Writer = w
}

func (e *echoContext) CtxLogger() *slog.Logger {
	return e.logger
}

func (e *echoContext) SetCtxLogger(logger *slog.Logger) {
	e.logger = logger
}

func newEchoHandler(handler Handler, logger *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return handler(newEchoContext(logger, ctx, nil))
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

	return &echoRouter{e, logger}
}

func (e *echoRouter) GET(path string, handler Handler) {
	e.Echo.GET(path, newEchoHandler(handler, e.logger))
}

func (e *echoRouter) POST(path string, handler Handler) {
	e.Echo.POST(path, newEchoHandler(handler, e.logger))
}

func (e *echoRouter) PUT(path string, handler Handler) {
	e.Echo.PUT(path, newEchoHandler(handler, e.logger))
}

func (e *echoRouter) DELETE(path string, handler Handler) {
	e.Echo.DELETE(path, newEchoHandler(handler, e.logger))
}

func (e *echoRouter) PATCH(path string, handler Handler) {
	e.Echo.PATCH(path, newEchoHandler(handler, e.logger))
}

func newEchoMiddleware(logger *slog.Logger, handler Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			return handler(newEchoContext(logger, ctx, next))
		}
	}
}

func newEchoMiddlewares(logger *slog.Logger, handler ...Handler) []echo.MiddlewareFunc {
	m := make([]echo.MiddlewareFunc, len(handler))
	for i, h := range handler {
		m[i] = newEchoMiddleware(logger, h)
	}
	return m
}

func (e *echoRouter) Use(m ...Handler) {
	e.Echo.Use(newEchoMiddlewares(e.logger, m...)...)
}

type echoGroup struct {
	EchoGroup *echo.Group
	logger    *slog.Logger
}

func (e *echoRouter) Group(prefix string, m ...echo.MiddlewareFunc) *echoGroup {
	g := e.Echo.Group(prefix, m...)
	return &echoGroup{g, e.logger}
}

func (g *echoGroup) GET(path string, handler Handler) {
	g.EchoGroup.GET(path, newEchoHandler(handler, g.logger))
}

func (g *echoGroup) POST(path string, handler Handler) {
	g.EchoGroup.POST(path, newEchoHandler(handler, g.logger))
}

func (g *echoGroup) PUT(path string, handler Handler) {
	g.EchoGroup.PUT(path, newEchoHandler(handler, g.logger))
}

func (g *echoGroup) DELETE(path string, handler Handler) {
	g.EchoGroup.DELETE(path, newEchoHandler(handler, g.logger))
}

func (g *echoGroup) PATCH(path string, handler Handler) {
	g.EchoGroup.PATCH(path, newEchoHandler(handler, g.logger))
}

func (g *echoGroup) Group(prefix string, m ...echo.MiddlewareFunc) *echoGroup {
	g2 := g.EchoGroup.Group(prefix, m...)
	return &echoGroup{g2, g.logger}
}

func (g *echoGroup) Use(m ...Handler) {
	g.EchoGroup.Use(newEchoMiddlewares(g.logger, m...)...)
}
