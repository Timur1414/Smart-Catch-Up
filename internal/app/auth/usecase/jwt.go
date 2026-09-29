package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JwtUseCase interface {
	CreatePair(ctx context.Context, user domain.User) (string, string, error)
	UpdateAccessToken(ctx context.Context, user domain.User) (string, error)
	GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error)
	DeleteByUuid(ctx context.Context, uuid string) error
	DeleteByUser(ctx context.Context, user domain.User) error
	CheckAccessToken(ctx context.Context, tokenStr string) (bool, int)
	CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int)
}

type Jwt struct {
	repository repository.RefreshTokenRepository
	mu         sync.RWMutex
	secret     []byte
	version    string
}

func NewJwt(repo repository.RefreshTokenRepository, secret string, version string) (*Jwt, error) {
	if len(secret) < 8 {
		return nil, errors.New("secret too short")
	}
	if version == "" {
		return nil, errors.New("version too short")
	}
	return &Jwt{
		repository: repo,
		secret:     []byte(secret),
		version:    version,
	}, nil
}

func (obj *Jwt) CreatePair(ctx context.Context, user domain.User) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *Jwt) UpdateAccessToken(ctx context.Context, user domain.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *Jwt) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	return obj.repository.GetByUuid(ctx, uuid)
}

func (obj *Jwt) DeleteByUuid(ctx context.Context, uuid string) error {
	return obj.repository.DeleteByUuid(ctx, uuid)
}

func (obj *Jwt) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (obj *Jwt) CheckAccessToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return obj.secret, nil
	})
	if err != nil {
		log.Warn("Token parse error", zap.Error(err))
		return false, -1
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn("invalid token")
		return false, -1
	}
	version, ok := claims["version"].(string)
	if !ok {
		log.Warn("invalid token (unable to claim version)")
		return false, -1
	}
	if version != obj.GetVersion() {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}
	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		log.Warn("invalid token (unable to claim user_id)")
		return false, -1
	}
	userId := int(userIdFloat)
	return true, userId
}

func (obj *Jwt) CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return obj.secret, nil
	})
	if err != nil {
		log.Warn("Token parse error", zap.Error(err))
		return false, -1
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Error("invalid token (unable to claim payload)")
		return false, -1
	}
	version, ok := claims["version"].(string)
	if !ok {
		log.Warn("invalid token (unable to claim version)")
		return false, -1
	}
	if version != obj.GetVersion() {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}
	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		log.Warn("invalid token (unable to claim user_id)")
		return false, -1
	}
	userId := int(userIdFloat)
	tokenId, ok := claims["id"].(string)
	if !ok {
		log.Error("invalid token (unable to claim token_id)")
		return false, -1
	}
	storedToken, err := obj.repository.GetByUuid(ctx, tokenId)
	if err != nil {
		return false, -1
	}
	if storedToken.UserId != userId {
		log.Error("invalid token (invalid userId)")
		return false, -1
	}
	return true, userId
}

func (obj *Jwt) GetVersion() string {
	obj.mu.RLock()
	defer obj.mu.RUnlock()
	return obj.version
}
