package domain

import "time"

type NotificationStatus string

const (
	NotificationStatusActive   NotificationStatus = "ACTIVE"
	NotificationStatusInactive NotificationStatus = "INACTIVE"
	NotificationStatusFinished NotificationStatus = "FINISHED"
)

type Notification struct {
	ID        int64
	Status    NotificationStatus
	UserID    int64
	ChatID    int64
	Text      string
	Schedule  string
	CreatedAt time.Time
	Until     time.Time
	Next      time.Time
}
