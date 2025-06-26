package handler

import (
	"github.com/go-chi/chi/v5"
	"img-compress/internal/app"
	"log/slog"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	Img    *app.Image
	Logger *slog.Logger
}

func NewImageHandler(logger *slog.Logger, img *app.Image) *ImageHandler {
	return &ImageHandler{Img: img, Logger: logger}
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.Logger.Error("GetImage: Invalid ID format", slog.Any("error", err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	path, err := h.Img.GetImage(id)

	if err != nil {
		h.Logger.Error("GetImage: Get image error", slog.Any("error", err))
		http.Error(w, "Img with id:"+strconv.Itoa(id)+" not found", http.StatusBadRequest)
		return
	}

	w.Write([]byte("GET img path: " + path))
}

func (h *ImageHandler) AddImages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /img called"))
}
