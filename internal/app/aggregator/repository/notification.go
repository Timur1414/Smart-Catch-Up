package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	GetAll(ctx context.Context) ([]domain.Notification, error)
	GetById(ctx context.Context, id int) (domain.Notification, error)
	GetByUser(ctx context.Context, userId int) ([]domain.Notification, error)
}

type NotificationPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationPostgres(db *pgxpool.Pool) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj *NotificationPostgres) GetAll(ctx context.Context) ([]domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, notification_type, recipient_id, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where cluster is null order by created_at;`
	var args []any
	res := make([]domain.Notification, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("fail to get all notifications", zap.Error(err))
		return []domain.Notification{}, err
	}
	defer rows.Close()
	for rows.Next() {
		notification := domain.Notification{Cluster: ""}
		var readAt pgtype.Timestamp
		err = rows.Scan(&notification.Id, &notification.NotificationType, &notification.RecipientId, &notification.ActorId, &notification.ActorName, &notification.ObjectId, &notification.ObjectType, &notification.CreatedAt, &readAt, &notification.Payload)
		if err != nil {
			log.Error("fail to scan notification", zap.Error(err))
			return []domain.Notification{}, err
		}
		if readAt.Valid {
			notification.ReadAt = readAt.Time
		}
		res = append(res, notification)
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

func (obj *NotificationPostgres) GetById(ctx context.Context, id int) (domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select cluster, notification_type, recipient_id, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where id = $1 order by created_at;`
	args := []any{id}
	res := domain.Notification{Id: id}
	var readAt pgtype.Timestamp
	var cluster pgtype.Text
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&cluster, &res.NotificationType, &res.RecipientId, &res.ActorId, &res.ActorName, &res.ObjectId, &res.ObjectType, &res.CreatedAt, &readAt, &res.Payload)
	if err != nil {
		log.Error("fail to get notification", zap.Int("id", id), zap.Error(err))
		return domain.Notification{}, err
	}
	if readAt.Valid {
		res.ReadAt = readAt.Time
	}
	if cluster.Valid {
		res.Cluster = cluster.String
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

func (obj *NotificationPostgres) GetByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, notification_type, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where cluster is null and recipient_id = $1 order by created_at;`
	args := []any{userId}
	res := make([]domain.Notification, 0)
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("fail to get all notifications by user", zap.Int("userId", userId), zap.Error(err))
		return []domain.Notification{}, err
	}
	defer rows.Close()
	for rows.Next() {
		notification := domain.Notification{Cluster: "", RecipientId: userId}
		var readAt pgtype.Timestamp
		err = rows.Scan(&notification.Id, &notification.NotificationType, &notification.ActorId, &notification.ActorName, &notification.ObjectId, &notification.ObjectType, &notification.CreatedAt, &readAt, &notification.Payload)
		if err != nil {
			log.Error("fail to scan notification", zap.Error(err))
			return []domain.Notification{}, err
		}
		if readAt.Valid {
			notification.ReadAt = readAt.Time
		}
		res = append(res, notification)
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}
