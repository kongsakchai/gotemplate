package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/logger"
	"github.com/kongsakchai/gotemplate/middleware"
	"github.com/kongsakchai/gotemplate/validator"
	"github.com/labstack/echo/v4"
)

func main() {
	logger.New()
	r := setupRoutes()

	idle := make(chan struct{})
	go gracefulShutdown(func(ctx context.Context) error {
		defer close(idle)
		return r.Shutdown(ctx)
	})

	if err := r.Start(":8080"); err != nil && err != http.ErrServerClosed {
		slog.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	slog.Info("bye bye")
}

func setupRoutes() app.App {
	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()

	r.Use(
		middleware.EchoRefID("X-REF-ID"),
		middleware.Logger("X-REF-ID", true, true),
	)

	r.GET("/ping", func(c echo.Context) error {
		return app.OkWithMessage(c, "pong", nil)
	})

	return r
}

func gracefulShutdown(close func(context.Context) error) {
	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-sig.Done()

	slog.Info("shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := close(ctx); err != nil {
		slog.Error("graceful shutdown failed: " + err.Error())
		return
	}
	slog.Info("graceful shutdown completed")
}
