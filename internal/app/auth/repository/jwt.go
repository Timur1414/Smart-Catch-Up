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
	Create(ctx context.Context, token domain.RefreshToken) (string, error)
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

func (obj *RefreshTokenRedis) Create(ctx context.Context, token domain.RefreshToken) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	err := obj.db.Set(ctx, token.Uuid, token.UserId, time.Hour).Err() // TODO swap key and value
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to create jwt", zap.Error(err))
		return "", err
	}
	log.Info("Created jwt", zap.String("uuid", token.Uuid), zap.Duration("duration", duration))
	return token.Uuid, nil
}

func (obj *RefreshTokenRedis) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	var userIdCmd *redis.StringCmd
	var ttlCmd *redis.DurationCmd
	startTime := time.Now()
	_, err := obj.db.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		userIdCmd = pipe.Get(ctx, uuid)
		ttlCmd = pipe.TTL(ctx, uuid)
		return nil
	})
	duration := time.Since(startTime)
	if errors.Is(err, redis.Nil) {
		return domain.RefreshToken{}, errors.New("uuid not found")
	} else if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	userId, err := userIdCmd.Result()
	if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	ttl, err := ttlCmd.Result()
	if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Warn("Failed to convert jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	expiredAt := time.Now().Add(ttl)
	log.Info("Get jwt", zap.String("uuid", uuid), zap.Duration("duration", duration))
	return domain.RefreshToken{
		Uuid:      uuid,
		UserId:    userIdInt,
		ExpiredAt: expiredAt,
	}, nil
}

func (obj *RefreshTokenRedis) DeleteByUuid(ctx context.Context, uuid string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	err := obj.db.Del(ctx, uuid).Err()
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to delete jwt", zap.Error(err))
		return err
	}
	log.Info("Delete jwt", zap.String("uuid", uuid), zap.Duration("duration", duration))
	return nil
}

func (obj *RefreshTokenRedis) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
