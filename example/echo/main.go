package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/logger"
)

var (
	conf config.Config
)

func init() {
	runtime.GOMAXPROCS(1)
	conf = config.Load()
	logger.SetLevel(conf.App.LogLevel)
}

func main() {
	log := logger.New()
	r := setupRoutes(app.NewEchoRoute(log)) // Initialize your router here

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

func setupRoutes(r app.Router) app.Router {
	r.Use(logger.LoggerRequest())
	r.Use(logger.LoggerResponse())

	r.GET("/hello", func(c app.Context) error {
		return c.JSON(http.StatusOK, "Hello, World")
	})

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
