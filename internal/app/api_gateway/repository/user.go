package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetAllIds(ctx context.Context) ([]int, error)
	GetAllShortUsers(ctx context.Context) ([]domain.User, []domain.Settings, error)
	Update(ctx context.Context, user domain.User) error
}

type SettingsRepository interface {
	GetById(ctx context.Context, id int) (domain.Settings, error)
	GetByUser(ctx context.Context, user domain.User) (domain.Settings, error)
	Update(ctx context.Context, settings domain.Settings) error
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNothingInTable
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNothingInTable
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return []int{}, domain.ErrNothingInTable
		}
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

func (obj *UserPostgres) GetAllShortUsers(ctx context.Context) ([]domain.User, []domain.Settings, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select U.id, S.first_name, S.last_name from "user" U join settings S on U.id = S.user_id where U.active = true;`
	var args []any
	resUsers := make([]domain.User, 0)
	resSettings := make([]domain.Settings, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get all users ids", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.User{}, []domain.Settings{}, domain.ErrNothingInTable
		}
		return []domain.User{}, []domain.Settings{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var firstName, lastName string
		err = rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Error("failed to get all users ids", zap.Error(err))
			return []domain.User{}, []domain.Settings{}, err
		}
		resUsers = append(resUsers, domain.User{Id: id})
		resSettings = append(resSettings, domain.Settings{FirstName: firstName, LastName: lastName})
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return resUsers, resSettings, nil
}

func (obj *UserPostgres) Update(ctx context.Context, user domain.User) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update "user" set email = $1 where id = $2 and active = true;`
	args := []any{user.Email, user.Id}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Warn("failed to update user", zap.Int("userId", user.Id), zap.Error(err))
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.ErrDuplicatedData
			}
		}
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}

type SettingsPostgres struct {
	db *pgxpool.Pool
}

func NewSettingsPostgres(db *pgxpool.Pool) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (obj *SettingsPostgres) GetById(ctx context.Context, id int) (domain.Settings, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select user_id, avatar_url, first_name, last_name, updated_at from settings where id = $1;`
	args := []any{id}
	res := domain.Settings{Id: id}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.UserId, &res.AvatarUrl, &res.FirstName, &res.LastName, &updatedAt)
	if err != nil {
		log.Error("failed to get settings by id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Settings{}, domain.ErrNothingInTable
		}
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
	query := `select id,avatar_url, first_name, last_name, updated_at from settings where user_id = $1;`
	args := []any{user.Id}
	res := domain.Settings{UserId: user.Id}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &res.AvatarUrl, &res.FirstName, &res.LastName, &updatedAt)
	if err != nil {
		log.Error("failed to get settings by user id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Settings{}, domain.ErrNothingInTable
		}
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

func (obj *SettingsPostgres) Update(ctx context.Context, settings domain.Settings) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update settings set first_name = $1, last_name = $2 where id = $3;`
	args := []any{settings.FirstName, settings.LastName, settings.Id}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Warn("failed to update settings", zap.Int("id", settings.Id), zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}

func (obj *SettingsPostgres) UpdateAvatar(ctx context.Context, id int, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}
