package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SettingsRepository interface {
	UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error
}

type SettingsPostgres struct {
	db *pgxpool.Pool
}

func NewSettingsPostgres(db *pgxpool.Pool) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (obj *SettingsPostgres) UpdateAvatar(ctx context.Context, userId int, avatarUrl string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update settings set avatar_url = $1 where user_id = $2;`
	args := []any{avatarUrl, userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to update avatar", zap.Int("user_id", userId), zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}
