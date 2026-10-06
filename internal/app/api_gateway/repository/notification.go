package repository

import (
	"context"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	GetById(ctx context.Context, id int) (domain.Notification, error)
	GetByUser(ctx context.Context, userId int, limit int) ([]domain.Notification, error)
	GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error)
}

type NotificationPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationPostgres(db *pgxpool.Pool) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (obj *NotificationPostgres) GetById(ctx context.Context, id int) (domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select cluster, notification_type, recipient_id, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where id = $1;`
	args := []any{id}
	res := domain.Notification{Id: id}
	start := time.Now()
	err := obj.db.QueryRow(ctx, query, args...).Scan(&res.Cluster, &res.NotificationType, &res.RecipientId, &res.ActorId, &res.ActorName, &res.ObjectId, &res.ObjectType, &res.CreatedAt, &res.ReadAt, &res.Payload)
	if err != nil {
		log.Error("failed to get notification", zap.Error(err))
		return domain.Notification{}, err
	}
	duration := time.Since(start)
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	log.Info("Query executed")
	return res, nil
}

func (obj *NotificationPostgres) GetByUser(ctx context.Context, userId int, limit int) ([]domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, cluster, notification_type, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where recipient_id = $1 order by notification.created_at limit $2;`
	args := []any{userId, limit}
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get notifications by user", zap.Error(err))
		return []domain.Notification{}, err
	}
	duration := time.Since(start)
	defer rows.Close()
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	res := make([]domain.Notification, 0)
	for rows.Next() {
		notification := domain.Notification{RecipientId: userId}
		err = rows.Scan(&notification.Id, &notification.Cluster, &notification.NotificationType, &notification.ActorId, &notification.ActorName, &notification.ObjectId, &notification.ObjectType, &notification.CreatedAt, &notification.ReadAt, &notification.Payload)
		if err != nil {
			log.Error("failed to scan notification", zap.Error(err))
			return []domain.Notification{}, err
		}
		res = append(res, notification)
	}
	log.Info("Query executed")
	return res, nil
}

func (obj *NotificationPostgres) GetAllByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	query := `select id, cluster, notification_type, actor_id, actor_name, object_id, object_type, created_at, read_at, payload from notification where recipient_id = $1 order by notification.created_at;`
	args := []any{userId}
	start := time.Now()
	rows, err := obj.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to get all notifications by user", zap.Error(err))
		return []domain.Notification{}, err
	}
	duration := time.Since(start)
	defer rows.Close()
	log = logger.ModifyLoggerWithDBQuery(log, query, args, duration)
	res := make([]domain.Notification, 0)
	for rows.Next() {
		notification := domain.Notification{RecipientId: userId}
		err = rows.Scan(&notification.Id, &notification.Cluster, &notification.NotificationType, &notification.ActorId, &notification.ActorName, &notification.ObjectId, &notification.ObjectType, &notification.CreatedAt, &notification.ReadAt, &notification.Payload)
		if err != nil {
			log.Error("failed to scan notification", zap.Error(err))
			return []domain.Notification{}, err
		}
		res = append(res, notification)
	}
	log.Info("Query executed")
	return res, nil
}
