package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) (int, error)
	BulkCreate(ctx context.Context, notifications []domain.Notification) ([]int, error)
}

type NotificationPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationPostgres(db *pgxpool.Pool) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj *NotificationPostgres) Create(ctx context.Context, notification domain.Notification) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *NotificationPostgres) BulkCreate(ctx context.Context, notifications []domain.Notification) ([]int, error) {
	//TODO implement me
	panic("implement me")
}
