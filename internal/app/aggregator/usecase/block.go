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
	Aggregate(ctx context.Context) error
}

type Block struct {
	repository       repository.BlockRepository
	notificationRepo repository.NotificationRepository
	digestRepo       repository.DigestRepository
}

func NewBlock(repository repository.BlockRepository, notifRepo repository.NotificationRepository, digestRepo repository.DigestRepository) *Block {
	return &Block{
		repository:       repository,
		notificationRepo: notifRepo,
		digestRepo:       digestRepo,
	}
}

func (obj *Block) Aggregate(ctx context.Context) error {
	log := logger.GetLoggerWithRequestId(ctx)
	log.Info("Start aggregation")
	allNotifications, err := obj.notificationRepo.GetAll(ctx)
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
	for userId, userNotifications := range userNotificationMap {
		clusterGroups := make(map[string][]domain.Notification)
		for _, notification := range userNotifications {
			cluster := GetClusterForNotificationType(notification.NotificationType)
			clusterGroups[cluster] = append(clusterGroups[cluster], notification)
		}
		parts := make([]domain.BlockPart, 0, len(clusterGroups))
		var overAllStart, overAllEnd time.Time
		for cluster, clusterNotifications := range clusterGroups {
			var clusterStart, clusterEnd time.Time
			var blockText []string
			for i, notification := range clusterNotifications {
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
				ClusterType: cluster,
				Text:        strings.Join(blockText, "\n"),
				Start:       clusterStart,
				End:         clusterEnd,
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
