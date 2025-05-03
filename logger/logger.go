package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var (
	logger   *slog.Logger
	logLevel slog.Level
)

func SetLevel(level string) {
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}
}

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	return logger
}

func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !logger.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, this function, this function's caller]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

func Debug(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	log(context.Background(), slog.LevelInfo, msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	log(context.Background(), slog.LevelWarn, msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	log(context.Background(), slog.LevelError, msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelError, msg, args...)
}
