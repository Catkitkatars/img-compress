package http

import (
	"img-compress/internal/config"
	"img-compress/internal/http/router"
	logs "img-compress/internal/logger"
	"net/http"
)

func Start() error {
	cfg := config.Cfg
	r := router.Init()

	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	logs.Logger.Info("server start")
	err := srv.ListenAndServe()

	return err
}
