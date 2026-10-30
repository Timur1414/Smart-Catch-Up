package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type BlockUseCase interface {
	GetById(ctx context.Context, id int) (domain.Block, error)
}

type Block struct {
	repository repository.BlockRepository
}

func NewBlock(repo repository.BlockRepository) *Block {
	return &Block{repository: repo}
}

func (obj Block) GetById(ctx context.Context, id int) (domain.Block, error) {
	return obj.repository.GetById(ctx, id)
}
