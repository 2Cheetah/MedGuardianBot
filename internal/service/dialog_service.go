package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

type DialogService struct {
	repo                domain.DialogRepository
	scheduleProcessor   ScheduleProcessor
	notificationService NotificationService
}

type DialogCreateNotification struct {
	Schedule string `json:"schedule"`
	Text     string `json:"text"`
}

type ScheduleProcessor interface {
	ParseSchedule(schedule string) (string, error)
}

func NewDialogService(repo domain.DialogRepository, sp ScheduleProcessor, ns NotificationService) *DialogService {
	return &DialogService{
		repo:                repo,
		scheduleProcessor:   sp,
		notificationService: ns,
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
				s, _ := ds.scheduleProcessor.ParseSchedule(d.Context)
				s = strings.TrimSpace(s)
				slog.Info("parsed schedule", "crontab", s)
				dialogContext := DialogCreateNotification{
					Schedule: s,
				}
				contextData, err := json.Marshal(dialogContext)
				if err != nil {
					return "", fmt.Errorf("couldn't marshal dialog to create notification, error: %w", err)
				}
				contextString := string(contextData)

				dDB.Context = contextString
				dDB.UpdatedAt = time.Now().UTC()
				if err := ds.repo.UpdateActiveDialog(dDB); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}
				msg := "What do you want me notify you about? What is the notificaiton text?"
				return msg, nil
			} else {
				var dialogContext DialogCreateNotification
				if err := json.Unmarshal([]byte(dDB.Context), &dialogContext); err != nil {
					return "", fmt.Errorf("couldn't unmarshal context, error: %w", err)
				}
				dialogContext.Text = d.Context
				contextData, err := json.Marshal(dialogContext)
				if err != nil {
					return "", fmt.Errorf("couldn't marshal dialog to create notification, error: %w", err)
				}
				contextString := string(contextData)
				dDB.Context = contextString
				dDB.UpdatedAt = time.Now().UTC()
				dDB.State = domain.DialogStatusFinished
				if err := ds.repo.UpdateActiveDialog(dDB); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}

				notification := &domain.Notification{
					Status:   domain.NotificationStatusActive,
					UserID:   d.UserID,
					ChatID:   d.UserID,
					Text:     dialogContext.Text,
					Schedule: dialogContext.Schedule,
					Until:    time.Now(),
					Next:     time.Now(),
				}

				slog.Info("about to create notification", "notification", notification)

				if err = ds.notificationService.CreateNotification(notification); err != nil {
					return "", fmt.Errorf("couldn't create notification, error: %w", err)
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
