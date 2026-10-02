package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type ClusterUseCase interface {
	GetById(ctx context.Context, id int) (domain.Cluster, error)
	GetByType(ctx context.Context, clusterType string) (domain.Cluster, error)
	GetAll(ctx context.Context) ([]domain.Cluster, error)
}

type Cluster struct {
	repository repository.ClusterRepository
}

func NewCluster(repo repository.ClusterRepository) *Cluster {
	return &Cluster{repository: repo}
}

func (obj *Cluster) GetById(ctx context.Context, id int) (domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *Cluster) GetByType(ctx context.Context, clusterType string) (domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *Cluster) GetAll(ctx context.Context) ([]domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}
