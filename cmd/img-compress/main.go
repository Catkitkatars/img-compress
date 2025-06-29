package main

import (
	"img-compress/internal/app"
	"img-compress/internal/config"
	"img-compress/internal/handler"
	srv "img-compress/internal/http"
	slogger "img-compress/internal/logger"
	"img-compress/internal/storage"
	"log"
	"log/slog"
	"os"
)

func main() {
	config.New()
	cfg := &config.Cfg

	logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("не удалось открыть лог-файл: %v", err)
	}

	slogger.New(logFile)
	logger := slogger.Logger

	logger.Info("starting img-compress", slog.String("env", cfg.Env))

	storeErr := storage.New(cfg.StoragePath)

	if storeErr != nil {
		logger.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	srvErr := srv.Start()

	if srvErr != nil {
		logger.Error("failed to start server", err)
		logger.Error("server stopped", nil)
		os.Exit(1)
	}
}
