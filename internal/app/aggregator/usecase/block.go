package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
)

type BlockUseCase interface {
	Aggregate(ctx context.Context, startUserID, endUserID int) error
}

type Block struct {
	repository       repository.BlockRepository
	notificationRepo repository.NotificationRepository
	digestRepo       repository.DigestRepository
	settingsRepo     repository.SettingsRepository
}

func NewBlock(repository repository.BlockRepository, notifRepo repository.NotificationRepository, digestRepo repository.DigestRepository, settingsRepo repository.SettingsRepository) *Block {
	return &Block{
		repository:       repository,
		notificationRepo: notifRepo,
		digestRepo:       digestRepo,
		settingsRepo:     settingsRepo,
	}
}

func (obj *Block) Aggregate(ctx context.Context, startUserID, endUserID int) error {
	log := logger.GetLoggerWithRequestId(ctx)
	log.Info("Start aggregation")
	allNotifications, err := obj.notificationRepo.GetAllByUserRange(ctx, startUserID, endUserID)
	if err != nil {
		return err
	}
	if len(allNotifications) == 0 {
		log.Info("No notifications to aggregate")
		return nil
	}
	userNotificationMap := make(map[int][]domain.Notification)
	for _, notification := range allNotifications {
		userNotificationMap[notification.RecipientId] = append(userNotificationMap[notification.RecipientId], notification)
	}

	userIds := make([]int, 0, len(userNotificationMap))
	for uId := range userNotificationMap {
		userIds = append(userIds, uId)
	}

	userSettingsMap, err := obj.settingsRepo.GetByUsers(ctx, userIds)
	if err != nil {
		log.Error("failed to get user settings", zap.Error(err))
		return err
	}

	for userId, userNotifications := range userNotificationMap {
		settings := userSettingsMap[userId]
		importantTypes := make(map[string]bool, len(settings.Important))
		for _, impType := range settings.Important {
			importantTypes[impType] = true
		}

		var importantNotifications []domain.Notification
		clusterGroups := make(map[string][]domain.Notification)

		for _, notification := range userNotifications {
			if importantTypes[notification.NotificationType] {
				importantNotifications = append(importantNotifications, notification)
			} else {
				cluster := GetClusterForNotificationType(notification.NotificationType)
				clusterGroups[cluster] = append(clusterGroups[cluster], notification)
			}
		}

		parts := make([]domain.BlockPart, 0, len(clusterGroups)+1)
		var overAllStart, overAllEnd time.Time

		for cluster, clusterNotifications := range clusterGroups {
			var clusterStart, clusterEnd time.Time
			var blockText []string
			notificationIds := make([]int, 0, len(clusterNotifications))
			for i, notification := range clusterNotifications {
				notificationIds = append(notificationIds, notification.Id)
				if i == 0 || notification.CreatedAt.Before(clusterStart) {
					clusterStart = notification.CreatedAt
				}
				if i == 0 || notification.CreatedAt.After(clusterEnd) {
					clusterEnd = notification.CreatedAt
				}
				blockText = append(blockText, notification.ActorName+": "+notification.Payload)
			}
			if overAllStart.IsZero() || clusterStart.Before(overAllStart) {
				overAllStart = clusterStart
			}
			if overAllEnd.IsZero() || clusterEnd.After(overAllEnd) {
				overAllEnd = clusterEnd
			}
			parts = append(parts, domain.BlockPart{
				ClusterType:     cluster,
				Text:            strings.Join(blockText, "\n"),
				Start:           clusterStart,
				End:             clusterEnd,
				NotificationIds: notificationIds,
			})
		}

		if len(importantNotifications) > 0 {
			var impStart, impEnd time.Time
			var impTexts []string
			impNotificationIds := make([]int, 0, len(importantNotifications))
			for i, notification := range importantNotifications {
				impNotificationIds = append(impNotificationIds, notification.Id)
				if i == 0 || notification.CreatedAt.Before(impStart) {
					impStart = notification.CreatedAt
				}
				if i == 0 || notification.CreatedAt.After(impEnd) {
					impEnd = notification.CreatedAt
				}
				impTexts = append(impTexts, notification.Payload)
			}
			if overAllStart.IsZero() || impStart.Before(overAllStart) {
				overAllStart = impStart
			}
			if overAllEnd.IsZero() || impEnd.After(overAllEnd) {
				overAllEnd = impEnd
			}
			parts = append(parts, domain.BlockPart{
				ClusterType:     "important",
				Text:            strings.Join(impTexts, "\n"),
				Start:           impStart,
				End:             impEnd,
				NotificationIds: impNotificationIds,
			})
		}

		if len(parts) == 0 {
			continue
		}
		block := domain.Block{
			UserId: userId,
			Start:  overAllStart,
			End:    overAllEnd,
			Status: "created",
			Parts:  parts,
		}
		blockId, err := obj.repository.Create(ctx, block)
		if err != nil {
			return err
		}
		log.Info("Successfully created Block and BlockParts",
			zap.Int("user_id", userId),
			zap.Int("block_id", blockId),
			zap.Int("parts_count", len(parts)),
		)
		err = obj.digestRepo.UpdateBlock(ctx, userId, blockId)
		if err != nil {
			return err
		}
		log.Info("Successfully updated digest")
	}
	log.Info("Finish aggregation")
	return nil
}
