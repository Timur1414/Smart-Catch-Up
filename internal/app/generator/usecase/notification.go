package usecase

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/repository"
)

type NotificationUseCase interface {
	Create(ctx context.Context, notificationType string, text string, userId int) (int, error)
	BulkCreate(ctx context.Context, notificationTypes []string, number int, userIds []int) ([]int, error)
}

type Notification struct {
	repository repository.NotificationRepository
}

func NewNotification(repo repository.NotificationRepository) *Notification {
	return &Notification{repository: repo}
}

func (obj *Notification) Create(ctx context.Context, notificationType string, text string, userId int) (int, error) {
	notification := domain.Notification{
		NotificationType: notificationType,
		RecipientId:      userId,
		ActorId:          0,
		ActorName:        "system",
		ObjectId:         0,
		ObjectType:       "",
		CreatedAt:        time.Now(),
		ReadAt:           time.Time{},
		Payload:          text,
		Actions:          nil,
	}
	return obj.repository.Create(ctx, notification)
}

func (obj *Notification) BulkCreate(ctx context.Context, notificationTypes []string, number int, userIds []int) ([]int, error) {
	notifications := make([]domain.Notification, number)
	for i := range number {
		notifications[i] = domain.Notification{
			NotificationType: notificationTypes[i%len(notificationTypes)],
			RecipientId:      userIds[i%len(userIds)],
			ActorId:          0,
			ActorName:        "system",
			ObjectId:         0,
			ObjectType:       "",
			CreatedAt:        time.Now(),
			ReadAt:           time.Time{},
			Payload:          "random text",
			Actions:          nil,
		}
	}
	return obj.repository.BulkCreate(ctx, notifications)
}
