package router

import (
	"img-compress/internal/handler"
	"net/http"
)

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

func GetRoutes() []Route {
	ih := handler.NewImageHandler()

	return []Route{
		{
			Method:  "GET",
			Path:    "/img/{id}",
			Handler: handler.Wrap(ih.GetImage),
		},
		{
			Method:  "POST",
			Path:    "/img",
			Handler: handler.Wrap(ih.AddImages),
		},
	}
}
