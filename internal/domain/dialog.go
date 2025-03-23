package domain

import "time"

type DialogState string

const (
	DialogStatusStarted   DialogState = "STARTED"
	DialogStatusFinished  DialogState = "FINISHED"
	DialogStatusCancelled DialogState = "CANCELLED"
)

type DialogCommand string

const (
	DialogCommandCreateNotification DialogCommand = "create_notification"
)

type Dialog struct {
	ID        int64
	UserID    int64
	State     DialogState
	Command   DialogCommand
	CreatedAt time.Time
	UpdatedAt time.Time
	Context   string
}
