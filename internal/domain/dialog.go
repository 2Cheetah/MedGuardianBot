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
	ID        int64         `db:"id"`
	UserID    int64         `db:"user_id"`
	ChatID    int64         `db:"chat_id"`
	State     DialogState   `db:"state"`
	Command   DialogCommand `db:"command"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	Context   string        `db:"context"`
}
