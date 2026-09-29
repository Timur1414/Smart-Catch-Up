package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error)
	DeleteByUuid(ctx context.Context, uuid string) error
	DeleteByUser(ctx context.Context, user domain.User) error
}

type RefreshTokenRedis struct {
	db *redis.Client
}

func NewRefreshTokenRedis(client *redis.Client) *RefreshTokenRedis {
	return &RefreshTokenRedis{db: client}
}

func (obj RefreshTokenRedis) Create(ctx context.Context, user domain.User) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	err := obj.db.Set(ctx, "uuid", "user_id", time.Hour).Err()
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to create jwt", zap.Error(err))
		return "", err
	}
	log.Info("Created jwt", zap.String("uuid", ""), zap.Duration("duration", duration))
	return "", nil
}

func (obj RefreshTokenRedis) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	userId, err := obj.db.Get(ctx, "uuid").Result()
	duration := time.Since(startTime)
	if errors.Is(err, redis.Nil) {
		return domain.RefreshToken{}, errors.New("uuid not found")
	} else if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Warn("Failed to convert jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	log.Info("Get jwt", zap.String("uuid", ""), zap.Duration("duration", duration))
	return domain.RefreshToken{
		Uuid:      uuid,
		UserId:    userIdInt,
		ExpiredAt: time.Time{},
	}, nil
}

func (obj RefreshTokenRedis) DeleteByUuid(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshTokenRedis) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
