package domain

type NotificationRepository interface {
	CreateNotification(n *Notification) error
	GetNotificationsByStatus(status NotificationStatus) ([]*Notification, error)
	GetActiveNotificationsByUserID(userID int64) ([]*Notification, error)
}
