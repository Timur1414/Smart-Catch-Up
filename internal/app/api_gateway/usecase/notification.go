package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type NotificationUseCase interface {
	GetById(ctx context.Context, id int) (domain.Notification, error)
	GetByUser(ctx context.Context, userId int) ([]domain.Notification, error)
	GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error)
	GetAllTypes(ctx context.Context) ([]string, error)
}

type Notification struct {
	repository repository.NotificationRepository
}

func NewNotification(repo repository.NotificationRepository) *Notification {
	return &Notification{repository: repo}
}

func (obj *Notification) GetById(ctx context.Context, id int) (domain.Notification, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj *Notification) GetByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	return obj.repository.GetByUser(ctx, userId, 10)
}

func (obj *Notification) GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	return obj.repository.GetAllByUser(ctx, userId)
}

func (obj *Notification) GetAllTypes(ctx context.Context) ([]string, error) {
	return obj.repository.GetAllTypes(ctx)
}
