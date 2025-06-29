package app

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"img-compress/internal/logger"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	logger *slog.Logger
}

func NewImage() *Image {
	return &Image{logger: logger.Logger}
}

func (i *Image) Process(f multipart.File, header multipart.FileHeader) (string, error) {
	name := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

	srcImg, err := imaging.Decode(f)
	if err != nil {
		// todo
		// Шаблон декоратор/враппер
		return "", fmt.Errorf("imaging.Decode: %w", err)
	}

	bg := imaging.New(srcImg.Bounds().Dx(), srcImg.Bounds().Dy(), color.White)
	img := imaging.Overlay(bg, srcImg, image.Pt(0, 0), 1.0)

	imgMarked, err := i.AddWaterMark(img, "assets/watermark.png")
	if err != nil {
		return "", fmt.Errorf("image.AddWaterMark: %w", err)
	}

	outPath := "assets/img/" + name + ".jpg"
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer outFile.Close()

	err = imaging.Encode(outFile, imgMarked, imaging.JPEG, imaging.JPEGQuality(30))
	if err != nil {
		return "", fmt.Errorf("imaging.Encode: %w", err)
	}

	return outPath, nil
}

func (i *Image) AddWaterMark(img image.Image, wmPath string) (image.Image, error) {
	m := "app.Image.AddWaterMark"

	watermark, err := imaging.Open(wmPath)
	if err != nil {
		i.logger.Error(m, "failed to open watermark image", slog.String("wmPath", wmPath), slog.Any("error", err))
		return nil, err
	}
	watermark = i.ResizeWaterMark(watermark, img)

	bw, bh := img.Bounds().Dx(), img.Bounds().Dy()
	ww, wh := watermark.Bounds().Dx(), watermark.Bounds().Dy()

	offset := image.Pt((bw-ww)/2, (bh-wh)/2)

	return imaging.Overlay(img, watermark, offset, 0.3), nil
}

func (i *Image) ResizeWaterMark(wm image.Image, base image.Image) image.Image {
	baseW := base.Bounds().Dx()
	baseH := base.Bounds().Dy()

	wmW := wm.Bounds().Dx()
	wmH := wm.Bounds().Dy()

	maxW := baseW
	maxH := baseH

	if wmW > maxW || wmH > maxH {
		return imaging.Resize(wm, maxW, 0, imaging.Lanczos)
	}
	return wm
}
