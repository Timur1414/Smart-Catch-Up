package repository

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
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
