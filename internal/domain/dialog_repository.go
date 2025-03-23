package domain

type DialogRepository interface {
	CreateDialog(userId int64, command string) error
	GetActiveDialogByUserId(userId int64) (*Dialog, error)
	UpdateActiveDialog(d *Dialog) error
}
