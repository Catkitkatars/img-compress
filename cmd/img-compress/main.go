package main

import (
	"img-compress/internal/config"
	srv "img-compress/internal/http-server"
	"img-compress/internal/lib/logger/slog"
	"img-compress/internal/storage/sqlite"
	"log/slog"
	"os"
)

func main() {
	cfg := config.Init()

	log := sl.New(cfg.Env)

	log.Info("starting img-compress", slog.String("env", cfg.Env))

	_, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	srv.HttpStart(cfg, log)
}
