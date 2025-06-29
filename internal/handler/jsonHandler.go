package handler

import (
	"encoding/json"
	log "img-compress/internal/logger"
	"net/http"
)

type JsonHandler func(*http.Request) (any, error)

func Wrap(handler JsonHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		rsBody, err := handler(req)

		if err != nil {
			log.Logger.Error("JsonHandler.Wrap: ", err)
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err := json.NewEncoder(res).Encode(rsBody); err != nil {
			http.Error(res, "Failed to encode response", http.StatusBadRequest)
		}
	}
}
