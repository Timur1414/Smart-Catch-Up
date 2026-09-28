package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/redis/go-redis/v9"
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
	//TODO implement me
	panic("implement me")
}

func (obj RefreshTokenRedis) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshTokenRedis) DeleteByUuid(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshTokenRedis) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
