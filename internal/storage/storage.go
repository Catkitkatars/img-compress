package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type StorageInterface interface {
	SaveImage(path string) error
	GetImage(id int) (string, error)
	GetAllImages() ([]string, error)
	Close() error
}

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	s := &Storage{db: db}
	if err := s.initSchema(); err != nil {
		return nil, fmt.Errorf("%s: initSchema: %w", fn, err)
	}

	return s, nil
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

func (s *Storage) SaveImage(path string) error {
	_, err := s.db.Exec(`INSERT INTO images (path) VALUES (?)`, path)
	return err
}

func (s *Storage) GetImage(id int) (string, error) {
	var path string
	err := s.db.QueryRow("SELECT path FROM images WHERE id = ?", id).Scan(&path)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("image with id %d not found", id)
		}
		return "", err
	}
	return path, nil
}

func (s *Storage) GetAllImages() ([]string, error) {
	rows, err := s.db.Query(`SELECT path FROM images`)
	if err != nil {
		return nil, err
	}

	var images []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		images = append(images, path)
	}
	return images, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
