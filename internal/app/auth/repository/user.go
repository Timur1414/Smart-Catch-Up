package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	SaveTotpSecret(ctx context.Context, userId int, secret string) error
	EnableTotp(ctx context.Context, userId int, backupCodes []string) error
	DisableTotp(ctx context.Context, userId int) error
	RemoveBackupCode(ctx context.Context, userId int, code string) error
}

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (obj *UserPostgres) Create(ctx context.Context, user domain.User) (int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `insert into "user"(email, password) values ($1, $2) returning id;`
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
			return -1, domain.ErrUserAlreadyExists
		case pgerrcode.CheckViolation:
			return -1, domain.ErrDuplicatedData
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create user (not db error)", zap.Error(err))
		return -1, err
	}
	log = logger.ModifyLoggerWithDBQuery(log, query, []any{}, duration)
	log.Info("Query executed")

	log = logger.GetLoggerWithRequestId(ctx)
	query = `insert into settings(user_id, first_name, last_name) values ($1, $2, $3);`
	args = []any{id, "", ""}
	start = time.Now()
	_, err = obj.db.Exec(ctx, query, args...)
	duration = time.Since(start)
	pgErr, ok = errors.AsType[*pgconn.PgError](err)
	if ok {
		log.Error("failed to create settings (db error)", zap.Error(pgErr))
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return -1, domain.ErrSettingsAlreadyExists
		case pgerrcode.CheckViolation:
			return -1, domain.ErrDuplicatedData
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create settings (not db error)", zap.Error(err))
		return -1, err
	}
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")

	log = logger.GetLoggerWithRequestId(ctx)
	query = `insert into digest(user_id) values ($1);`
	args = []any{id, "", ""}
	start = time.Now()
	_, err = obj.db.Exec(ctx, query, args...)
	duration = time.Since(start)
	pgErr, ok = errors.AsType[*pgconn.PgError](err)
	if ok {
		log.Error("failed to create digest (db error)", zap.Error(pgErr))
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return -1, domain.ErrDigestAlreadyExists
		case pgerrcode.CheckViolation:
			return -1, domain.ErrDuplicatedData
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create digest (not db error)", zap.Error(err))
		return -1, err
	}
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return id, nil
}

func (obj *UserPostgres) GetById(ctx context.Context, id int) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select is_staff, created_at, updated_at, email, password, totp_secret, totp_enabled, totp_backup_codes from "user" where id = $1 and active = true;`
	args := []any{id}
	res := domain.User{Id: id, Active: true}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.IsStaff, &res.CreatedAt, &updatedAt, &res.Email, &res.Password, &res.TotpSecret, &res.TotpEnabled, &res.TotpBackupCodes)
	if err != nil {
		log.Error("failed to get user by id", zap.Error(err))
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
	query := `select id, is_staff, created_at, updated_at, password, totp_secret, totp_enabled, totp_backup_codes from "user" where email = $1 and active = true;`
	args := []any{email}
	res := domain.User{Email: email, Active: true}
	var updatedAt pgtype.Timestamp
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Id, &res.IsStaff, &res.CreatedAt, &updatedAt, &res.Password, &res.TotpSecret, &res.TotpEnabled, &res.TotpBackupCodes)
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

func (obj *UserPostgres) Update(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (obj *UserPostgres) SaveTotpSecret(ctx context.Context, userId int, secret string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update "user" set totp_secret = $1, totp_enabled = false where id = $2 and active = true;`
	args := []any{secret, userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to save TOTP secret", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}

func (obj *UserPostgres) EnableTotp(ctx context.Context, userId int, backupCodes []string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update "user" set totp_enabled = true, totp_backup_codes = $1 where id = $2 and active = true;`
	args := []any{backupCodes, userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to enable TOTP", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}

func (obj *UserPostgres) DisableTotp(ctx context.Context, userId int) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update "user" set totp_enabled = false, totp_secret = '', totp_backup_codes = '{}' where id = $1 and active = true;`
	args := []any{userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to disable TOTP", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}

func (obj *UserPostgres) RemoveBackupCode(ctx context.Context, userId int, code string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `update "user" set totp_backup_codes = array_remove(totp_backup_codes, $1) where id = $2 and active = true;`
	args := []any{code, userId}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to remove backup code", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}
