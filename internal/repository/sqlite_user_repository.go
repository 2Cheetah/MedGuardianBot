// internal/repository/sqlite_user_repository.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) CreateUser(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, first_name, last_name, username) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO NOTHING",
		user.ID, user.FirstName, user.LastName, user.Username)
	return err
}

func (r *SQLiteUserRepository) GetUser(id int64) (*domain.User, error) {
	rows, err := r.db.Query("SELECT first_name, last_name, username  FROM users WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, fmt.Errorf("couldn't query user from db, error %v", err)
	}
	defer rows.Close()

	rows.Next()
	var first_name string
	var last_name string
	var username string
	if err = rows.Scan(&first_name, &last_name, &username); err != nil {
		return nil, fmt.Errorf("couldn't scan user from db, error %v", err)
	}

	return &domain.User{
		FirstName: first_name,
		LastName:  last_name,
		ID:        id,
		Username:  username,
	}, nil
}
