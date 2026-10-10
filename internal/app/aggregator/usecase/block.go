package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
)

type BlockUseCase interface {
	Create(ctx context.Context)
}

type Block struct {
	repository repository.BlockRepository
}

func NewBlock(repository repository.BlockRepository) *Block {
	return &Block{repository: repository}
}
