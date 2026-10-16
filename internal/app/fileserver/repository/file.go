package repository

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileRepository interface {
	StoreFile(ctx context.Context, data io.Reader, extension string) (string, error)
}

const (
	StaticDir = "static"
)

type FileLocalStore struct {
	dir string
}

func NewFileLocalStore() (*FileLocalStore, error) {
	err := os.MkdirAll(StaticDir, 0o755)
	if err != nil {
		return nil, domain.ErrFailedToCreateDir
	}
	return &FileLocalStore{dir: StaticDir}, nil
}

func (obj *FileLocalStore) StoreFile(ctx context.Context, data io.Reader, extension string) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	fileName := filepath.Base(uuid.New().String() + extension)
	filePath := filepath.Join(obj.dir, fileName)
	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = dst.Close(); err != nil {
			log.Error("Error when file close", zap.Error(err))
		}
	}()
	if _, err := io.Copy(dst, data); err != nil {
		log.Error("Error when file copy", zap.Error(err))
		_ = os.Remove(filePath)
		return "", err
	}
	return fileName, nil
}
