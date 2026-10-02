package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (obj *UserPostgres) Create(ctx context.Context, user domain.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *UserPostgres) GetById(ctx context.Context, id int) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *UserPostgres) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *UserPostgres) Update(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

type SettingsPostgres struct {
	db *pgxpool.Pool
}

func NewSettingsPostgres(db *pgxpool.Pool) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (obj *SettingsPostgres) Create(ctx context.Context, settings domain.Settings) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *SettingsPostgres) GetById(ctx context.Context, id int) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *SettingsPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Settings, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *SettingsPostgres) Update(ctx context.Context, user domain.Settings) error {
	//TODO implement me
	panic("implement me")
}

func (obj *SettingsPostgres) UpdateAvatar(ctx context.Context, id int, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}
