package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"img-compress/internal/handler"
)

const (
	get  = "GET"
	post = "POST"
)

func Init(imgHandler *handler.ImageHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	routes := GetRoutes(imgHandler)

	for _, route := range routes {

		switch route.Method {
		case get:
			router.Get(route.Path, route.Handler)
		case post:
			router.Post(route.Path, route.Handler)
		}

	}

	return router
}
