package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetById(ctx context.Context, id int) (domain.User, error)
	GetActors(ctx context.Context, curUserId int) ([]domain.NotificationActor, error)
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

func (obj *UserPostgres) GetActors(ctx context.Context, curUserId int) ([]domain.NotificationActor, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select "user".id, first_name, last_name from "user" join settings on "user".id = settings.user_id where "user".id <> $1;`
	args := []any{curUserId}
	res := make([]domain.NotificationActor, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get users", zap.Error(err))
		return []domain.NotificationActor{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var firstName, lastName string
		if err = rows.Scan(&id, &firstName, &lastName); err != nil {
			log.Error("failed to scan user", zap.Error(err))
			return []domain.NotificationActor{}, err
		}
		res = append(res, domain.NotificationActor{Id: id, Name: firstName + " " + lastName})
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}
