package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"img-compress/internal/app"
	"img-compress/internal/logger"
	"img-compress/internal/storage"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	Store  *storage.Storage
	Logger *slog.Logger
	Image  *app.Image
}

type FileResult struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Error    error  `json:"error,omitempty"`
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		Store:  storage.Store,
		Logger: logger.Logger,
		Image:  app.NewImage(),
	}
}

func (h *ImageHandler) GetImage(r *http.Request) (any, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		return "", fmt.Errorf("strconv.Atoi: %w", err)
	}

	path, err := h.Store.GetImage(id)

	if err != nil {
		return "", fmt.Errorf("h.Store.GetImage: %w", err)
	}

	return path, nil
}

func (h *ImageHandler) AddImages(r *http.Request) (any, error) {
	var results []FileResult
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return results, err
	}
	files := r.MultipartForm.File["img"]
	if len(files) == 0 {
		return results, fmt.Errorf("files not found")
	}

	resultCh := make(chan FileResult)

	for _, header := range files {
		go func(header *multipart.FileHeader) {
			file, err := header.Open()
			if err != nil {
				h.Logger.Error("header.Open: ", slog.String("filename", header.Filename), slog.Any("error", err))
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    fmt.Errorf("header.Open: %w", err),
				}
				return
			}

			path, processError := h.Image.Process(file, *header)
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

			_, saveErr := h.Store.SaveImage(path)

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
	return results, nil
}
