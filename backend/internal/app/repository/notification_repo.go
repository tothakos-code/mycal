package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
)

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (n *NotificationRepo) CreateNotification(ctx context.Context, notification models.Notification) error {
	query := `INSERT INTO notification (id, user_id, event_id, title, message, shown_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := n.db.ExecContext(ctx, query,
		notification.ID, notification.UserID, notification.EventID,
		notification.Title, notification.Message, notification.ShownAt)
	return err
}

func (n *NotificationRepo) GetNotificationsByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	query := `SELECT * FROM notification WHERE user_id = $1`
	rows, err := n.db.QueryContext(ctx, query, userID)
  if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.EventID,
			&notif.Title, &notif.Message, &notif.ShownAt, &notif.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}
