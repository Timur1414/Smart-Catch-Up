package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DigestRepository interface {
	Create(ctx context.Context, digest domain.Digest) (int, error)
	GetById(ctx context.Context, id int) (domain.Digest, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
	Update(ctx context.Context, digest domain.Digest) error
}

type DigestPostgres struct {
	db *pgxpool.Pool
}

func NewDigestPostgres(db *pgxpool.Pool) *DigestPostgres {
	return &DigestPostgres{db: db}
}

func (obj DigestPostgres) Create(ctx context.Context, digest domain.Digest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj DigestPostgres) GetById(ctx context.Context, id int) (domain.Digest, error) {
	//TODO implement me
	panic("implement me")
}

func (obj DigestPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Digest, error) {
	//TODO implement me
	panic("implement me")
}

func (obj DigestPostgres) Update(ctx context.Context, digest domain.Digest) error {
	//TODO implement me
	panic("implement me")
}
