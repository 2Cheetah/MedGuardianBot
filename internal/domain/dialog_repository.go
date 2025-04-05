package domain

type DialogRepository interface {
	CreateDialog(d Dialog) error
	GetActiveDialogByUserId(userId int64) (*Dialog, error)
	UpdateActiveDialog(d *Dialog) error
}
