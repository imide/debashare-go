package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

type FileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

type Storage interface {
	Upload(filename string, content io.Reader) (*FileInfo, error)

	Download(filename string) (io.ReadCloser, error)

	Delete(filename string) error

	List() ([]FileInfo, error)
}

type LocalStorage struct {
	basePath string
}

// Upload implements Storage.Upload
func (ls *LocalStorage) Upload(filename string, content io.Reader) (*FileInfo, error) {
	// Create a safe filename
	safeFilename := filepath.Clean(filename)
	fullPath := filepath.Join(ls.basePath, safeFilename)

	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
			return
		}
	}(file)

	// Copy content to file
	written, err := io.Copy(file, content)
	if err != nil {
		err := os.Remove(fullPath)
		if err != nil {
			return nil, err
		} // Clean up on error
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &FileInfo{
		Name:      safeFilename,
		Size:      written,
		CreatedAt: time.Now(),
	}, nil
}

// Download implements Storage.Download
func (ls *LocalStorage) Download(filename string) (io.ReadCloser, error) {
	safeFilename := filepath.Clean(filename)
	fullPath := filepath.Join(ls.basePath, safeFilename)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete implements Storage.Delete
func (ls *LocalStorage) Delete(filename string) error {
	safeFilename := filepath.Clean(filename)
	fullPath := filepath.Join(ls.basePath, safeFilename)

	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// List implements Storage.List
func (ls *LocalStorage) List() ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(ls.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:      info.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return files, nil
}
