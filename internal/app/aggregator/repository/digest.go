package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DigestRepository interface {
	UpdateBlock(ctx context.Context, userId int, newBlockId int) error
}

type DigestPostgres struct {
	db *pgxpool.Pool
}

func NewDigestPostgres(db *pgxpool.Pool) *DigestPostgres {
	return &DigestPostgres{db: db}
}

func (obj *DigestPostgres) UpdateBlock(ctx context.Context, userId int, newBlockId int) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update digest set block_id = $1 where user_id = $2`
	args := []any{newBlockId, userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to update digest", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}
