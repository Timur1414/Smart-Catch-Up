package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/config"
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
	ExtractClaims(ctx context.Context, tokenStr string) (domain.Claims, error)
	CreateTempToken(ctx context.Context, userId int) (string, error)
	CheckTempToken(ctx context.Context, tokenStr string) (bool, int)
	DeleteTempToken(ctx context.Context, tokenStr string) error
	IncrementTempTokenAttempts(ctx context.Context, userId int) (int, error)
	Mark2FACodeUsed(ctx context.Context, userId int, code string) (bool, error)
}

type Jwt struct {
	refreshTokenRepository repository.RefreshTokenRepository
	tempTokenRepository    repository.TwoFactorTokenRepository
	secret                 []byte
	version                string
	cfg                    config.Config
}

func NewJwt(repo repository.RefreshTokenRepository, twoFactorRepo repository.TwoFactorTokenRepository, jwtConfig config.JwtConfig, cfg config.Config) (*Jwt, error) {
	if len(jwtConfig.Secret) < 8 {
		return nil, domain.ErrSecretTooShort
	}
	if jwtConfig.Version == "" {
		return nil, domain.ErrVersionIsEmpty
	}
	return &Jwt{
		refreshTokenRepository: repo,
		tempTokenRepository:    twoFactorRepo,
		secret:                 []byte(jwtConfig.Secret),
		version:                jwtConfig.Version,
		cfg:                    cfg,
	}, nil
}

func (obj *Jwt) CreatePair(ctx context.Context, userId int) (string, string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userIdStr := strconv.Itoa(userId)

	err := obj.refreshTokenRepository.DeleteByUser(ctx, userId)
	if err != nil {
		return "", "", err
	}

	claims := domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(obj.cfg.Common.AccessTokenExpirationTime).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    obj.version,
			Subject:   userIdStr,
		},
		Type: obj.cfg.Auth.AccessTokenType,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(obj.secret)
	if err != nil {
		log.Error("failed to sign access token",
			zap.Int("user_id", userId),
			zap.Error(err))
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(obj.cfg.Common.RefreshTokenExpirationTime).UTC()
	refreshId := uuid.New().String()
	refreshClaims := domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshId,
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    obj.version,
			Subject:   userIdStr,
		},
		Type: obj.cfg.Auth.RefreshTokenType,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(obj.secret)
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
	_, err = obj.refreshTokenRepository.Create(ctx, refreshTokenToSave)
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
	if claims.Type != obj.cfg.Auth.RefreshTokenType {
		log.Warn("invalid refresh token type")
		return domain.ErrInvalidTokenType
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Warn("failed to convert subject to int", zap.Error(err))
		return err
	}
	return obj.refreshTokenRepository.DeleteByUser(ctx, userId)
}

func (obj *Jwt) CheckAccessToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, tokenStr)
	if err != nil {
		return false, -1
	}
	if claims.Type != obj.cfg.Auth.AccessTokenType {
		log.Warn("invalid access token type")
		return false, -1
	}
	if claims.Issuer != obj.version {
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
	if claims.Type != obj.cfg.Auth.RefreshTokenType {
		log.Warn("invalid refresh token type")
		return false, -1
	}
	if claims.Issuer != obj.version {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return false, -1
	}

	storedToken, err := obj.refreshTokenRepository.GetByUser(ctx, userId)
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

func (obj *Jwt) ExtractClaims(ctx context.Context, tokenStr string) (domain.Claims, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	var claims domain.Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return obj.secret, nil
	})
	if err != nil || !token.Valid {
		log.Warn("Token parse error", zap.Error(err))
		return domain.Claims{}, err
	}
	return claims, nil
}

func (obj *Jwt) CreateTempToken(ctx context.Context, userId int) (string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userIdStr := strconv.Itoa(userId)
	err := obj.tempTokenRepository.Delete(ctx, userId)
	if err != nil {
		return "", err
	}
	tempTokenExpirationDate := time.Now().Add(obj.cfg.Auth.TwoFactorTokenExpirationTime).UTC()
	tempUuid := uuid.New().String()
	claims := domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tempUuid,
			ExpiresAt: jwt.NewNumericDate(tempTokenExpirationDate),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    obj.version,
			Subject:   userIdStr,
		},
		Type: obj.cfg.Auth.TempTokenType,
	}
	tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tempTokenStr, err := tempToken.SignedString(obj.secret)
	if err != nil {
		log.Error("failed to sign temp token", zap.Error(err))
		return "", err
	}

	tempTokenToSave := domain.RefreshToken{
		Uuid:      tempUuid,
		UserId:    userId,
		ExpiredAt: tempTokenExpirationDate,
	}
	err = obj.tempTokenRepository.Create(ctx, tempTokenToSave)
	if err != nil {
		log.Error("failed to save temp token", zap.Error(err))
		return "", err
	}

	return tempTokenStr, nil
}

func (obj *Jwt) CheckTempToken(ctx context.Context, tokenStr string) (bool, int) {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, tokenStr)
	if err != nil {
		return false, -1
	}
	if claims.Type != obj.cfg.Auth.TempTokenType {
		log.Warn("invalid temp token type")
		return false, -1
	}
	if claims.Issuer != obj.version {
		log.Warn("invalid token (invalid version)")
		return false, -1
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Warn("failed to convert subject to int", zap.Error(err))
		return false, -1
	}
	storedToken, err := obj.tempTokenRepository.GetByUser(ctx, userId)
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

func (obj *Jwt) DeleteTempToken(ctx context.Context, tokenStr string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	claims, err := obj.ExtractClaims(ctx, tokenStr)
	if err != nil {
		return err
	}
	if claims.Type != obj.cfg.Auth.TempTokenType {
		log.Warn("invalid temp token type")
		return domain.ErrInvalidTokenType
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Warn("failed to convert subject to int", zap.Error(err))
		return err
	}
	return obj.tempTokenRepository.Delete(ctx, userId)
}

func (obj *Jwt) IncrementTempTokenAttempts(ctx context.Context, userId int) (int, error) {
	return obj.tempTokenRepository.IncrementAttempts(ctx, userId)
}

func (obj *Jwt) Mark2FACodeUsed(ctx context.Context, userId int, code string) (bool, error) {
	return obj.tempTokenRepository.Mark2FACodeUsed(ctx, userId, code)
}
