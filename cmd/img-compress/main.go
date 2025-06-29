package main

import (
	"img-compress/internal/config"
	srv "img-compress/internal/http"
	slogger "img-compress/internal/logger"
	"img-compress/internal/storage"
	"log"
	"os"
)

func main() {
	config.New()
	cfg := &config.Cfg

	logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("fail to opened log file: %v", err)
	}

	slogger.New(logFile)
	logger := slogger.Logger

	logger.Info("starting img-compress")

	storeErr := storage.New(cfg.StoragePath)

	logger.Info("storage start")

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
