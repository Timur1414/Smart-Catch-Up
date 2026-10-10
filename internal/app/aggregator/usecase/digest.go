package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
)

type DigestUseCase interface {
	Update(ctx context.Context, userId int, newBlockId int) error
}

type Digest struct {
	repository repository.DigestPostgres
}

func NewDigest(repository repository.DigestPostgres) *Digest {
	return &Digest{repository: repository}
}

func (obj *Digest) Update(ctx context.Context, userId int, newBlockId int) error {
	return obj.repository.UpdateBlock(ctx, userId, newBlockId)
}
