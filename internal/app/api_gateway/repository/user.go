package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type SettingsRepository interface {
	Create(ctx context.Context, settings domain.Settings) (int, error)
	GetById(ctx context.Context, id int) (domain.Settings, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Settings, error)
	Update(ctx context.Context, user domain.Settings) error
	UpdateAvatar(ctx context.Context, id int, avatarUrl string) error
}

type UserPostgres struct {
	db database.DB
}

func NewUserPostgres(db database.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (obj UserPostgres) Create(ctx context.Context, user domain.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj UserPostgres) GetById(ctx context.Context, id int) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj UserPostgres) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj UserPostgres) Update(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

type SettingsPostgres struct {
	db database.DB
}

func NewSettingsPostgres(db database.DB) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (obj SettingsPostgres) Create(ctx context.Context, settings domain.Settings) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj SettingsPostgres) GetById(ctx context.Context, id int) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj SettingsPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj SettingsPostgres) Update(ctx context.Context, user domain.Settings) error {
	//TODO implement me
	panic("implement me")
}

func (obj SettingsPostgres) UpdateAvatar(ctx context.Context, id int, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}
