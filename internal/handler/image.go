package handler

import (
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"image"
	"image/color"
	"img-compress/internal/app"
	"img-compress/internal/dto"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ImageHandler struct {
	Img       *app.Image
	ConfigDto *dto.Config
	Logger    *slog.Logger
}

type FileResult struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

func NewImageHandler(configDto *dto.Config, img *app.Image) *ImageHandler {
	return &ImageHandler{
		Img:       img,
		ConfigDto: configDto,
		Logger:    configDto.Log,
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
					Error:    "не удалось открыть файл",
				}
				return
			}

			path, errorText := h.ProcessImage(file, *header)
			file.Close()

			if errorText != "" {
				h.Logger.Error(errorText, slog.String("filename", header.Filename), slog.Any("error", err))
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    errorText,
				}
				return
			}

			_, saveErr := h.Img.Save(path)

			if saveErr != nil {
				h.Logger.Error("AddImages: Save image error", slog.String("path", path), slog.Any("error", err))
				resultCh <- FileResult{
					Filename: header.Filename,
					Status:   "error",
					Error:    "AddImages: Save image error",
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

func (h *ImageHandler) ProcessImage(f multipart.File, header multipart.FileHeader) (string, string) {
	name := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

	srcImg, err := imaging.Decode(f)
	if err != nil {
		return "", "Не удалось декодировать изображение: " + err.Error()
	}

	bg := imaging.New(srcImg.Bounds().Dx(), srcImg.Bounds().Dy(), color.White)
	img := imaging.Overlay(bg, srcImg, image.Pt(0, 0), 1.0)

	imgMarked, err := h.Img.AddWaterMark(img, "assets/watermark.png")
	if err != nil {
		return "", "Не удалось добавить водяной знак: " + err.Error()
	}

	outPath := "assets/img/" + name + ".jpg"
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", "Не удалось создать файл для сохранения изображения: " + err.Error()
	}
	defer outFile.Close()

	err = imaging.Encode(outFile, imgMarked, imaging.JPEG, imaging.JPEGQuality(30))
	if err != nil {
		return "", "Не удалось сохранить изображение: " + err.Error()
	}

	return outPath, ""
}
