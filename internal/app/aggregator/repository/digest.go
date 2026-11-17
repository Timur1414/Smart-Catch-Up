package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DigestRepository interface {
	UpdateBlock(ctx context.Context, userId int, newBlockId int) error
}

type CacheRepository interface {
	Add(ctx context.Context, userId int, block domain.Block) error
}

const (
	cacheTTL    = 24 * time.Hour
	cachePrefix = "cache:digest:"
)

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

type CacheRedis struct {
	db *redis.Client
}

func NewCacheRedis(client *redis.Client) *CacheRedis {
	return &CacheRedis{db: client}
}

func (obj *CacheRedis) Add(ctx context.Context, userId int, block domain.Block) error {
	log := logger.GetLogger()
	cachedDigest := domain.NewDigestCache(block)
	data, err := json.Marshal(cachedDigest)
	if err != nil {
		log.Error("failed to marshal cached_digest", zap.Error(err))
		return err
	}
	key := cachePrefix + strconv.Itoa(userId)
	start := time.Now()
	err = obj.db.Set(ctx, key, data, cacheTTL).Err()
	if err != nil {
		log.Error("failed to set cache", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log.Info("Query executed", zap.Duration("duration", duration))
	return nil
}
