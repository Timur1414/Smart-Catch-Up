package usecase

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JwtUseCase interface {
	CreatePair(ctx context.Context, userId int) (string, string, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	CheckAccessToken(ctx context.Context, tokenStr string) (bool, int)
	CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int)
	ExtractClaims(ctx context.Context, tokenStr string) (jwt.RegisteredClaims, error)
}

type Jwt struct {
	repository repository.RefreshTokenRepository
	mu         sync.RWMutex
	secret     []byte
	version    string
}

func NewJwt(repo repository.RefreshTokenRepository, secret string, version string) (*Jwt, error) {
	if len(secret) < 8 {
		return nil, domain.ErrSecretTooShort
	}
	if version == "" {
		return nil, domain.ErrVersionIsEmpty
	}
	return &Jwt{
		repository: repo,
		secret:     []byte(secret),
		version:    version,
	}, nil
}

func (obj *Jwt) CreatePair(ctx context.Context, userId int) (string, string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userIdStr := strconv.Itoa(userId)

	err := obj.repository.DeleteByUser(ctx, userId)
	if err != nil {
		return "", "", err
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(repository.AccessTokenExpirationTime).UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    obj.GetVersion(),
		Subject:   userIdStr,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(obj.GetSecret())
	if err != nil {
		log.Error("failed to sign access token",
			zap.Int("user_id", userId),
			zap.Error(err))
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(repository.RefreshTokenExpirationTime).UTC()
	refreshId := uuid.New().String()
	refreshClaims := jwt.RegisteredClaims{
		ID:        refreshId,
		ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    obj.GetVersion(),
		Subject:   userIdStr,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(obj.GetSecret())
	if err != nil {
		log.Error("failed to sign refresh token",
			zap.Int("user_id", userId),
			zap.Error(err))
		return "", "", err
	}

	refreshTokenToSave := domain.RefreshToken{
		Uuid:      refreshId,
		UserId:    userId,
		ExpiredAt: refreshExpirationTime,
	}
	_, err = obj.repository.Create(ctx, refreshTokenToSave)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (obj *Jwt) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, refreshToken)
	if err != nil {
		return err
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Warn("failed to convert subject to int", zap.Error(err))
		return err
	}
	return obj.repository.DeleteByUser(ctx, userId)
}

func (obj *Jwt) CheckAccessToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, tokenStr)
	if err != nil {
		return false, -1
	}

	if claims.Issuer != obj.GetVersion() {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return false, -1
	}
	return true, userId
}

func (obj *Jwt) CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, tokenStr)
	if err != nil {
		return false, -1
	}

	if claims.Issuer != obj.GetVersion() {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return false, -1
	}

	storedToken, err := obj.repository.GetByUser(ctx, userId)
	if err != nil {
		return false, -1
	}
	if storedToken.ExpiredAt.Before(time.Now().UTC()) {
		log.Warn("token is expired")
		return false, -1
	}
	if storedToken.Uuid != claims.ID {
		log.Warn("invalid token (invalid token Id)")
		return false, -1
	}
	return true, userId
}

func (obj *Jwt) ExtractClaims(ctx context.Context, tokenStr string) (jwt.RegisteredClaims, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return obj.GetSecret(), nil
	})
	if err != nil || !token.Valid {
		log.Warn("Token parse error", zap.Error(err))
		return jwt.RegisteredClaims{}, err
	}
	return claims, nil
}

func (obj *Jwt) GetVersion() string {
	obj.mu.RLock()
	defer obj.mu.RUnlock()
	return obj.version
}

func (obj *Jwt) GetSecret() []byte {
	obj.mu.RLock()
	defer obj.mu.RUnlock()
	return obj.secret
}
