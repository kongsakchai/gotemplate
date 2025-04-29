package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/internal/ping"
)

func main() {
	conf := config.Load()
	e := app.NewEcho(conf.App.LogLevel)

	{
		pingHandler := ping.NewHandler()

		e.GET("/ping", pingHandler.Ping)
	}

	go func() {
		if err := e.Start(":" + conf.App.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()

	gracefulShutdown(func(ctx context.Context) {
		e.Logger.Info("gracefully shutting down the server")

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal("gracefully shutting down the server", err)
		}
	})
}

func gracefulShutdown(close func(context.Context)) {
	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-sig.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	close(ctx)
}
