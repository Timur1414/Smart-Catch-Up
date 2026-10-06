package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) (int, error)
	BulkCreate(ctx context.Context, notifications []domain.Notification) error
}

type NotificationPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationPostgres(db *pgxpool.Pool) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj *NotificationPostgres) Create(ctx context.Context, notification domain.Notification) (int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `insert into notification(notification_type, recipient_id, actor_id, actor_name, payload) values ($1, $2, $3, $4, $5) returning id;`
	args := []any{notification.NotificationType, notification.RecipientId, notification.ActorId, notification.ActorName, notification.Payload}
	var id int
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&id)
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if ok {
		log.Error("failed to create notification (db error)", zap.Error(pgErr))
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return -1, errors.New("notification already exists")
		case pgerrcode.CheckViolation:
			return -1, errors.New("DuplicatedData")
		default:
			return -1, pgErr
		}
	}
	if err != nil {
		log.Error("failed to create notification (not db error)", zap.Error(err))
		return -1, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return id, nil
}

func (obj *NotificationPostgres) BulkCreate(ctx context.Context, notifications []domain.Notification) error {
	log := logger.GetLoggerWithRequestId(ctx)
	rows := pgx.CopyFromSlice(len(notifications), func(i int) ([]any, error) {
		return []any{notifications[i].NotificationType, notifications[i].RecipientId, notifications[i].ActorId, notifications[i].ActorName, notifications[i].Payload}, nil
	})
	_, err := obj.db.CopyFrom(
		ctx,
		pgx.Identifier{"notification"},
		[]string{"notification_type", "recipient_id", "actor_id", "actor_name", "payload"},
		rows,
	)
	if err != nil {
		return err
	}
	log.Info("Query executed")
	return nil
}
