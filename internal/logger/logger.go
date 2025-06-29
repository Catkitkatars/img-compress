package logger

import (
	"img-compress/internal/config"
	"io"
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var Logger *slog.Logger

func New(file *os.File) {
	cfg := config.Cfg
	logWriter := io.MultiWriter(os.Stdout, file)
	var slogHandler slog.Handler

	switch cfg.Env {
	case EnvLocal:
		slogHandler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		slogHandler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		slogHandler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	Logger = slog.New(slogHandler)
}
