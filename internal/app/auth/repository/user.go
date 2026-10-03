package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	log := logger.GetLoggerWithRequestId(ctx)
	query := `insert into "user"(email, password) values ('$1', '$2') returning id;`
	args := []any{user.Email, user.Password}
	var id int
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&id)
	duration := time.Since(start)
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if ok {
		log.Error("failed to create user (db error)", zap.Error(pgErr))
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return -1, errors.New("user already exists")
		case pgerrcode.CheckViolation:
			return -1, errors.New("DuplicatedData")
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create user (not db error)", zap.Error(err))
		return -1, err
	}
	log = logger.ModifyLoggerWithDBQuery(log, query, []any{}, duration)
	log.Info("created user")

	query = `insert into settings(user_id, first_name, last_name) values ();`
	args = []any{id, "", ""}
	start = time.Now()
	_, err = obj.db.Exec(ctx, query, args...)
	duration = time.Since(start)
	pgErr, ok = errors.AsType[*pgconn.PgError](err)
	if ok {
		log.Error("failed to create settings (db error)", zap.Error(pgErr))
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return -1, errors.New("settings already exists")
		case pgerrcode.CheckViolation:
			return -1, errors.New("DuplicatedData")
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create settings (not db error)", zap.Error(err))
		return -1, err
	}
	return id, nil
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
