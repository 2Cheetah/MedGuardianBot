package repository

import (
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

func (r *Repository) CreateDialog(d domain.Dialog) error {
	slog.Info("creating dialog", "userID", d.UserID, "command", d.Command)
	q := "INSERT INTO dialogs (user_id, chat_id, command) VALUES (?, ?, ?)"
	_, err := r.db.Exec(
		q,
		d.UserID,
		d.ChatID,
		d.Command,
	)
	if err != nil {
		return fmt.Errorf("couldn't create a dialog, error: %w", err)
	}
	return nil
}

func (r *Repository) GetActiveDialogByUserId(userID int64) (*domain.Dialog, error) {
	q := "SELECT * FROM dialogs WHERE user_id = ? AND state = ?"
	rows, err := r.db.Query(q, userID, domain.DialogStatusStarted)
	if err != nil {
		return nil, fmt.Errorf("couldn't get active dialogs by userID, error: %w", err)
	}
	defer rows.Close()

	var dialog domain.Dialog

	if rows.Next() {
		if err := rows.Scan(
			&dialog.ID,
			&dialog.UserID,
			&dialog.ChatID,
			&dialog.State,
			&dialog.CreatedAt,
			&dialog.UpdatedAt,
			&dialog.Command,
			&dialog.Context,
		); err != nil {
			return nil, fmt.Errorf("couldn't scan DB row to dialog, error: %w", err)
		}
	} else {
		return nil, nil
	}

	if !rows.Next() {
		return &dialog, nil
	} else {
		return nil, fmt.Errorf("more than one STARTED dialogs found for userID: %d", userID)
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
