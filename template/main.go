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

	"github.com/kongsakchai/gotemplate/common/logger"
	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/kongsakchai/gotemplate/template/config"
	migrate "github.com/kongsakchai/simple-sql-migrate"
)

var gracefulTimeout = time.Second * 10

func init() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(1) // 0 - 999m
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1) // GOMEMLIMIT
	}
}

func main() {
	cfg := config.Load(config.Env)
	log := logger.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	external, closeExternal := setupExternalService(cfg)
	defer closeExternal(context.Background())

	app := router(cfg, external, log)

	idle := make(chan struct{})
	go gracefulShutdown(idle, ctx, app)

	log.Info("starting "+cfg.App.Name, "version", cfg.App.Version, "env", config.Env)
	log.Info("listening on port " + cfg.App.Port)

	if err := app.Start(ctx, ":"+cfg.App.Port); err != nil && err != http.ErrServerClosed {
		log.Error("shutting down the server: " + err.Error())
		return
	}

	<-idle
	log.Info("bye bye")
}

func gracefulShutdown(idle chan struct{}, signal context.Context, app app.App) {
	defer close(idle)
	<-signal.Done()

	slog.Info("shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		panic("graceful shutdown failed: " + err.Error())
	}
	slog.Info("graceful shutdown completed")
}

func migrateDB(db *sql.DB, cfg config.Migration) {
	if !cfg.Enable {
		return
	}
	if err := migrate.Migrate(db, migrate.Options{
		Source:  cfg.Directory,
		Version: cfg.Version,
		Repeat:  migrate.GetRepeatAction(cfg.Repeat),
	}); err != nil {
		panic("migration failed: " + err.Error())
	}

	version, _ := migrate.Version(db)
	slog.Info("migration completed", "version", version)
}
