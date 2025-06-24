package main

import (
	"img-compress/internal/config"
	srv "img-compress/internal/http"
	"img-compress/internal/storage"
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	cfg := config.Init()

	var handler slog.Handler

	switch cfg.Env {
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	log := slog.New(handler)

	log.Info("starting img-compress", slog.String("env", cfg.Env))

	_, err := storage.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	srvErr := srv.Start(cfg)

	if srvErr != nil {
		log.Error("failed to start server", err)
		log.Error("server stopped", nil)
		os.Exit(1)
	}
}
