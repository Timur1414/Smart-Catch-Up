package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
)

type UserUseCase interface {
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type SettingsUseCase interface {
	GetById(ctx context.Context, id int) (domain.Settings, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Settings, error)
	Update(ctx context.Context, settings domain.Settings) error
	UpdateAvatar(ctx context.Context, avatarUrl string) error
}

type User struct {
	repository repository.UserRepository
}

func NewUser(repo repository.UserRepository) *User {
	return &User{repository: repo}
}

func (obj *User) GetById(ctx context.Context, id int) (domain.User, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj *User) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (obj *User) Update(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

type Settings struct {
	repository repository.SettingsRepository
}

func NewSettings(repo repository.SettingsRepository) *Settings {
	return &Settings{repository: repo}
}

func (obj *Settings) GetById(ctx context.Context, id int) (domain.Settings, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj *Settings) GetByUser(ctx context.Context, user domain.User) (domain.Settings, error) {
	return obj.repository.GetByUser(ctx, user)
}

func (obj *Settings) Update(ctx context.Context, settings domain.Settings) error {
	//TODO implement me
	panic("implement me")
}

func (obj *Settings) UpdateAvatar(ctx context.Context, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}
