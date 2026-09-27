package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
)

type BlockRepository interface {
	Create(ctx context.Context, block domain.Block) (int, error)
	GetById(ctx context.Context, id int) (domain.Block, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Block, error)
	Update(ctx context.Context, block domain.Block) error
	Delete(ctx context.Context, id int) error
}

type BlockPartRepository interface {
	Create(ctx context.Context, blockPart domain.BlockPart) (int, error)
	GetById(ctx context.Context, id int) (domain.BlockPart, error)
	GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error)
	GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error)
}

type BlockPostgres struct {
	db database.DB
}

func NewBlockPostgres(db database.DB) *BlockPostgres {
	return &BlockPostgres{db: db}
}

type BlockPartPostgres struct {
	db database.DB
}

func NewBlockPartPostgres(db database.DB) *BlockPartPostgres {
	return &BlockPartPostgres{db: db}
}
