package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/database"
	"github.com/kongsakchai/gotemplate/logger"
	"github.com/kongsakchai/gotemplate/middleware"
	"github.com/kongsakchai/gotemplate/validator"
	migrate "github.com/kongsakchai/simple-migrate"
	"github.com/labstack/echo/v4"
)

var gracefulTimeout = time.Second * 10

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
	cfg := config.Load(config.Env)
	log := logger.New()

	r, close := setupRoutes(cfg)

	idle := make(chan struct{})
	go gracefulShutdown(idle, close)

	log.Info("starting "+cfg.App.Name, "version", cfg.App.Version, "env", config.Env)
	log.Info("listening on port " + cfg.App.Port)

	if err := r.Start(":" + cfg.App.Port); err != nil && err != http.ErrServerClosed {
		slog.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	log.Info("bye bye")
}

func setupRoutes(cfg config.Config) (app.App, func(context.Context) error) {
	db, closeDB := database.NewMySQL(cfg.Database)
	setMigration(db, cfg.Migration)

	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()

	shutdown := func(ctx context.Context) error {
		if err := r.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown %w", err)
		}
		if err := closeDB(); err != nil {
			return fmt.Errorf("database shutdown %w", err)
		}
		return nil
	}

	r.Use(
		middleware.RefID(cfg.Header.RefIDKey),
		middleware.Logger(cfg.Header.RefIDKey, true, true),
	)

	r.GET("/health", healthCheck(db))

	return r, shutdown
}

func healthCheck(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if db.Ping() != nil {
			return app.Fail(ctx, app.InternalServer(app.ErrCode, app.ErrDatabaseMsg, nil))
		}
		return app.OkWithMessage(ctx, nil, "healthy")
	}
}

func gracefulShutdown(idle chan struct{}, shutdown func(context.Context) error) {
	defer close(idle)

	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-sig.Done()

	slog.Info("shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed: " + err.Error())
		panic("force shutdown")
	}
	slog.Info("graceful shutdown completed")
}

func setMigration(db *sql.DB, cfg config.Migration) {
	var err error
	if cfg.Version != "" {
		err = migrate.New(db, cfg.Directory).SetVersion(cfg.Version)
	} else {
		err = migrate.New(db, cfg.Directory).Up()
	}

	if err != nil {
		panic("migration failed: " + err.Error())
	}
}
