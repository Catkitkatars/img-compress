package app

import (
	"github.com/disintegration/imaging"
	"image"
	"img-compress/internal/dto"
	"img-compress/internal/storage"
	"log/slog"
)

type Image struct {
	st     *storage.Storage
	logger *slog.Logger
}

func NewImage(imgAppDto *dto.Config) *Image {
	return &Image{st: imgAppDto.Storage, logger: imgAppDto.Log}
}

func (i *Image) Get(id int) (string, error) {
	m := "app.Image.GetImage"

	res, err := i.st.GetImage(id)

	if err != nil {
		i.logger.Error(m, "failed to get image", slog.Int("id", id), slog.Any("error", err))
		return "", err
	}

	return res, nil
}
func (i *Image) Save(path string) (int, error) {
	m := "app.Image.SaveImage"

	id, err := i.st.SaveImage(path)

	if err != nil {
		i.logger.Error(m, "failed to save image", slog.String("path", path), slog.Any("error", err))
		return 0, err
	}

	return id, nil
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
