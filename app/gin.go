package app

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ginContext struct {
	logger *slog.Logger
	*gin.Context
}

var ginContextPool = &sync.Pool{
	New: func() any {
		return &ginContext{}
	},
}

func putGinContext(ctx *ginContext) {
	ginContextPool.Put(ctx)
}

func newGinContext(logger *slog.Logger, ctx *gin.Context) *ginContext {
	c := ginContextPool.Get().(*ginContext)
	c.reset(ctx, logger)
	return c
}

func (g *ginContext) reset(ctx *gin.Context, logger *slog.Logger) {
	g.Context = ctx
	g.logger = logger.With(slog.String("traceID", ctx.GetString("traceID")))
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

func (g *ginContext) OKWithMessage(message string, obj any) error {
	return g.JSON(200, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
	})
}

func (g *ginContext) Created(obj any) error {
	return g.JSON(201, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (g *ginContext) CreatedWithMessage(message string, obj any) error {
	return g.JSON(201, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
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

func (g *ginContext) Logger() *slog.Logger {
	return g.logger
}

func newGinHandler(handler Handler, logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := newGinContext(logger, ctx)
		defer putGinContext(c)
		handler(c)
	}
}

type ginRouter struct {
	*gin.Engine
	logger *slog.Logger
	serv   *http.Server
}

func NewGinRouter(logger *slog.Logger) *ginRouter {
	r := gin.New()
	r.Use(gin.Recovery())

	return &ginRouter{
		Engine: r,
		logger: logger,
	}
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

func (g *ginRouter) GET(path string, handler Handler, m ...gin.HandlerFunc) {
	m = append(m, newGinHandler(handler, g.logger))
	g.Engine.GET(path, m...)
}

func (g *ginRouter) POST(path string, handler Handler, m ...gin.HandlerFunc) {
	m = append(m, newGinHandler(handler, g.logger))
	g.Engine.POST(path, m...)
}

func (g *ginRouter) PUT(path string, handler Handler, m ...gin.HandlerFunc) {
	m = append(m, newGinHandler(handler, g.logger))
	g.Engine.PUT(path, m...)
}

func (g *ginRouter) DELETE(path string, handler Handler, m ...gin.HandlerFunc) {
	m = append(m, newGinHandler(handler, g.logger))
	g.Engine.DELETE(path, m...)
}

func (g *ginRouter) PATCH(path string, handler Handler, m ...gin.HandlerFunc) {
	m = append(m, newGinHandler(handler, g.logger))
	g.Engine.PATCH(path, m...)
}

type ginGroup struct {
	*gin.RouterGroup
	logger *slog.Logger
}

func (g *ginRouter) Group(prefix string, m ...gin.HandlerFunc) *ginGroup {
	grp := g.Engine.Group(prefix, m...)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
	}
}

func (g *ginGroup) GET(path string, handler Handler, m ...gin.HandlerFunc) {
	g.RouterGroup.GET(path, newGinHandler(handler, g.logger))
}

func (g *ginGroup) POST(path string, handler Handler, m ...gin.HandlerFunc) {
	g.RouterGroup.POST(path, newGinHandler(handler, g.logger))
}

func (g *ginGroup) PUT(path string, handler Handler, m ...gin.HandlerFunc) {
	g.RouterGroup.PUT(path, newGinHandler(handler, g.logger))
}

func (g *ginGroup) DELETE(path string, handler Handler, m ...gin.HandlerFunc) {
	g.RouterGroup.DELETE(path, newGinHandler(handler, g.logger))
}

func (g *ginGroup) PATCH(path string, handler Handler, m ...gin.HandlerFunc) {
	g.RouterGroup.PATCH(path, newGinHandler(handler, g.logger))
}

func (g *ginGroup) Group(prefix string, m ...gin.HandlerFunc) *ginGroup {
	grp := g.RouterGroup.Group(prefix, m...)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
	}
}
