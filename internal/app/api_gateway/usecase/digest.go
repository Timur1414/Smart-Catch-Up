package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type DigestUseCase interface {
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
}

type Digest struct {
	repository repository.DigestRepository
}

func NewDigest(repo repository.DigestRepository) *Digest {
	return &Digest{repository: repo}
}

func (obj *Digest) GetByUser(ctx context.Context, user domain.User) (domain.Digest, error) {
	return obj.repository.GetByUser(ctx, user)
}
