package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, postgresCfg config.PostgresConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(postgresCfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}
	cfg.MaxConns = 90
	cfg.MinConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 10 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	return pool, nil
}
