package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BlockRepository interface {
	Create(ctx context.Context, block domain.Block) (int, error)
	GetById(ctx context.Context, id int) (domain.Block, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Block, error)
	Update(ctx context.Context, block domain.Block) error
	Delete(ctx context.Context, id int) error
}

type BlockPartRepository interface {
	Create(ctx context.Context, blockPart domain.BlockPart) (int, error)
	GetById(ctx context.Context, id int) (domain.BlockPart, error)
	GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error)
	GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error)
}

type BlockPostgres struct {
	db *pgxpool.Pool
}

func NewBlockPostgres(db *pgxpool.Pool) *BlockPostgres {
	return &BlockPostgres{db: db}
}

func (obj *BlockPostgres) Create(ctx context.Context, block domain.Block) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) GetById(ctx context.Context, id int) (domain.Block, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select user_id, start_at, end_at, status from block where id = $1;`
	args := []any{id}
	res := domain.Block{Id: id}
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.UserId, &res.Start, &res.End, &res.Status)
	if err != nil {
		log.Error("failed to get block", zap.Int("id", id), zap.Error(err))
		return domain.Block{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")

	log = logger.GetLoggerWithRequestId(ctx)
	query = `select id, cluster, text, start_at, end_at from block_part where block_id = $1;`
	args = []any{id}
	parts := make([]domain.BlockPart, 0)
	start = time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get block parts", zap.Int("blockId", id), zap.Error(err))
		return domain.Block{}, err
	}
	defer rows.Close()
	for rows.Next() {
		part := domain.BlockPart{BlockId: id}
		err = rows.Scan(&part.Id, &part.ClusterType, &part.Text, &part.Start, &part.End)
		if err != nil {
			log.Error("failed to scan block part", zap.Int("blockId", id), zap.Error(err))
			return domain.Block{}, err
		}
		parts = append(parts, part)
	}
	duration = time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	res.Parts = parts
	return res, nil
}

func (obj *BlockPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) Update(ctx context.Context, block domain.Block) error {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPostgres) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

type BlockPartPostgres struct {
	db *pgxpool.Pool
}

func NewBlockPartPostgres(db *pgxpool.Pool) *BlockPartPostgres {
	return &BlockPartPostgres{db: db}
}

func (obj *BlockPartPostgres) Create(ctx context.Context, blockPart domain.BlockPart) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetById(ctx context.Context, id int) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetByBlock(ctx context.Context, block domain.Block) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *BlockPartPostgres) GetByType(ctx context.Context, clusterType string) (domain.BlockPart, error) {
	//TODO implement me
	panic("implement me")
}
