package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/repository"
)

type NotificationUseCase interface {
	Create(ctx context.Context, notificationType string, text string, userId int) (int, error)
	BulkCreate(ctx context.Context, notificationTypes []string, number int, userIds []int) error
}

type Notification struct {
	notificationRepository repository.NotificationRepository
	userRepository         repository.UserRepository
}

func NewNotification(repo repository.NotificationRepository, userRepo repository.UserRepository) *Notification {
	return &Notification{
		notificationRepository: repo,
		userRepository:         userRepo,
	}
}

func (obj *Notification) Create(ctx context.Context, notificationType string, text string, userId int) (int, error) {
	allowedActors, err := obj.userRepository.GetActors(ctx, []int{userId})
	if err != nil {
		return 0, err
	}
	notification := GenerateNotification(notificationType, text, userId, allowedActors)
	return obj.notificationRepository.Create(ctx, notification)
}

func (obj *Notification) BulkCreate(ctx context.Context, notificationTypes []string, number int, userIds []int) error {
	allowedActors, err := obj.userRepository.GetActors(ctx, userIds)
	if err != nil {
		return err
	}
	notifications := make([]domain.Notification, number)
	for i := range number {
		notifications[i] = GenerateNotification(notificationTypes[i%len(notificationTypes)], "", userIds[i%len(userIds)], allowedActors)
	}
	return obj.notificationRepository.BulkCreate(ctx, notifications)
}
