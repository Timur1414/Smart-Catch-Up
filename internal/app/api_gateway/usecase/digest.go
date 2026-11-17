package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type DigestUseCase interface {
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
	GetFromCache(ctx context.Context, user domain.User) (domain.DigestCache, error)
	UpdateCache(ctx context.Context, user domain.User, digest domain.Block) error
}

type Digest struct {
	repository repository.DigestRepository
	cache      repository.CacheRepository
}

func NewDigest(repo repository.DigestRepository, cache repository.CacheRepository) *Digest {
	return &Digest{
		repository: repo,
		cache:      cache,
	}
}

func (obj *Digest) GetByUser(ctx context.Context, user domain.User) (domain.Digest, error) {
	return obj.repository.GetByUser(ctx, user)
}

func (obj *Digest) GetFromCache(ctx context.Context, user domain.User) (domain.DigestCache, error) {
	return obj.cache.Get(ctx, user.Id)
}

func (obj *Digest) UpdateCache(ctx context.Context, user domain.User, block domain.Block) error {
	return obj.cache.Update(ctx, user.Id, block)
}
