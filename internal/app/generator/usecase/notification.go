package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/repository"
)

type NotificationUseCase interface {
	Create(ctx context.Context, notificationType string, text string, userId int) (int, error)
	BulkCreate(ctx context.Context, notificationType []string, number int, userId []int) ([]int, error)
}

type Notification struct {
	repository repository.NotificationRepository
}

func NewNotification(repo repository.NotificationRepository) *Notification {
	return &Notification{repository: repo}
}

func (obj *Notification) Create(ctx context.Context, notificationType string, text string, userId int) (int, error) {
	notification := domain.Notification{
		Id:               0,
		NotificationType: notificationType,
		RecipientId:      userId,
		ActorId:          0,
		ActorName:        "",
		ObjectId:         0,
		ObjectType:       "",
		CreatedAt:        time.Now(),
		ReadAt:           time.Time{},
		Payload:          "random text",
		Actions:          nil,
	}
	return obj.repository.Create(ctx, notification)
}

func (obj *Notification) BulkCreate(ctx context.Context, notificationType []string, number int, userId []int) ([]int, error) {
	if len(notificationType) != number || len(userId) != number {
		return nil, errors.New("invalid arguments len")
	}
	notifications := make([]domain.Notification, number)
	for i := range number {
		notifications[i] = domain.Notification{
			Id:               0 + i,
			NotificationType: notificationType[i],
			RecipientId:      userId[i],
			ActorId:          0,
			ActorName:        "",
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
