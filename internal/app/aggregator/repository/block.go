package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockRepository interface {
	Create(ctx context.Context)
}

type BlockPostgres struct {
	db *pgxpool.Pool
}

func NewBlockPostgres(db *pgxpool.Pool) *BlockPostgres {
	return &BlockPostgres{db: db}
}
