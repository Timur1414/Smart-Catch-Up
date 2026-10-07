package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	IsExists(ctx context.Context, user domain.User) (domain.User, error)
}

type User struct {
	repository repository.UserRepository
}

func NewUser(repo repository.UserRepository) *User {
	return &User{repository: repo}
}

func (obj *User) Create(ctx context.Context, user domain.User) (int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userToCreate := domain.User{
		Email: user.Email,
	}
	origPassword := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(origPassword, bcrypt.DefaultCost)
	if err != nil {
		log.Warn("failed to hash password", zap.Error(err))
		return -1, domain.ErrFailedToHashPassword
	}
	userToCreate.Password = string(hashedPassword)
	return obj.repository.Create(ctx, userToCreate)
}

func (obj *User) GetById(ctx context.Context, id int) (domain.User, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj *User) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return obj.repository.GetByEmail(ctx, email)
}

func (obj *User) Update(ctx context.Context, user domain.User) error {
	return obj.repository.Update(ctx, user)
}

func (obj *User) IsExists(ctx context.Context, user domain.User) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userFromDb, err := obj.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password))
	if err != nil {
		log.Warn("failed to compare password", zap.Error(err))
		return domain.User{}, err
	}
	return userFromDb, nil
}
