package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type BlockUseCase interface {
	Create(ctx context.Context, block domain.Block) (int, error)
	GetById(ctx context.Context, id int) (domain.Block, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Block, error)
	Update(ctx context.Context, block domain.Block) error
	Delete(ctx context.Context, id int) error
}

type BlockPartUseCase interface {
	Create(ctx context.Context, blockPart domain.BlockPart) (int, error)
	GetById(ctx context.Context, id int) (domain.BlockPart, error)
	GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error)
	GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error)
}

type Block struct {
	repository repository.BlockRepository
}

func NewBlock(repo repository.BlockRepository) *Block {
	return &Block{repository: repo}
}

func (obj Block) Create(ctx context.Context, block domain.Block) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj Block) GetById(ctx context.Context, id int) (domain.Block, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj Block) GetByUser(ctx context.Context, user domain.User) (domain.Block, error) {
	return obj.repository.GetByUser(ctx, user)
}

func (obj Block) Update(ctx context.Context, block domain.Block) error {
	//TODO implement me
	panic("implement me")
}

func (obj Block) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

type BlockPart struct {
	repository repository.BlockPartRepository
}

func NewBlockPart(repo repository.BlockPartRepository) *BlockPart {
	return &BlockPart{repository: repo}
}

func (obj *BlockPart) Create(ctx context.Context, blockPart domain.BlockPart) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPart) GetById(ctx context.Context, id int) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPart) GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPart) GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}
