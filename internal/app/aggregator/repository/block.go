package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BlockRepository interface {
	Create(ctx context.Context, block domain.Block) (int, error)
}

type BlockPostgres struct {
	db *pgxpool.Pool
}

func NewBlockPostgres(db *pgxpool.Pool) *BlockPostgres {
	return &BlockPostgres{db: db}
}

func (obj *BlockPostgres) Create(ctx context.Context, block domain.Block) (int, error) {
	blockId := 0
	err := pgx.BeginFunc(ctx, obj.db, func(tx pgx.Tx) error {
		log := logger.GetLoggerWithRequestId(ctx)
		query := `insert into block(user_id, start_at, end_at, status) values ($1, $2, $3, $4) returning id;`
		args := []any{block.UserId, block.Start, block.End, block.Status}
		start := time.Now()
		err := tx.QueryRow(ctx, query, args...).Scan(&blockId)
		if err != nil {
			log.Error("failed to create block", zap.Error(err))
			return err
		}
		duration := time.Since(start)
		log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
		log.Info("Query executed")

		query = `insert into block_part(block_id, cluster, text, start_at, end_at) values ($1, $2, $3, $4, $5);`
		for _, part := range block.Parts {
			args = []any{blockId, part.ClusterType, part.Text, part.Start, part.End}
			start = time.Now()
			_, err = tx.Exec(ctx, query, args...)
			if err != nil {
				log.Error("failed to create block", zap.Error(err))
				return err
			}
			duration = time.Since(start)
			log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
			log.Info("Query executed")
		}
		return nil
	})
	return blockId, err
}
