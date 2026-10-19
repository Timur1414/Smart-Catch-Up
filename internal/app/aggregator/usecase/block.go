package usecase

import (
	"context"
	"slices"
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
	userIds := make([]int, endUserID-startUserID)
	for i := range userIds {
		userIds[i] = startUserID + i
	}
	settings, err := obj.settingsRepo.GetByUsers(ctx, userIds)
	if err != nil {
		return err
	}
	userNotificationMap := make(map[int][]domain.Notification)
	for _, notification := range allNotifications {
		userNotificationMap[notification.RecipientId] = append(userNotificationMap[notification.RecipientId], notification)
	}
	for userId, userNotifications := range userNotificationMap {
		clusterGroups := make(map[string][]domain.Notification)
		curUserSettings := settings[userId]
		for _, notification := range userNotifications {
			cluster := ""
			if slices.Contains(curUserSettings.Important, notification.NotificationType) {
				cluster = "important"
			} else {
				cluster = GetClusterForNotificationType(notification.NotificationType)
			}
			clusterGroups[cluster] = append(clusterGroups[cluster], notification)
		}
		parts := make([]domain.BlockPart, 0, len(clusterGroups))
		var overAllStart, overAllEnd time.Time
		for cluster, clusterNotifications := range clusterGroups {
			var clusterStart, clusterEnd time.Time
			notificationIds := make([]int, 0, len(clusterNotifications))
			for i, notification := range clusterNotifications {
				notificationIds = append(notificationIds, notification.Id)
				if i == 0 || notification.CreatedAt.Before(clusterStart) {
					clusterStart = notification.CreatedAt
				}
				if i == 0 || notification.CreatedAt.After(clusterEnd) {
					clusterEnd = notification.CreatedAt
				}
			}
			if overAllStart.IsZero() || clusterStart.Before(overAllStart) {
				overAllStart = clusterStart
			}
			if overAllEnd.IsZero() || clusterEnd.After(overAllEnd) {
				overAllEnd = clusterEnd
			}
			var textGenerator TextGeneratorStrategy
			if cluster == "important" {
				textGenerator = NewImportantTextGeneratorStrategy(clusterNotifications)
			} else {
				textGenerator = NewBaseTextGeneratorStrategy(clusterNotifications, cluster)
			}
			parts = append(parts, domain.BlockPart{
				ClusterType:     cluster,
				Text:            obj.generateBlockText(true, clusterNotifications),
				Start:           clusterStart,
				End:             clusterEnd,
				NotificationIds: notificationIds,
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

func (obj *Block) generateBlockText(important bool, notifications []domain.Notification) string {
	// ToDo important and no important strategies
	var blockText []string
	for _, notification := range notifications {
		blockText = append(blockText, notification.ActorName+": "+notification.Payload)
	}
	return strings.Join(blockText, "\n")
}

type TextGeneratorStrategy interface {
	Generate() string
}

type BaseTextGeneratorStrategy struct {
	notifications []domain.Notification
	cluster       string
}

func NewBaseTextGeneratorStrategy(notifications []domain.Notification, cluster string) *BaseTextGeneratorStrategy {
	return &BaseTextGeneratorStrategy{
		notifications: notifications,
		cluster:       cluster,
	}
}

func (obj *BaseTextGeneratorStrategy) Generate() string {
	// ToDo llm
	var blockText []string
	for _, notification := range obj.notifications {
		blockText = append(blockText, notification.ActorName+": "+notification.Payload)
	}
	return strings.Join(blockText, "\n")
}

type ImportantTextGeneratorStrategy struct {
	notifications []domain.Notification
}

func NewImportantTextGeneratorStrategy(notifications []domain.Notification) *ImportantTextGeneratorStrategy {
	return &ImportantTextGeneratorStrategy{
		notifications: notifications,
	}
}

func (obj *ImportantTextGeneratorStrategy) Generate() string {
	var blockText []string
	for _, notification := range obj.notifications {
		blockText = append(blockText, notification.ActorName+": "+notification.Payload)
	}
	return strings.Join(blockText, "\n")
}
