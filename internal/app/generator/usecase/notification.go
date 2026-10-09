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
	settingsRepository     repository.SettingsRepository
}

func NewNotification(repo repository.NotificationRepository, userRepo repository.UserRepository, settingsRepo repository.SettingsRepository) *Notification {
	return &Notification{
		notificationRepository: repo,
		userRepository:         userRepo,
		settingsRepository:     settingsRepo,
	}
}

func (obj *Notification) Create(ctx context.Context, notificationType string, text string, userId int) (int, error) {
	allowedActors, err := obj.userRepository.GetActors(ctx, []int{userId})
	if err != nil {
		return 0, err
	}
	notification := GenerateNotification(notificationType, text, userId, allowedActors)
	if len(notification.Actions) > 0 {
		curUserSettings, err := obj.settingsRepository.GetByUser(ctx, userId)
		if err != nil {
			return 0, err
		}
		for i := range notification.Actions {
			if notification.Actions[i].ActionType == "text" {
				notification.Actions[i].ActionTarget = "@" + curUserSettings.FirstName
			}
		}
	}
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
		if len(notifications[i].Actions) > 0 {
			curUserSettings, err := obj.settingsRepository.GetByUser(ctx, userIds[i%len(userIds)])
			if err != nil {
				return err
			}
			for j := range notifications[i].Actions {
				if notifications[i].Actions[j].ActionType == "text" {
					notifications[i].Actions[j].ActionTarget = "@" + curUserSettings.FirstName
				}
			}
		}
	}
	return obj.notificationRepository.BulkCreate(ctx, notifications)
}
