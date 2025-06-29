package handler

import (
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"image"
	"image/color"
	"img-compress/internal/app"
	"img-compress/internal/dto"
	"img-compress/internal/storage"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ImageHandler struct {
	store  *storage.Storage
	Logger *slog.Logger
}

type FileResult struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Error    error  `json:"error,omitempty"`
}

func NewImageHandler() *ImageHandler {

	return &ImageHandler{
		Logger: configDto.Log,
	}
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.Logger.Error("GetImage: Invalid ID format", slog.Any("error", err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	path, err := h.Img.Get(id)

	if err != nil {
		h.Logger.Error("GetImage: Get image error", slog.Any("error", err))
		http.Error(w, "Img with id:"+strconv.Itoa(id)+" not found", http.StatusBadRequest)
		return
	}

	w.Write([]byte("GET img path: " + path))
}

func (h *ImageHandler) AddImages(w http.ResponseWriter, r *http.Request) {
	var results []FileResult
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Не удалось распарсить multipart-форму", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["img"]
	if len(files) == 0 {
		http.Error(w, "Файлы не найдены", http.StatusBadRequest)
		return
	}

	resultCh := make(chan FileResult)

	for _, header := range files {
		go func(header *multipart.FileHeader) {
			file, err := header.Open()
			if err != nil {
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    fmt.Errorf("header.Open: %w", err),
				}
				return
			}

			path, processError := h.Img.Process(file, *header)
			file.Close()

			if processError != nil {
				h.Logger.Error("h.Img.Process: ", slog.String("filename", header.Filename), slog.Any("error", err))
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    processError,
				}
				return
			}

			_, saveErr := h.Img.Save(path)

			if saveErr != nil {
				h.Logger.Error("AddImages: Save image error", slog.String("path", path), slog.Any("error", err))
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    fmt.Errorf("h.Img.Save: %w", saveErr),
				}
				return
			}

			resultCh <- FileResult{
				Filename: header.Filename,
				Status:   "success",
			}
			return
		}(header)
	}

	for _ = range files {
		results = append(results, <-resultCh)
	}

	close(resultCh)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
