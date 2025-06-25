package app

import (
	"img-compress/internal/dto"
	"img-compress/internal/storage"
	"log/slog"
)

type Image struct {
	st  *storage.Storage
	log *slog.Logger
}

func NewImage(imgAppDto *dto.ImageApp) *Image {
	return &Image{st: imgAppDto.Storage, log: imgAppDto.Log}
}

func (i *Image) GetImage(id int) (string, error) {
	fn := "app.Image.GetImage"

	res, err := i.st.GetImage(id)

	if err != nil {
		i.log.Error(fn, "failed to get image", slog.Int("id", id), slog.Any("error", err))
		return "", err
	}

	return res, nil
}
