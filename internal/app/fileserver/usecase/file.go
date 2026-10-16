package usecase

import (
	"context"
	"io"
	"slices"
	"strings"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/repository"
)

type FileUseCase interface {
	StoreFile(ctx context.Context, data io.Reader, extension string) (string, error)
}

type FileUseCaseImplementation struct {
	repository repository.FileRepository
}

func NewFileUseCaseImplementation(repository repository.FileRepository) *FileUseCaseImplementation {
	return &FileUseCaseImplementation{repository: repository}
}

func (obj *FileUseCaseImplementation) StoreFile(ctx context.Context, data io.Reader, extension string) (string, error) {
	if data == nil {
		return "", domain.ErrEmptyName
	}
	ext := strings.ToLower(strings.TrimSpace(extension))
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	allowed := []string{".jpg", ".jpeg", ".png", ".bmp", ".webp"}
	if !slices.Contains(allowed, ext) {
		return "", domain.ErrInvalidExtension
	}
	return obj.repository.StoreFile(ctx, data, extension)
}
