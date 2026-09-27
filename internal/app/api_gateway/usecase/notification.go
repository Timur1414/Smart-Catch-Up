package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type NotificationUseCase interface {
	Create(ctx context.Context, notification domain.Notification) (int, error)
	GetById(ctx context.Context, id int) (domain.Notification, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Notification, error)
	Update(ctx context.Context, notification domain.Notification) error
}

type Notification struct {
	repository repository.NotificationRepository
}

func NewNotification(repo repository.NotificationRepository) *Notification {
	return &Notification{repository: repo}
}
