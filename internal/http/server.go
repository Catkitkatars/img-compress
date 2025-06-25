package http

import (
	"img-compress/internal/config"
	"img-compress/internal/handler"
	"img-compress/internal/http/router"
	"net/http"
)

func Start(cfg *config.Config, imageHandler *handler.ImageHandler) error {
	r := router.Init(imageHandler)

	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	err := srv.ListenAndServe()

	return err
}
