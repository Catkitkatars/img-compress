package main

import (
	"img-compress/internal/app"
	"img-compress/internal/config"
	"img-compress/internal/dto"
	"img-compress/internal/handler"
	srv "img-compress/internal/http"
	"img-compress/internal/storage"
	"io"
	"log"
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

	logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("не удалось открыть лог-файл: %v", err)
	}

	logWriter := io.MultiWriter(os.Stdout, logFile)
	var slogHandler slog.Handler

	switch cfg.Env {
	case EnvLocal:
		slogHandler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		slogHandler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		slogHandler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(slogHandler)

	logger.Info("starting img-compress", slog.String("env", cfg.Env))

	store, err := storage.New(cfg.StoragePath, logger)

	if err != nil {
		logger.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	configDto := &dto.Config{
		Cfg:     cfg,
		Storage: store,
		Log:     logger,
	}

	imageApp := app.NewImage(configDto)
	imageHandler := handler.NewImageHandler(configDto, imageApp)

	srvErr := srv.Start(configDto, imageHandler)

	if srvErr != nil {
		logger.Error("failed to start server", err)
		logger.Error("server stopped", nil)
		os.Exit(1)
	}
}
