package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/repository"
)

type SettingsUseCase interface {
	UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error
}

type Settings struct {
	repository repository.SettingsRepository
}

func NewSettings(repository repository.SettingsRepository) *Settings {
	return &Settings{repository: repository}
}

func (obj *Settings) UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error {
	return obj.repository.UpdateAvatar(ctx, userId, avatarUrl)
}
