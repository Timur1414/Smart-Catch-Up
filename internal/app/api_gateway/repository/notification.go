package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
)

type NotificationRepository interface {
	Create(ctx context.Context, settings domain.Settings) (int, error)
	GetById(ctx context.Context, id int) (domain.Settings, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Settings, error)
	Update(ctx context.Context, settings domain.Settings) error
}

type NotificationPostgres struct {
	db database.DB
}

func NewNotificationPostgres(db database.DB) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj NotificationPostgres) Create(ctx context.Context, settings domain.Settings) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj NotificationPostgres) GetById(ctx context.Context, id int) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj NotificationPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj NotificationPostgres) Update(ctx context.Context, settings domain.Settings) error {
	//TODO implement me
	panic("implement me")
}
