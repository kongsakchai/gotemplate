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
	echo.Context
}

func newEchoContext(logger *slog.Logger, ctx echo.Context) *echoContext {
	return &echoContext{
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

func (e *echoContext) Request() *http.Request {
	return e.Context.Request()
}

func (e *echoContext) Logger() *slog.Logger {
	return e.logger
}

func (e *echoContext) SetLogger(logger *slog.Logger) {
	e.logger = logger
}

func (e *echoContext) Writer() http.ResponseWriter {
	return e.Context.Response().Writer
}

func (e *echoContext) SetWriter(w http.ResponseWriter) {
	e.Context.Response().Writer = w
}

func newEchoHandler(handler Handler, middlewares []Middleware, logger *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		h := applyMiddleware(handler, middlewares)
		return h(newEchoContext(logger, ctx))
	}
}

type echoRouter struct {
	*echo.Echo
	logger      *slog.Logger
	middlewares []Middleware
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

func (e *echoRouter) Origin() any {
	return e.Echo
}

func (e *echoRouter) GET(path string, handler Handler) {
	e.Echo.GET(path, newEchoHandler(handler, e.middlewares, e.logger))
}

func (e *echoRouter) POST(path string, handler Handler) {
	e.Echo.POST(path, newEchoHandler(handler, e.middlewares, e.logger))
}

func (e *echoRouter) PUT(path string, handler Handler) {
	e.Echo.PUT(path, newEchoHandler(handler, e.middlewares, e.logger))
}

func (e *echoRouter) DELETE(path string, handler Handler) {
	e.Echo.DELETE(path, newEchoHandler(handler, e.middlewares, e.logger))
}

func (e *echoRouter) PATCH(path string, handler Handler) {
	e.Echo.PATCH(path, newEchoHandler(handler, e.middlewares, e.logger))
}

func (e *echoRouter) Use(middlewares ...Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

type echoGroup struct {
	EchoGroup   *echo.Group
	logger      *slog.Logger
	middlewares []Middleware
}

func (e *echoRouter) Group(prefix string, middlewares ...Middleware) RouterGroup {
	g := e.Echo.Group(prefix)
	return &echoGroup{
		EchoGroup:   g,
		logger:      e.logger,
		middlewares: append(copyMiddlewares(e.middlewares, middlewares...), middlewares...),
	}
}

func (g *echoGroup) Origin() any {
	return g.EchoGroup
}

func (g *echoGroup) GET(path string, handler Handler) {
	g.EchoGroup.GET(path, newEchoHandler(handler, g.middlewares, g.logger))
}

func (g *echoGroup) POST(path string, handler Handler) {
	g.EchoGroup.POST(path, newEchoHandler(handler, g.middlewares, g.logger))
}

func (g *echoGroup) PUT(path string, handler Handler) {
	g.EchoGroup.PUT(path, newEchoHandler(handler, g.middlewares, g.logger))
}

func (g *echoGroup) DELETE(path string, handler Handler) {
	g.EchoGroup.DELETE(path, newEchoHandler(handler, g.middlewares, g.logger))
}

func (g *echoGroup) PATCH(path string, handler Handler) {
	g.EchoGroup.PATCH(path, newEchoHandler(handler, g.middlewares, g.logger))
}

func (g *echoGroup) Group(prefix string, middlewares ...Middleware) RouterGroup {
	g2 := g.EchoGroup.Group(prefix)
	return &echoGroup{
		EchoGroup:   g2,
		logger:      g.logger,
		middlewares: append(copyMiddlewares(g.middlewares, middlewares...), middlewares...),
	}
}

func (g *echoGroup) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}
