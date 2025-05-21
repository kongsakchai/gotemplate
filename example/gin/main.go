package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/app/todo"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/logger"
	"github.com/kongsakchai/gotemplate/middleware"
)

var (
	conf config.Config
)

func init() {
	conf = config.Load()
	logger.SetLevel(conf.App.LogLevel)
}

func main() {
	log := logger.New()
	r := setupRoutes(log)

	idle := make(chan struct{})
	go gracefulShutdown(func(ctx context.Context) error {
		defer close(idle)
		return r.Shutdown(ctx)
	})

	if err := r.Start(":" + conf.App.Port); err != nil && err != http.ErrServerClosed {
		log.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	log.Info("bye bye")
}

func setupRoutes(log *slog.Logger) app.App {
	r := app.NewGinRouter(log)
	r.Validator = middleware.NewReqValidator()
	r.Use(logger.GinLogger())

	r.GET("/hello", r.NewHandler(func(ctx app.Context) error {
		return ctx.OK("hello world")
	}))

	{
		service := todo.NewService()
		h := todo.NewHandler(service)

		r.GET("/todos", r.NewHandler(h.Todos))
		r.GET("/todos/:id", r.NewHandler(h.Todo))
		r.POST("/todos", r.NewHandler(h.Create))
		r.PUT("/todos/:id", r.NewHandler(h.Update))
		r.DELETE("/todos/:id", r.NewHandler(h.Delete))
	}

	return r
}

func gracefulShutdown(close func(context.Context) error) {
	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-sig.Done()

	logger.Info("shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := close(ctx); err != nil {
		logger.Error("graceful shutdown failed: " + err.Error())
		return
	}
	logger.Info("graceful shutdown completed")
}
