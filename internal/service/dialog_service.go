package service

import (
	"errors"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

type DialogService struct {
	repo domain.DialogRepository
}

func NewDialogService(repo domain.DialogRepository) *DialogService {
	return &DialogService{
		repo: repo,
	}
}

func (ds *DialogService) CreateDialog(dialog *domain.Dialog) error {
	slog.Info("calling repository to create a dialog", "dialog", dialog)
	return ds.repo.CreateDialog(dialog.UserID, string(dialog.Command))
}

func (ds *DialogService) GetActiveDialogByUserId(userID int64) (*domain.Dialog, error) {
	if userID <= 0 {
		return nil, errors.New("userID can't be zero or negative")
	}
	return ds.repo.GetActiveDialogByUserId(userID)
}

func (ds *DialogService) UpdateActiveDialog(dialog *domain.Dialog) error {
	return ds.repo.UpdateActiveDialog(dialog)
}
