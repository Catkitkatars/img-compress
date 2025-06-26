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

func GetRoutes(imgHandler *handler.ImageHandler) []Route {
	return []Route{
		{
			Method:  "GET",
			Path:    "/img/{id}",
			Handler: imgHandler.GetImage,
		},
		{
			Method:  "POST",
			Path:    "/img",
			Handler: imgHandler.AddImages,
		},
	}
}
