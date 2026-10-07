package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetAllIds(ctx context.Context) ([]int, error)
	Update(ctx context.Context, user domain.User) error
}

type SettingsRepository interface {
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

func (obj *UserPostgres) GetById(ctx context.Context, id int) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select is_staff, created_at, updated_at, email from "user" where id = $1 and active = true;`
	args := []any{id}
	res := domain.User{Id: id, Active: true}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.IsStaff, &res.CreatedAt, &updatedAt, &res.Email)
	if err != nil {
		log.Error("failed to get user by id", zap.Int("id", id), zap.Error(err))
		return domain.User{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	log.Info("Query executed")
	return res, nil
}

func (obj *UserPostgres) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, is_staff, created_at, updated_at from "user" where email = $1 and active = true;`
	args := []any{email}
	res := domain.User{Email: email, Active: true}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &res.IsStaff, &res.CreatedAt, &updatedAt)
	if err != nil {
		log.Error("failed to get user by email", zap.Error(err))
		return domain.User{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	log.Info("Query executed")
	return res, nil
}

func (obj *UserPostgres) GetAllIds(ctx context.Context) ([]int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id from "user" where active = true;`
	var args []any
	res := make([]int, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get all users ids", zap.Error(err))
		return []int{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Error("failed to get all users ids", zap.Error(err))
			return []int{}, err
		}
		res = append(res, id)
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
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

func (obj *SettingsPostgres) GetById(ctx context.Context, id int) (domain.Settings, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select user_id, "interval", avatar_url, first_name, last_name, updated_at from settings where id = $1;`
	args := []any{id}
	res := domain.Settings{Id: id}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.UserId, &res.Interval, &res.AvatarUrl, &res.FirstName, &res.LastName, &updatedAt)
	if err != nil {
		log.Error("failed to get settings by id", zap.Error(err))
		return domain.Settings{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	log.Info("Query executed")
	return res, nil
}

func (obj *SettingsPostgres) GetByUser(ctx context.Context, user domain.User) (domain.Settings, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, "interval", avatar_url, first_name, last_name, updated_at from settings where user_id = $1;`
	args := []any{user.Id}
	res := domain.Settings{UserId: user.Id}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &res.Interval, &res.AvatarUrl, &res.FirstName, &res.LastName, &updatedAt)
	if err != nil {
		log.Error("failed to get settings by user id", zap.Error(err))
		return domain.Settings{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Time
	}
	log.Info("Query executed")
	return res, nil
}

func (obj *SettingsPostgres) Update(ctx context.Context, user domain.Settings) error {
	//TODO implement me
	panic("implement me")
}

func (obj *SettingsPostgres) UpdateAvatar(ctx context.Context, id int, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}
