package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
)

type ClusterUseCase interface {
	GetAllClusters(ctx context.Context) ([]string, error)
}

type Cluster struct {
	repository repository.ClusterRepository
}
