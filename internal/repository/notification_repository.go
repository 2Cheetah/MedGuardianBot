package repository

import (
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

func (r *Repository) CreateNotification(n *domain.Notification) error {
	slog.Info("creating notification", "notification", n)
	q := "INSERT INTO notifications (status, user_id, chat_id, text, schedule, until, next) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(
		q,
		n.Status,
		n.UserID,
		n.ChatID,
		n.Text,
		n.Schedule,
		n.Until,
		n.Next,
	)
	if err != nil {
		return fmt.Errorf("couldn't create a notification, error %w", err)
	}
	return nil
}

func (r *Repository) GetNotificationsByStatus(status domain.NotificationStatus) ([]*domain.Notification, error) {
	q := "SELECT id, status, user_id, chat_id, text, schedule, created_at, to, next FROM notifications WHERE status = ?"
	rows, err := r.db.Query(q, status)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute query to get active notifications, error: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification

	for rows.Next() {
		var notificaiton domain.Notification
		if err := rows.Scan(
			&notificaiton.ID,
			&notificaiton.Status,
			&notificaiton.UserID,
			&notificaiton.ChatID,
			&notificaiton.Text,
			&notificaiton.Schedule,
			&notificaiton.CreatedAt,
			&notificaiton.Until,
			&notificaiton.Next,
		); err != nil {
			return nil, fmt.Errorf("couldn't scan rows, error: %w", err)
		}
		notifications = append(notifications, &notificaiton)
	}
	return notifications, nil
}

func (r *Repository) GetActiveNotificationsByUserID(userID int64) ([]*domain.Notification, error) {
	q := "SELECT * FROM notifications WHERE user_id = ? AND status = ?"
	rows, err := r.db.Query(q, userID, domain.NotificationStatusActive)
	if err != nil {
		return nil, fmt.Errorf("couldn't get notifications, error: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification

	for rows.Next() {
		var notification domain.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.Status,
			&notification.UserID,
			&notification.ChatID,
			&notification.Text,
			&notification.Schedule,
			&notification.CreatedAt,
			&notification.Until,
			&notification.Next,
		); err != nil {
			return nil, fmt.Errorf("couldn't scan DB row to notification, error: %w", err)
		}
		notifications = append(notifications, &notification)
	}
	return notifications, nil
}
