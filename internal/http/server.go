package http

import (
	"img-compress/internal/dto"
	"img-compress/internal/handler"
	"img-compress/internal/http/router"
	"net/http"
)

func Start(dto *dto.Config, imgHandler *handler.ImageHandler) error {
	r := router.Init(imgHandler)

	srv := &http.Server{
		Addr:         dto.Cfg.HTTP.Host + ":" + dto.Cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  dto.Cfg.HTTP.Timeout,
		WriteTimeout: dto.Cfg.HTTP.Timeout,
		IdleTimeout:  dto.Cfg.HTTP.IdleTimeout,
	}

	err := srv.ListenAndServe()

	return err
}
