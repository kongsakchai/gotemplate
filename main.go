package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/logger"
	migrate "github.com/kongsakchai/simple-migrate"
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

	app, shutdowns := router(cfg)

	idle := make(chan struct{})
	go gracefulShutdown(idle, shutdowns...)

	log.Info("starting "+cfg.App.Name, "version", cfg.App.Version, "env", config.Env)
	log.Info("listening on port " + cfg.App.Port)

	if err := app.Start(":" + cfg.App.Port); err != nil && err != http.ErrServerClosed {
		slog.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	log.Info("bye bye")
}

type shutdownFunc func(context.Context) error

func gracefulShutdown(idle chan struct{}, shutdowns ...shutdownFunc) {
	defer close(idle)

	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-sig.Done()

	slog.Info("shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	for _, shutdown := range shutdowns {
		if err := shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed: " + err.Error())
			panic("force shutdown")
		}
	}
	slog.Info("graceful shutdown completed")
}

func setMigration(db *sql.DB, cfg config.Migration) {
	if !cfg.Enable {
		return
	}

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
