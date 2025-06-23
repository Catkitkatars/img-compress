package main

import (
	"img-compress/internal/config"
	"img-compress/internal/lib/logger/slog"
	"img-compress/internal/storage/sqlite"
	"log/slog"
	"os"
)

func main() {
	cfg := config.Init()

	log := sl.New(cfg.Env)

	log.Info("starting img-compress", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	_ = storage

	// TODO
	// config - cleanenv
	// logger - slog
	// storage - sqlite
	// router - chi
}
