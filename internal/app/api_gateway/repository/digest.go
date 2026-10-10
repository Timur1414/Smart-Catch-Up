package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DigestRepository interface {
	Create(ctx context.Context, digest domain.Digest) (int, error)
	GetById(ctx context.Context, id int) (domain.Digest, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
	Update(ctx context.Context, digest domain.Digest) error
}

type DigestPostgres struct {
	db *pgxpool.Pool
}

func NewDigestPostgres(db *pgxpool.Pool) *DigestPostgres {
	return &DigestPostgres{db: db}
}

func (obj *DigestPostgres) Create(ctx context.Context, digest domain.Digest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *DigestPostgres) GetById(ctx context.Context, id int) (domain.Digest, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select user_id, block_id, created_at, updated_at from digest where id = $1;`
	args := []any{id}
	res := domain.Digest{Id: id}
	var blockId pgtype.Int8
	var updatedAt pgtype.Timestamptz
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.UserId, &blockId, &res.CreatedAt, &updatedAt)
	if err != nil {
		log.Error("failed to get digest by id", zap.Int("id", id), zap.Error(err))
		return domain.Digest{}, err
	}
	if blockId.Valid {
		res.BlockId = int(blockId.Int64)
	}
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

func (obj *DigestPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Digest, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, block_id, created_at, updated_at from digest where user_id = $1;`
	args := []any{user.Id}
	res := domain.Digest{UserId: user.Id}
	var blockId pgtype.Int8
	var updatedAt pgtype.Timestamptz
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &blockId, &res.CreatedAt, &updatedAt)
	if err != nil {
		log.Error("failed to get digest by user", zap.Int("userId", user.Id), zap.Error(err))
		return domain.Digest{}, err
	}
	if blockId.Valid {
		res.BlockId = int(blockId.Int64)
	}
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

func (obj *DigestPostgres) Update(ctx context.Context, digest domain.Digest) error {
	//TODO implement me
	panic("implement me")
}
