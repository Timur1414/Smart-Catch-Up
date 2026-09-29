package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
)

type UserUseCase interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type User struct {
	repository repository.UserRepository
}

func NewUser(repo repository.UserRepository) *User {
	return &User{repository: repo}
}

func (obj *User) Create(ctx context.Context, user domain.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *User) GetById(ctx context.Context, id int) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *User) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *User) Update(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
