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
	id := -1
	err := pgx.BeginFunc(ctx, obj.db, func(tx pgx.Tx) error {
		log := logger.GetLoggerWithRequestId(ctx)
		query := `insert into notification(notification_type, recipient_id, actor_id, actor_name, created_at, payload) values ($1, $2, $3, $4, $5, $6) returning id;`
		args := []any{notification.NotificationType, notification.RecipientId, notification.ActorId, notification.ActorName, notification.CreatedAt, notification.Payload}
		start := time.Now()
		err := obj.db.QueryRow(ctx, query, args...).Scan(&id)
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok {
			log.Error("failed to create notification (db error)", zap.Error(pgErr))
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return errors.New("notification already exists")
			case pgerrcode.CheckViolation:
				return errors.New("DuplicatedData")
			default:
				return pgErr
			}
		}
		if err != nil {
			log.Error("failed to create notification (not db error)", zap.Error(err))
			return err
		}
		duration := time.Since(start)
		log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
		log.Info("Query executed")
		for i := range notification.Actions {
			err = obj.AddNotificationActions(ctx, id, notification.Actions[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	return id, err
}

func (obj *NotificationPostgres) BulkCreate(ctx context.Context, notifications []domain.Notification) error {
	err := pgx.BeginFunc(ctx, obj.db, func(tx pgx.Tx) error {
		log := logger.GetLoggerWithRequestId(ctx)
		types := make([]string, len(notifications))
		recipientIds := make([]int, len(notifications))
		actorIds := make([]int, len(notifications))
		actorNames := make([]string, len(notifications))
		createdAts := make([]time.Time, len(notifications))
		payloads := make([]string, len(notifications))
		for i, notification := range notifications {
			types[i] = notification.NotificationType
			recipientIds[i] = notification.RecipientId
			actorIds[i] = notification.ActorId
			actorNames[i] = notification.ActorName
			createdAts[i] = notification.CreatedAt
			payloads[i] = notification.Payload
		}
		query := `insert into notification(notification_type, recipient_id, actor_id, actor_name, created_at, payload) select * from unnest(
$1::notification_type_enum[], $2::int[], $3::int[], $4::int[], $5::timestamptz[], $6::text[]
) returning id;`
		args := []any{types, recipientIds, actorIds, actorNames, payloads}
		start := time.Now()
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			log.Error("BulkCreate query failed", zap.Error(err))
			return err
		}
		defer rows.Close()
		type actionRow struct {
			notificationId int
			actionType     string
			actionTarget   string
		}
		var allActions []actionRow
		curNotificationIndex := 0
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				log.Error("BulkCreate id scan failed", zap.Error(err))
				return err
			}
			for _, action := range notifications[curNotificationIndex].Actions {
				allActions = append(allActions, actionRow{
					notificationId: id,
					actionType:     action.ActionType,
					actionTarget:   action.ActionTarget,
				})
			}
			curNotificationIndex++
		}
		if len(allActions) == 0 {
			return nil
		}
		actionsCopySource := pgx.CopyFromSlice(len(allActions), func(i int) ([]any, error) {
			return []any{allActions[i].notificationId, allActions[i].actionType, allActions[i].actionTarget}, nil
		})
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"notification_action"},
			[]string{"notification_id", "action_type", "action_target"},
			actionsCopySource,
		)
		if err != nil {
			log.Error("BulkCreate actions failed", zap.Error(err))
			return err
		}
		duration := time.Since(start)
		log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
		log.Info("Query executed")
		return nil
	})
	return err
}

func (obj *NotificationPostgres) AddNotificationActions(ctx context.Context, notificationId int, action domain.NotificationAction) error {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `insert into notification_action(notification_id, action_type, action_target) values ($1, $2, $3);`
	args := []any{notificationId, action.ActionType, action.ActionTarget}
	start := time.Now()
	_, err := obj.db.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to create notification action", zap.Error(err))
		return err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return nil
}
