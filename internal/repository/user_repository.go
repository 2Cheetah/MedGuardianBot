package repository

import (
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

func (r *Repository) CreateUser(user *domain.User) error {
	q := "INSERT INTO users (id, first_name, last_name, username) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO NOTHING"
	_, err := r.db.Exec(
		q,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Username,
	)
	return err
}

func (r *Repository) GetUser(id int64) (*domain.User, error) {
	slog.Info("querying user with", "id", id)
	rows, err := r.db.Query("SELECT first_name, last_name, username  FROM users WHERE id = ?", id)
	slog.Info("found", "rows", rows)
	if err != nil {
		return nil, fmt.Errorf("couldn't query user from db, error %w", err)
	}
	defer rows.Close()

	var first_name string
	var last_name string
	var username string

	slog.Info("scanning rows...")
	if rows.Next() {
		if err := rows.Scan(&first_name, &last_name, &username); err != nil {
			slog.Warn("error while scanning rows", "error", err)
			return nil, fmt.Errorf("couldn't scan user from db, error %w", err)
		}
	} else {
		slog.Info("no user found with id", "id", id)
		return nil, nil
	}

	return &domain.User{
		FirstName: first_name,
		LastName:  last_name,
		ID:        id,
		Username:  username,
	}, nil
}
