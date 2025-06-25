package main

import (
	"img-compress/internal/app"
	"img-compress/internal/config"
	"img-compress/internal/dto"
	"img-compress/internal/handler"
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

	var slogHandler slog.Handler

	switch cfg.Env {
	case EnvLocal:
		slogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	log := slog.New(slogHandler)

	log.Info("starting img-compress", slog.String("env", cfg.Env))

	store, err := storage.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	ImageAppDto := dto.ImageApp{
		Cfg:     cfg,
		Storage: store,
		Log:     log,
	}

	//todo add app dto with ImageAppDto, ImageHandlerDto, etc.

	imageApp := app.NewImage(ImageAppDto)
	imageHandler := handler.NewImageHandler(imageApp, ImageAppDto)

	srvErr := srv.Start(cfg, imageHandler)

	if srvErr != nil {
		log.Error("failed to start server", err)
		log.Error("server stopped", nil)
		os.Exit(1)
	}
}
