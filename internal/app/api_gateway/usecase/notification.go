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

func (obj Notification) Create(ctx context.Context, notification domain.Notification) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Notification) GetById(ctx context.Context, id int) (domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Notification) GetByUser(ctx context.Context, user domain.User) (domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Notification) Update(ctx context.Context, notification domain.Notification) error {
	//TODO implement me
	panic("implement me")
}
