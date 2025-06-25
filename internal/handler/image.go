package handler

import (
	"github.com/go-chi/chi/v5"
	"img-compress/internal/app"
	"img-compress/internal/dto"
	"log/slog"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	Img    *app.Image
	Logger *slog.Logger
}

func NewImageHandler(img *app.Image, dto *dto.ImageApp) *ImageHandler {
	return &ImageHandler{Img: img, Logger: dto.Log}
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	//
	//path, err := h.Img.GetImage(id)
	//
	//if err != nil {
	//	http.Error(w, "Img with id:"+strconv.Itoa(id)+" not found", http.StatusBadRequest)
	//	return
	//}
	//
	//w.Write([]byte("GET img path: " + path))
	w.Write([]byte("GET img path: " + strconv.Itoa(id)))
}

func (h *ImageHandler) AddImages(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /img called"))
}
