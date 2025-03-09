// internal/repository/sqlite_user_repository.go
package repository

import (
	"database/sql"

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
