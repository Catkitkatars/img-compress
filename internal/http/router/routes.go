package router

import (
	"encoding/json"
	"fmt"
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
			Handler: ih.GetImage,
		},
		{
			Method:  "POST",
			Path:    "/img",
			Handler: ih.AddImages,
		},
		{
			Method:  "GET",
			Path:    "/ping",
			Handler: wrap(ping),
		},
	}
}

// todo
// Теперь можно обернуть методы хенлера и получать ошибку или ответ.
// Обрабатывать ошибку в врапере

// Разворачиваение и сворачивание ошибок

type MyHandler func(*http.Request) (any, error)

func wrap(myHandler MyHandler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		//...
		rsBody, err := myHandler(req)

		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(resp).Encode(rsBody)
	}
}

func ping(*http.Request) (any, error) {
	return "pong", fmt.Errorf("ping")
}
