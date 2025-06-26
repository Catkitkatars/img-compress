package dto

import (
	"img-compress/internal/config"
	"img-compress/internal/storage"
	"log/slog"
)

type Config struct {
	Cfg     *config.Config
	Storage *storage.Storage
	Log     *slog.Logger
}
