package http_server

import (
	"img-compress/internal/config"
	"img-compress/internal/http-server/router"
	sl "img-compress/internal/lib/logger/slog"
	"net/http"
	"os"
)

func HttpStart(cfg *config.Config, log *sl.Logger) {
	r := router.Init()

	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
		os.Exit(1)
	}

	log.Error("server stopped", nil)
}
