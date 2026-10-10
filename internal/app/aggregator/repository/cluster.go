package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ClusterRepository interface {
	GetAllClusters(ctx context.Context) ([]string, error)
}

type ClusterPostgres struct {
	db *pgxpool.Pool
}

func NewClusterPostgres(db *pgxpool.Pool) *ClusterPostgres {
	return &ClusterPostgres{db: db}
}

func (obj *ClusterPostgres) GetAllClusters(ctx context.Context) ([]string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select enumlabel from pg_enum where enumtypid = 'cluster_type_enum'::regtype order by enumsortorder;`
	var args []any
	res := make([]string, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get clusters", zap.Error(err))
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var cluster string
		err = rows.Scan(&cluster)
		if err != nil {
			log.Error("failed to scan clusters", zap.Error(err))
			return []string{}, err
		}
		res = append(res, cluster)
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}
