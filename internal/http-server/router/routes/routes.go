package routes

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

var Routes = []Route{
	{
		Method: "GET",
		Path:   "/img",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("GET /img called"))
		},
	},
	{
		Method: "POST",
		Path:   "/img",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("POST /img called"))
		},
	},
}
