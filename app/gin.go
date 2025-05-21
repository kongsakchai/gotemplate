package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kongsakchai/gotemplate/pkg/generate"
)

var (
	ginContextPool *pool[*ginContext]
)

func init() {
	ginContextPool = createPool[*ginContext](func() any {
		return &ginContext{}
	})
}

type ginContext struct {
	*gin.Context
	logger    *slog.Logger
	validator Validator
}

func (g *ginContext) reset(ctx *gin.Context, logger *slog.Logger, validator Validator) {
	g.Context = ctx
	g.logger = logger
	g.validator = validator
}

func (g *ginContext) Next() error {
	g.Context.Next()
	return nil
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

func (g *ginContext) Validate(obj any) error {
	if g.validator == nil {
		return nil
	}
	return g.validator.Validate(obj)
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

func (g *ginContext) Error(err Error) error {
	g.logger.Error(err.Error())
	return g.JSON(err.StatusCd, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (g *ginContext) Original() any {
	return g.Context
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

func newGinHandler(logger *slog.Logger, validator Validator, handler Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := ginContextPool.Get()
		defer ginContextPool.Put(c)
		c.reset(ctx, logger, validator)

		handler(c)
	}
}

func newGinHandlers(logger *slog.Logger, validator Validator, handlers ...Handler) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = newGinHandler(logger, validator, handler)
	}
	return ginHandlers
}

type ginRouter struct {
	*gin.Engine
	logger *slog.Logger
	serv   *http.Server

	Validator Validator
}

func NewGinRouter(logger *slog.Logger) *ginRouter {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("traceID", generate.UUID())
		c.Next()
	})

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

func (g *ginRouter) GET(path string, handlers ...Handler) {
	g.Engine.GET(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginRouter) POST(path string, handlers ...Handler) {
	g.Engine.POST(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginRouter) PUT(path string, handlers ...Handler) {
	g.Engine.PUT(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginRouter) DELETE(path string, handlers ...Handler) {
	g.Engine.DELETE(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginRouter) PATCH(path string, handlers ...Handler) {
	g.Engine.PATCH(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginRouter) Use(m ...Handler) {
	g.Engine.Use(newGinHandlers(g.logger, g.Validator, m...)...)
}

type ginGroup struct {
	*gin.RouterGroup
	logger    *slog.Logger
	Validator Validator
}

func (g *ginRouter) Group(prefix string, m ...Handler) RouteGroup {
	grp := g.Engine.Group(prefix, newGinHandlers(g.logger, g.Validator, m...)...)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
		Validator:   g.Validator,
	}
}

func (g *ginGroup) GET(path string, handlers ...Handler) {
	g.RouterGroup.GET(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginGroup) POST(path string, handlers ...Handler) {
	g.RouterGroup.POST(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginGroup) PUT(path string, handlers ...Handler) {
	g.RouterGroup.PUT(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginGroup) DELETE(path string, handlers ...Handler) {
	g.RouterGroup.DELETE(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginGroup) PATCH(path string, handlers ...Handler) {
	g.RouterGroup.PATCH(path, newGinHandlers(g.logger, g.Validator, handlers...)...)
}

func (g *ginGroup) Group(prefix string, m ...Handler) RouteGroup {
	grp := g.RouterGroup.Group(prefix, newGinHandlers(g.logger, g.Validator, m...)...)
	return &ginGroup{
		RouterGroup: grp,
		logger:      g.logger,
	}
}

func (g *ginGroup) Use(m ...Handler) {
	g.RouterGroup.Use(newGinHandlers(g.logger, g.Validator, m...)...)
}
