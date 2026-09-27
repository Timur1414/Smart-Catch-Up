package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
)

type DigestRepository interface {
	Create(ctx context.Context, digest domain.Digest) (int, error)
	GetById(ctx context.Context, id int) (domain.Digest, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
	Update(ctx context.Context, digest domain.Digest) error
}

type DigestPostgres struct {
	db database.DB
}

func NewDigestPostgres(db database.DB) *DigestPostgres {
	return &DigestPostgres{db: db}
}
