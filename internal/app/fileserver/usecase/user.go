package usecase

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SettingsUseCase interface {
	UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error
	StoreFile(ctx context.Context, data io.Reader, extension string) (string, error)
}

const (
	dir = "static"
)

type Settings struct {
	repository repository.SettingsRepository
}

func NewSettings(repository repository.SettingsRepository) *Settings {
	return &Settings{repository: repository}
}

func (obj *Settings) UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error {
	return obj.repository.UpdateAvatar(ctx, userId, avatarUrl)
}

func (obj *Settings) StoreFile(ctx context.Context, data io.Reader, extension string) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := filepath.Base(uuid.New().String() + extension)
	if name == "" || name == "." || name == "/" {
		return "", domain.ErrEmptyName
	}
	dst, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = dst.Close(); err != nil {
			log.Error("Error when file close", zap.Error(err))
		}
	}()
	if _, err := io.Copy(dst, data); err != nil {
		_ = os.Remove(filepath.Join(dir, name))
		return "", err
	}
	return name, nil
}
