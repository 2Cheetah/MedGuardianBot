package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

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

// ArbitraryDialog handles dialogs before completing flow for commands
func (ds *DialogService) HandleDialog(d *domain.Dialog) (string, error) {
	// Check if any dialog for the user is in STARTED state
	slog.Info("checking if user has STARTED dialogs", "userID", d.UserID)
	dDB, err := ds.repo.GetActiveDialogByUserId(d.UserID)
	if err != nil {
		slog.Error("couldn't GetActiveDialogByUserId", "error", err)
		return "", fmt.Errorf("couldn't GetActiveDialogByUserId, error: %w", err)
	}
	slog.Info("received from GetActiveDialogByUserId", "dialog", dDB)
	if dDB != nil {
		switch dDB.Command {
		case "create_notification":
			if dDB.Context == "" {
				dDB.Context = "schedule: " + d.Context
				dDB.UpdatedAt = time.Now().UTC()
				if err := ds.repo.UpdateActiveDialog(dDB); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}
				msg := "What do you want me notify you about? What is the notificaiton text?"
				return msg, nil
			} else {
				dDB.Context += " text: " + d.Context
				dDB.UpdatedAt = time.Now().UTC()
				dDB.State = domain.DialogStatusFinished
				if err := ds.repo.UpdateActiveDialog(dDB); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}
				msg := fmt.Sprintf("Success! Notification created! %s", dDB.Context)
				return msg, nil
			}
		default:
			msg := fmt.Sprintf("Found STARTED dialog for command %s", dDB.Command)
			return msg, nil
		}
	}
	// If no dialogs in STARTED state, respond with help message
	msg := "No active dialogs found. Select a command."
	return msg, nil
}
