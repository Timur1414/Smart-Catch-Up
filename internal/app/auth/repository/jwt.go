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
	GetByUser(ctx context.Context, userId int) (domain.RefreshToken, error)
	DeleteByUser(ctx context.Context, userId int) error
}

const (
	AccessTokenExpirationTime  = time.Minute * 5
	RefreshTokenExpirationTime = time.Hour * 24 * 7
)

type RefreshTokenRedis struct {
	db *redis.Client
}

func NewRefreshTokenRedis(client *redis.Client) *RefreshTokenRedis {
	return &RefreshTokenRedis{db: client}
}

func (obj *RefreshTokenRedis) Create(ctx context.Context, token domain.RefreshToken) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	key := "auth:" + strconv.Itoa(token.UserId)
	err := obj.db.Set(ctx, key, token.Uuid, RefreshTokenExpirationTime).Err()
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to create jwt", zap.Error(err))
		return "", err
	}
	log.Info("Created jwt", zap.String("uuid", token.Uuid), zap.Duration("duration", duration))
	return token.Uuid, nil
}

func (obj *RefreshTokenRedis) GetByUser(ctx context.Context, userId int) (domain.RefreshToken, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	var uuidCmd *redis.StringCmd
	var ttlCmd *redis.DurationCmd
	key := "auth:" + strconv.Itoa(userId)
	startTime := time.Now()
	_, err := obj.db.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		uuidCmd = pipe.Get(ctx, key)
		ttlCmd = pipe.TTL(ctx, key)
		return nil
	})
	duration := time.Since(startTime)
	if errors.Is(err, redis.Nil) {
		return domain.RefreshToken{}, errors.New("uuid not found")
	} else if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	uuid, err := uuidCmd.Result()
	if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	ttl, err := ttlCmd.Result()
	if err != nil {
		log.Warn("Failed to get jwt", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	expiredAt := time.Now().Add(ttl)
	log.Info("Get jwt", zap.String("uuid", uuid), zap.Duration("duration", duration))
	return domain.RefreshToken{
		Uuid:      uuid,
		UserId:    userId,
		ExpiredAt: expiredAt,
	}, nil
}

func (obj *RefreshTokenRedis) DeleteByUser(ctx context.Context, userId int) error {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	key := "auth:" + strconv.Itoa(userId)
	err := obj.db.Del(ctx, key).Err()
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to delete jwt", zap.Error(err))
		return err
	}
	log.Info("Delete jwt", zap.Duration("duration", duration))
	return nil
}
