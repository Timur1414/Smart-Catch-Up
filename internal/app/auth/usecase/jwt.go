package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
)

type JwtUseCase interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error)
	DeleteByUuid(ctx context.Context, uuid string) error
	DeleteByUser(ctx context.Context, user domain.User) error
	CheckAccessToken(tokenStr string) (bool, int)
	CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int)
}

type Jwt struct {
	repository repository.RefreshTokenRepository
}

func NewJwt(repo repository.RefreshTokenRepository) *Jwt {
	return &Jwt{repository: repo}
}

func (obj Jwt) Create(ctx context.Context, user domain.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Jwt) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Jwt) DeleteByUuid(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}

func (obj Jwt) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (obj Jwt) CheckAccessToken(tokenStr string) (bool, int) {
	//TODO implement me
	panic("implement me")
}

func (obj Jwt) CheckRefreshToken(ctx context.Context, tokenStr string) (bool, int) {
	//TODO implement me
	panic("implement me")
}
