package domain

type NotificationState string
type NotificationEvent string

const (
	StateWaitingSchedule NotificationState = "waiting_schedule"
	StateWaitingUntil    NotificationState = "waiting_until"
	StateWaitingText     NotificationState = "waiting_text"
	StateCreated         NotificationState = "created"
	StateEnded           NotificationState = "ended"
	StateCancelled       NotificationState = "cancelled"
	StateRemoved         NotificationState = "removed"

	EventCreate          NotificationEvent = "notification_created"
	EventScheduleReceive NotificationEvent = "schedule_received"
	EventTextReceive     NotificationEvent = "text_received"
	EventCancel          NotificationEvent = "cancelled"
)

type NotificationFSM struct {
	UserID              int64
	ChatID              int64
	State               NotificationState
	PartialNotification Notification
}
