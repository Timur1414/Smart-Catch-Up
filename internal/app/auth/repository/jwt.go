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

type TwoFactorTokenRepository interface {
	Create(ctx context.Context, twoFactorToken domain.RefreshToken) error
	GetByUser(ctx context.Context, userId int) (domain.RefreshToken, error)
	Delete(ctx context.Context, userId int) error
	IncrementAttempts(ctx context.Context, userId int) (int, error)
	GetUsed2FACode(ctx context.Context, userId int) (string, error)
	AddUsed2FACode(ctx context.Context, userId int, code string) error
}

const (
	refreshTokenPrefix           = "auth:refresh:"
	AccessTokenExpirationTime    = time.Minute * 5
	RefreshTokenExpirationTime   = time.Hour * 24 * 7
	twoFactorFlowPrefix          = "auth:2fa:flow:"
	twoFactorAttemptsPrefix      = "auth:2fa:attempts:"
	TwoFactorTokenExpirationTime = time.Minute * 3
	usedCodePrefix               = "auth:2fa:used_code:"
	UsedCodeExpirationTime       = time.Second * 30
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
	key := refreshTokenPrefix + strconv.Itoa(token.UserId)
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
	key := refreshTokenPrefix + strconv.Itoa(userId)
	startTime := time.Now()
	_, err := obj.db.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		uuidCmd = pipe.Get(ctx, key)
		ttlCmd = pipe.TTL(ctx, key)
		return nil
	})
	duration := time.Since(startTime)
	if errors.Is(err, redis.Nil) {
		return domain.RefreshToken{}, domain.ErrNotFound
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
	expiredAt := time.Now().Add(ttl).UTC()
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
	key := refreshTokenPrefix + strconv.Itoa(userId)
	err := obj.db.Del(ctx, key).Err()
	duration := time.Since(startTime)
	if err != nil {
		log.Warn("Failed to delete jwt", zap.Error(err))
		return err
	}
	log.Info("Delete jwt", zap.Duration("duration", duration))
	return nil
}

type TwoFactorTokenRedis struct {
	db *redis.Client
}

func NewTwoFactorTokenRedis(client *redis.Client) *TwoFactorTokenRedis {
	return &TwoFactorTokenRedis{db: client}
}

func (obj *TwoFactorTokenRedis) Create(ctx context.Context, twoFactorToken domain.RefreshToken) error {
	log := logger.GetLoggerWithRequestId(ctx)
	startTime := time.Now()
	key := twoFactorFlowPrefix + strconv.Itoa(twoFactorToken.UserId)
	err := obj.db.Set(ctx, key, twoFactorToken.Uuid, TwoFactorTokenExpirationTime).Err()
	if err != nil {
		log.Warn("Failed to create temp token", zap.Error(err))
		return err
	}
	duration := time.Since(startTime)
	log.Info("Created temp token", zap.String("uuid", twoFactorToken.Uuid), zap.Duration("duration", duration))
	return nil
}

func (obj *TwoFactorTokenRedis) GetByUser(ctx context.Context, userId int) (domain.RefreshToken, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	var tokenCmd *redis.StringCmd
	var ttlCmd *redis.DurationCmd
	key := twoFactorFlowPrefix + strconv.Itoa(userId)
	startTime := time.Now()
	_, err := obj.db.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		tokenCmd = pipe.Get(ctx, key)
		ttlCmd = pipe.TTL(ctx, key)
		return nil
	})
	duration := time.Since(startTime)
	if errors.Is(err, redis.Nil) {
		log.Warn("temp token not found or expired", zap.Error(err))
		return domain.RefreshToken{}, domain.ErrNotFound
	} else if err != nil {
		log.Warn("Failed to get temp token", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	tokenUuid, err := tokenCmd.Result()
	if err != nil {
		log.Warn("Failed to get temp token", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	ttl, err := ttlCmd.Result()
	if err != nil {
		log.Warn("Failed to get temp token", zap.Error(err))
		return domain.RefreshToken{}, err
	}
	expiredAt := time.Now().Add(ttl).UTC()
	log.Info("Get temp token", zap.String("uuid", tokenUuid), zap.Duration("duration", duration))
	return domain.RefreshToken{
		Uuid:      tokenUuid,
		UserId:    userId,
		ExpiredAt: expiredAt,
	}, nil
}

func (obj *TwoFactorTokenRedis) Delete(ctx context.Context, userId int) error {
	log := logger.GetLoggerWithRequestId(ctx)
	userIdStr := strconv.Itoa(userId)
	tokenKey := twoFactorFlowPrefix + userIdStr
	attemptsKey := twoFactorAttemptsPrefix + userIdStr
	startTime := time.Now()
	err := obj.db.Del(ctx, tokenKey, attemptsKey).Err()
	if err != nil {
		log.Warn("Failed to delete temp token", zap.Error(err))
		return err
	}
	duration := time.Since(startTime)
	log.Info("Delete temp token", zap.String("uuid", tokenKey), zap.Duration("duration", duration))
	return nil
}

func (obj *TwoFactorTokenRedis) IncrementAttempts(ctx context.Context, userId int) (int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	key := twoFactorAttemptsPrefix + strconv.Itoa(userId)
	startTime := time.Now()
	attempts, err := obj.db.Incr(ctx, key).Result()
	if err != nil {
		log.Warn("Failed to increment temp token", zap.Error(err))
		return 0, err
	}
	if attempts == 1 {
		err = obj.db.Expire(ctx, key, TwoFactorTokenExpirationTime).Err()
		if err != nil {
			log.Warn("Failed to expire temp token", zap.Error(err))
		}
	}
	duration := time.Since(startTime)
	log.Info("Increment temp token", zap.String("uuid", key), zap.Duration("duration", duration), zap.Int64("attempts", attempts))
	return int(attempts), err
}

func (obj *TwoFactorTokenRedis) GetUsed2FACode(ctx context.Context, userId int) (string, error) {
	// ToDo implement
	panic("implement me")
}

func (obj *TwoFactorTokenRedis) AddUsed2FACode(ctx context.Context, userId int, code string) error {
	// ToDo implement
	panic("implement me")
}
