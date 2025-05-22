package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kongsakchai/gotemplate/pkg/generate"
)

var (
	ginContextPool *pool[ginContext]
)

func init() {
	ginContextPool = createPool[ginContext]()
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

func (g *ginRouter) NewHandler(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ginContextPool.Get()
		defer ginContextPool.Put(ctx)
		ctx.reset(c, g.logger, g.Validator)
		handler(ctx)
	}
}
