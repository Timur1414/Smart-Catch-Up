package usecase

import (
	"context"
	"time"

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
	notification := GenerateNotification(notificationType, text, userId, []int{})
	return obj.notificationRepository.Create(ctx, notification)
}

func (obj *Notification) BulkCreate(ctx context.Context, notificationTypes []string, number int, userIds []int) error {
	notifications := make([]domain.Notification, number)
	for i := range number {
		notifications[i] = domain.Notification{
			NotificationType: notificationTypes[i%len(notificationTypes)],
			RecipientId:      userIds[i%len(userIds)],
			ActorId:          1,
			ActorName:        "system",
			ObjectId:         0,
			ObjectType:       "",
			CreatedAt:        time.Now(),
			ReadAt:           time.Time{},
			Payload:          "random text",
			Actions:          nil,
		}
	}
	return obj.notificationRepository.BulkCreate(ctx, notifications)
}
