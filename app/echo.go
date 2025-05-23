package app

import (
	"context"
	"log/slog"

	"github.com/kongsakchai/gotemplate/pkg/generate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	echoContextPool *pool[*echoContext]
)

func init() {
	echoContextPool = createPool[*echoContext](func() any {
		return &echoContext{}
	})
}

type echoContext struct {
	echo.Context
	logger *slog.Logger
	next   func(c echo.Context) error
}

func (e *echoContext) reset(ctx echo.Context, logger *slog.Logger) {
	e.Context = ctx
	e.logger = logger
	e.next = nil
}

func (e *echoContext) Next() error {
	if e.next != nil {
		return e.next(e.Context)
	}
	return nil
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

func (e *echoContext) Validate(obj any) error {
	return e.Context.Validate(obj)
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

func (e *echoContext) Error(err Error) error {
	e.logger.Error(err.Error())
	return e.JSON(err.StatusCd, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *echoContext) Original() any {
	return e.Context
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

func newEchoMiddleware(logger *slog.Logger, handler Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := echoContextPool.Get()
			defer echoContextPool.Put(c)
			c.reset(ctx, logger)
			c.next = next

			return handler(c)
		}
	}
}

func newEchoMiddlewares(logger *slog.Logger, handlers ...Handler) []echo.MiddlewareFunc {
	m := make([]echo.MiddlewareFunc, len(handlers))
	for i, handler := range handlers {
		m[i] = newEchoMiddleware(logger, handler)
	}
	return m
}

func newEchoHandler(logger *slog.Logger, handler Handler) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := echoContextPool.Get()
		defer echoContextPool.Put(c)
		c.reset(ctx, logger)
		return handler(c)
	}
}

func newEchoHandlers(logger *slog.Logger, handlers ...Handler) (echo.HandlerFunc, []echo.MiddlewareFunc) {
	return newEchoHandler(logger, handlers[len(handlers)-1]), newEchoMiddlewares(logger, handlers[:len(handlers)-1]...)
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
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("traceID", generate.UUID())
			return next(c)
		}
	})

	return &echoRouter{
		Echo:   e,
		logger: logger,
	}
}

func (e *echoRouter) GET(path string, handlers ...Handler) {
	h, m := newEchoHandlers(e.logger, handlers...)
	e.Echo.GET(path, h, m...)
}

func (e *echoRouter) POST(path string, handlers ...Handler) {
	h, m := newEchoHandlers(e.logger, handlers...)
	e.Echo.POST(path, h, m...)
}

func (e *echoRouter) PUT(path string, handlers ...Handler) {
	h, m := newEchoHandlers(e.logger, handlers...)
	e.Echo.PUT(path, h, m...)
}

func (e *echoRouter) DELETE(path string, handlers ...Handler) {
	h, m := newEchoHandlers(e.logger, handlers...)
	e.Echo.DELETE(path, h, m...)
}

func (e *echoRouter) PATCH(path string, handlers ...Handler) {
	h, m := newEchoHandlers(e.logger, handlers...)
	e.Echo.PATCH(path, h, m...)
}

func (e *echoRouter) Use(middlewares ...Handler) {
	e.Echo.Use(newEchoMiddlewares(e.logger, middlewares...)...)
}

type echoGroup struct {
	EchoGroup *echo.Group
	logger    *slog.Logger
}

func (e *echoRouter) Group(prefix string, m ...Handler) RouteGroup {
	grp := e.Echo.Group(prefix, newEchoMiddlewares(e.logger, m...)...)
	return &echoGroup{
		EchoGroup: grp,
		logger:    e.logger,
	}
}

func (g *echoGroup) GET(path string, handlers ...Handler) {
	h, m := newEchoHandlers(g.logger, handlers...)
	g.EchoGroup.GET(path, h, m...)
}

func (g *echoGroup) POST(path string, handlers ...Handler) {
	h, m := newEchoHandlers(g.logger, handlers...)
	g.EchoGroup.POST(path, h, m...)
}

func (g *echoGroup) PUT(path string, handlers ...Handler) {
	h, m := newEchoHandlers(g.logger, handlers...)
	g.EchoGroup.PUT(path, h, m...)
}

func (g *echoGroup) DELETE(path string, handlers ...Handler) {
	h, m := newEchoHandlers(g.logger, handlers...)
	g.EchoGroup.DELETE(path, h, m...)
}

func (g *echoGroup) PATCH(path string, handlers ...Handler) {
	h, m := newEchoHandlers(g.logger, handlers...)
	g.EchoGroup.PATCH(path, h, m...)
}

func (g *echoGroup) Group(prefix string, m ...Handler) RouteGroup {
	grp := g.EchoGroup.Group(prefix, newEchoMiddlewares(g.logger, m...)...)
	return &echoGroup{
		EchoGroup: grp,
		logger:    g.logger,
	}
}

func (g *echoGroup) Use(middlewares ...Handler) {
	g.EchoGroup.Use(newEchoMiddlewares(g.logger, middlewares...)...)
}
