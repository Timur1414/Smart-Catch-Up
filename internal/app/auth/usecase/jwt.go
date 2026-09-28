package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
)

type RefreshTokenUseCase interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error)
	DeleteByUuid(ctx context.Context, uuid string) error
	DeleteByUser(ctx context.Context, user domain.User) error
}

type RefreshToken struct {
	repository repository.RefreshTokenRepository
}

func NewRefreshToken(repo repository.RefreshTokenRepository) *RefreshToken {
	return &RefreshToken{repository: repo}
}

func (obj RefreshToken) Create(ctx context.Context, user domain.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshToken) GetByUuid(ctx context.Context, uuid string) (domain.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshToken) DeleteByUuid(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}

func (obj RefreshToken) DeleteByUser(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
