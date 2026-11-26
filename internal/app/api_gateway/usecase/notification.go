package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/config"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type NotificationUseCase interface {
	GetByUser(ctx context.Context, userId int) ([]domain.Notification, error)
	GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error)
	GetAllTypes(ctx context.Context) ([]string, error)
}

type Notification struct {
	repository repository.NotificationRepository
	cfg        *config.Config
}

func NewNotification(repo repository.NotificationRepository, cfg *config.Config) *Notification {
	return &Notification{
		repository: repo,
		cfg:        cfg,
	}
}

func (obj *Notification) GetByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	return obj.repository.GetByUser(ctx, userId, obj.cfg.Api.CountRawNotifications)
}

func (obj *Notification) GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	return obj.repository.GetAllByUser(ctx, userId)
}

func (obj *Notification) GetAllTypes(ctx context.Context) ([]string, error) {
	return obj.repository.GetAllTypes(ctx)
}
