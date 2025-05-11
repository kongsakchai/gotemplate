package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginContext struct {
	logger *slog.Logger
	*gin.Context
}

func newGinContext(logger *slog.Logger, ctx *gin.Context) *ginContext {
	return &ginContext{
		logger:  logger,
		Context: ctx,
	}
}

func (g *ginContext) Query(key string) string {
	return g.Context.Query(key)
}

func (g *ginContext) Param(key string) string {
	return g.Context.Param(key)
}

func (g *ginContext) Bind(obj any) error {
	return g.Context.Bind(obj)
}

func (g *ginContext) JSON(code int, obj any) error {
	g.Context.JSON(code, obj)
	return nil
}

func (g *ginContext) OK(obj any) error {
	return g.JSON(200, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (g *ginContext) Created(obj any) error {
	return g.JSON(201, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (g *ginContext) NotFound(err Error) error {
	g.logger.Error(err.Error())
	return g.JSON(404, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (g *ginContext) InternalServer(err Error) error {
	g.logger.Error(err.Error())
	return g.JSON(500, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (g *ginContext) BadRequest(err Error) error {
	g.logger.Error(err.Error())
	return g.JSON(400, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (g *ginContext) Ctx() context.Context {
	return g.Context.Request.Context()
}

func (g *ginContext) Get(key string) any {
	val, ok := g.Context.Get(key)
	if !ok {
		return nil
	}
	return val
}

func (g *ginContext) Set(key string, value any) {
	g.Context.Set(key, value)
}

func (g *ginContext) Request() *http.Request {
	return g.Context.Request
}

func (g *ginContext) Logger() *slog.Logger {
	return g.logger
}

func (g *ginContext) SetLogger(logger *slog.Logger) {
	g.logger = logger
}

func (g *ginContext) Writer() http.ResponseWriter {
	return g.Context.Writer
}

type ginResponseWriter struct {
	gin.ResponseWriter
	custom http.ResponseWriter
}

func (g *ginResponseWriter) Write(b []byte) (int, error) {
	return g.custom.Write(b)
}

func (g *ginResponseWriter) WriteHeader(code int) {
	g.custom.WriteHeader(code)
}

func (g *ginContext) SetWriter(w http.ResponseWriter) {
	g.Context.Writer = &ginResponseWriter{
		ResponseWriter: g.Context.Writer,
		custom:         w,
	}
}

func newGinHandler(handler Handler, middlewares []Middleware, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := applyMiddleware(handler, middlewares)
		h(newGinContext(logger, c))
	}
}

type ginRouter struct {
	*gin.Engine
	middlewares []Middleware
	logger      *slog.Logger
	serv        *http.Server
}

func NewGinRouter(logger *slog.Logger) *ginRouter {
	r := gin.New()
	r.Use(gin.Recovery())

	return &ginRouter{
		Engine: r,
		logger: logger,
	}
}

func (g *ginRouter) Origin() any {
	return g.Engine
}

func (g *ginRouter) Shutdown(ctx context.Context) error {
	if g.serv == nil {
		return nil
	}

	return g.serv.Shutdown(ctx)
}

func (g *ginRouter) Start(addr string) error {
	g.serv = &http.Server{
		Addr:    addr,
		Handler: g.Engine,
	}

	return g.serv.ListenAndServe()
}

func (g *ginRouter) GET(path string, handler Handler) {
	g.Engine.GET(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginRouter) POST(path string, handler Handler) {
	g.Engine.POST(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginRouter) PUT(path string, handler Handler) {
	g.Engine.PUT(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginRouter) DELETE(path string, handler Handler) {
	g.Engine.DELETE(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginRouter) PATCH(path string, handler Handler) {
	g.Engine.PATCH(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginRouter) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

type ginGroup struct {
	*gin.RouterGroup
	middlewares []Middleware
	logger      *slog.Logger
}

func (g *ginRouter) Group(prefix string, middlewares ...Middleware) RouterGroup {
	grp := g.Engine.Group(prefix)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
		middlewares: append(copyMiddlewares(g.middlewares), middlewares...),
	}
}

func (g *ginGroup) Origin() any {
	return g.RouterGroup
}

func (g *ginGroup) GET(path string, handler Handler) {
	g.RouterGroup.GET(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginGroup) POST(path string, handler Handler) {
	g.RouterGroup.POST(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginGroup) PUT(path string, handler Handler) {
	g.RouterGroup.PUT(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginGroup) DELETE(path string, handler Handler) {
	g.RouterGroup.DELETE(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginGroup) PATCH(path string, handler Handler) {
	g.RouterGroup.PATCH(path, newGinHandler(handler, g.middlewares, g.logger))
}

func (g *ginGroup) Use(middleware ...Middleware) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *ginGroup) Group(prefix string, middlewares ...Middleware) RouterGroup {
	grp := g.RouterGroup.Group(prefix)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
		middlewares: append(copyMiddlewares(g.middlewares), middlewares...),
	}
}
