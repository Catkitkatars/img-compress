package main

import (
	"img-compress/internal/config"
	"img-compress/internal/lib/logger/slog"
	"img-compress/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var cfg = config.Init()

	log := initLogger(cfg.Env)

	log.Info("starting img-compress", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", sl.Error(err))
		os.Exit(1)
	}

	_ = storage

	// TODO
	// config - cleanenv
	// logger - slog
	// storage - sqlite
	// router - chi
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
