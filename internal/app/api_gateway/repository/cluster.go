package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClusterRepository interface {
	GetById(ctx context.Context, id int) (domain.Cluster, error)
	GetByType(ctx context.Context, clusterType string) (domain.Cluster, error)
	GetAll(ctx context.Context) ([]domain.Cluster, error)
}

type ClusterPostgres struct {
	db *pgxpool.Pool
}

func NewClusterPostgres(db *pgxpool.Pool) *ClusterPostgres {
	return &ClusterPostgres{db: db}
}

func (obj *ClusterPostgres) GetById(ctx context.Context, id int) (domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *ClusterPostgres) GetByType(ctx context.Context, clusterType string) (domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *ClusterPostgres) GetAll(ctx context.Context) ([]domain.Cluster, error) {
	//TODO implement me
	panic("implement me")
}
