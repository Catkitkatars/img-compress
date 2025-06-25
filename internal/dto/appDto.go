package dto

import (
	"img-compress/internal/config"
	"img-compress/internal/storage"
	"log/slog"
)

type AppDto struct {
	Cfg        *config.Config
	Storage    *storage.Storage
	Log        *slog.Logger
	imgHandler *HandlersDto
}
