package repository

import (
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

func (r *Repository) CreateDialog(userID int64, command string) error {
	slog.Info("creating dialog", "userID", userID, "command", command)
	q := "INSERT INTO dialogs (user_id, command) VALUES (?, ?)"
	_, err := r.db.Exec(
		q,
		userID,
		command,
	)
	if err != nil {
		return fmt.Errorf("couldn't create a dialog, error: %w", err)
	}
	return nil
}

func (r *Repository) GetActiveDialogByUserId(userId int64) (*domain.Dialog, error) {
	q := "SELECT * FROM dialogs WHERE user_id = ? AND state = 'CREATED'"
	rows, err := r.db.Query(q, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get active dialogs by userId, error: %w", err)
	}
	defer rows.Close()

	var dialog domain.Dialog

	if rows.Next() {
		if err := rows.Scan(&dialog); err != nil {
			return nil, fmt.Errorf("couldn't scan DB row to dialog, error: %w", err)
		}
	} else {
		return nil, nil
	}

	if !rows.Next() {
		return &dialog, nil
	} else {
		return nil, fmt.Errorf("more than one STARTED dialogs found for userId: %d", userId)
	}
}

func (r *Repository) UpdateActiveDialog(dialog *domain.Dialog) error {
	q := "UPDATE dialogs SET state = ?, updated_at = ?, context = ? WHERE state = 'STARTED' AND user_id = ?"
	res, err := r.db.Exec(
		q,
		dialog.State,
		dialog.UpdatedAt,
		dialog.Context,
		dialog.UserID,
	)
	if err != nil {
		return fmt.Errorf("couldn't update active dialog, dialog: %v, error: %w", *dialog, err)
	}
	if n, _ := res.RowsAffected(); n > 1 {
		slog.Error("more than one row was updated by UpdateActiveDialog sql")
	}
	return nil
}
