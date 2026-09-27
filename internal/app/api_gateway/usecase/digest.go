package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type DigestUseCase interface {
	Create(ctx context.Context, digest domain.Digest) (int, error)
	GetById(ctx context.Context, id int) (domain.Digest, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
	Update(ctx context.Context, digest domain.Digest) error
}

type Digest struct {
	repository repository.DigestRepository
}

func NewDigest(repo repository.DigestRepository) *Digest {
	return &Digest{repository: repo}
}

func (obj Digest) Create(ctx context.Context, digest domain.Digest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Digest) GetById(ctx context.Context, id int) (domain.Digest, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Digest) GetByUser(ctx context.Context, user domain.User) (domain.Digest, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Digest) Update(ctx context.Context, digest domain.Digest) error {
	//TODO implement me
	panic("implement me")
}
