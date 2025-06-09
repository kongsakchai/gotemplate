package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/logger"
	"github.com/kongsakchai/gotemplate/middleware"
	"github.com/kongsakchai/gotemplate/validator"
	migrate "github.com/kongsakchai/simple-migrate"
	"github.com/labstack/echo/v4"
)

func init() {
	if os.Getenv("GOMAXPROCS") != "" {
		runtime.GOMAXPROCS(0) // GOMAXPROCS
	} else {
		runtime.GOMAXPROCS(1) // 0 - 999m
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1) // GOMEMLIMIT
	}
}

func main() {
	cfg := config.Load()
	log := logger.New()
	r := setupRoutes(cfg)

	setMigration(cfg.Migration)

	idle := make(chan struct{})
	go gracefulShutdown(func(ctx context.Context) error {
		defer close(idle)
		return r.Shutdown(ctx)
	})

	log.Info("starting " + cfg.App.Name + " version " + cfg.App.Version)
	log.Info("listening on port " + cfg.App.Port)

	if err := r.Start(":" + cfg.App.Port); err != nil && err != http.ErrServerClosed {
		log.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	log.Info("bye bye")
}

func setupRoutes(cfg config.Config) app.App {
	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()

	r.Use(
		middleware.RefID(cfg.Header.RefIDKey),
		middleware.Logger(cfg.Header.RefIDKey, true, true),
	)

	r.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
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

func setMigration(cfg config.Migration) {
	var err error
	if cfg.Version != "" {
		err = migrate.New(nil, cfg.Directory).SetVersion(cfg.Version)
	} else {
		err = migrate.New(nil, cfg.Directory).Up()
	}

	if err != nil {
		panic("migration failed: " + err.Error())
	}
}
