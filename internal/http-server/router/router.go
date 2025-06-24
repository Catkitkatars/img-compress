package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	r "img-compress/internal/http-server/router/routes"
)

const (
	get  = "GET"
	post = "POST"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	for _, route := range r.Routes {

		switch route.Method {
		case get:
			router.Get(route.Path, route.Handler)
		case post:
			router.Post(route.Path, route.Handler)
		}

	}

	return router
}
