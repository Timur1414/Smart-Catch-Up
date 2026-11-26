package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/config"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DigestRepository interface {
	GetByUser(ctx context.Context, user domain.User) (domain.Digest, error)
}

type CacheRepository interface {
	Get(ctx context.Context, userId int) (domain.DigestCache, error)
	Update(ctx context.Context, userId int, block domain.Block) error
}

type DigestPostgres struct {
	db *pgxpool.Pool
}

func NewDigestPostgres(db *pgxpool.Pool) *DigestPostgres {
	return &DigestPostgres{db: db}
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

type CacheRedis struct {
	db  *redis.Client
	cfg *config.Config
}

func NewCacheRedis(client *redis.Client, cfg *config.Config) *CacheRedis {
	return &CacheRedis{
		db:  client,
		cfg: cfg,
	}
}

func (obj *CacheRedis) Get(ctx context.Context, userId int) (domain.DigestCache, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	key := obj.cfg.Common.DigestCachePrefix + strconv.Itoa(userId)
	start := time.Now()
	data, err := obj.db.Get(ctx, key).Bytes()
	if err != nil {
		log.Error("failed to get cache", zap.Int("userId", userId), zap.Error(err))
		return domain.DigestCache{}, err
	}
	duration := time.Since(start)
	var digest domain.DigestCache
	err = json.Unmarshal(data, &digest)
	if err != nil {
		log.Error("failed to unmarshal cache", zap.Int("userId", userId), zap.Error(err))
		return domain.DigestCache{}, err
	}
	log.Info("Query executed", zap.Duration("duration", duration))
	return digest, nil
}

func (obj *CacheRedis) Update(ctx context.Context, userId int, block domain.Block) error {
	log := logger.GetLoggerWithRequestId(ctx)
	cachedDigest := domain.NewDigestCache(block)
	data, err := json.Marshal(cachedDigest)
	if err != nil {
		log.Error("failed to marshal cached_digest", zap.Error(err))
		return err
	}
	key := obj.cfg.Common.DigestCachePrefix + strconv.Itoa(userId)
	start := time.Now()
	err = obj.db.Set(ctx, key, data, obj.cfg.Common.DigestCacheTtl).Err()
	if err != nil {
		log.Error("failed to set cache", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log.Info("Query executed", zap.Duration("duration", duration))
	return nil
}
