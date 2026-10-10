package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetById(ctx context.Context, id int) (domain.User, error)
}

type SettingsRepository interface {
	GetByUser(ctx context.Context, id int) (domain.Settings, error)
}

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (obj *UserPostgres) GetById(ctx context.Context, id int) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select email from "user" where id = $1 and active = true;`
	args := []any{id}
	res := domain.User{Id: id, Active: true}
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Email)
	if err != nil {
		log.Error("failed to get user by id", zap.Int("id", id), zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNothingInTable
		}
		return domain.User{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

type SettingsPostgres struct {
	db *pgxpool.Pool
}

func NewSettingsPostgres(db *pgxpool.Pool) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (obj *SettingsPostgres) GetByUser(ctx context.Context, id int) (domain.Settings, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, avatar_url, first_name, last_name from settings where user_id = $1;`
	args := []any{id}
	res := domain.Settings{UserId: id}
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &res.AvatarUrl, &res.FirstName, &res.LastName)
	if err != nil {
		log.Error("failed to get user by id", zap.Int("id", id), zap.Error(err))
		return domain.Settings{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}
