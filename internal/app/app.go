// internal/app/app.go
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/2Cheetah/MedGuardianBot/internal/bot"
	"github.com/2Cheetah/MedGuardianBot/internal/repository"
	"github.com/2Cheetah/MedGuardianBot/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func Run(ctx context.Context, apiToken string, dbPath string) {
	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Run migrations before starting the app
	if err := runMigrations(dbPath); err != nil {
		log.Fatal("Migration failed:", err)
	}

	userRepo := repository.NewSQLiteUserRepository(db)
	userService := service.NewUserService(userRepo)
	telegramBot, err := bot.NewTelegramBot(apiToken, userService)
	if err != nil {
		log.Fatal("Failed to initialize bot:", err)
	}

	telegramBot.Start(ctx)
}

// runMigrations applies database migrations
func runMigrations(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("could not open DB: %w", err)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("could not create SQLite driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver,
	)
	if err != nil {
		return fmt.Errorf("migration initialization failed: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}
