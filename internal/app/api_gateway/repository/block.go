package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func NewBlockPostgres(db *pgxpool.Pool) *BlockPostgres {
	return &BlockPostgres{db: db}
}

func (obj *BlockPostgres) Create(ctx context.Context, block domain.Block) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) GetById(ctx context.Context, id int) (domain.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) Update(ctx context.Context, block domain.Block) error {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

type BlockPartPostgres struct {
	db *pgxpool.Pool
}

func NewBlockPartPostgres(db *pgxpool.Pool) *BlockPartPostgres {
	return &BlockPartPostgres{db: db}
}

func (obj *BlockPartPostgres) Create(ctx context.Context, blockPart domain.BlockPart) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetById(ctx context.Context, id int) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}
