package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	GetById(ctx context.Context, id int) (domain.Notification, error)
	GetByUser(ctx context.Context, userId int) ([]domain.Notification, error)
}

type NotificationPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationPostgres(db *pgxpool.Pool) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj *NotificationPostgres) GetById(ctx context.Context, id int) (domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *NotificationPostgres) GetByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}
