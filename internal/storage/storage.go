package storage

import (
	"database/sql"
	"fmt"
	slogger "img-compress/internal/logger"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type StorageInterface interface {
	SaveImage(path string) error
	GetImage(id int) (string, error)
	GetAllImages() ([]string, error)
	Close() error
}

type Storage struct {
	db     *sql.DB
	logger *slog.Logger
}

var Store *Storage

func New(storagePath string) error {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {

		return fmt.Errorf("%s: %w", fn, err)
	}

	logger := slogger.Logger

	Store := &Storage{db: db, logger: logger}

	if err := Store.initSchema(); err != nil {
		return fmt.Errorf("%s: initSchema: %w", fn, err)
	}

	return nil
}

func (s *Storage) initSchema() error {
	stmt, err := s.db.Prepare(`
		CREATE TABLE IF NOT EXISTS images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	return err
}

func (s *Storage) SaveImage(path string) (int, error) {
	mt := "storage.sqlite.SaveImage"

	r, execErr := s.db.Exec(`INSERT INTO images (path) VALUES (?)`, path)

	if execErr != nil {
		s.logger.Error(mt, "failed to save image", slog.String("path", path), slog.Any("error", execErr))
		return 0, execErr
	}

	id, err := r.LastInsertId()

	if err != nil {
		s.logger.Error(mt, "failed to get last insert id", slog.String("path", path), slog.Any("error", err))
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetImage(id int) (string, error) {
	mt := "storage.sqlite.GetImage"

	var path string
	err := s.db.QueryRow("SELECT path FROM images WHERE id = ?", id).Scan(&path)

	if err != nil {
		s.logger.Error(mt, "failed to get image", slog.Int("id", id), slog.Any("error", err))
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("image with id %d not found", id)
		}
		return "", err
	}
	return path, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
