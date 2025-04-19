// internal/app/app.go
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/bot"
	crontabninja "github.com/2Cheetah/MedGuardianBot/internal/crontabNinja"
	"github.com/2Cheetah/MedGuardianBot/internal/repository"
	"github.com/2Cheetah/MedGuardianBot/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// App encapsulates all application components and manages their lifecycle
type App struct {
	db          *sql.DB
	telegramBot *bot.TelegramBot
}

// NewApp initializes a new application instance with all dependencies
func NewApp(apiToken string, dbPath string) (*App, error) {
	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Run migrations before starting the app
	if err := runMigrations(dbPath); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	repo := repository.NewRepository(db)
	userService := service.NewUserService(repo)
	scheduleProcessor := crontabninja.NewClient("https://cronly.app/api/ai/generate")
	notificationFSMService := service.NewNotificationFSMService(scheduleProcessor)
	notificationService := service.NewNotificationService(repo)
	dialogService := service.NewDialogService(repo, scheduleProcessor, *notificationService)
	telegramBot, err := bot.NewTelegramBot(apiToken, userService, notificationFSMService, dialogService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot: %w", err)
	}

	return &App{
		db:          db,
		telegramBot: telegramBot,
	}, nil
}

// Start launches all application components and blocks until context is canceled
func (a *App) Start(ctx context.Context) error {
	slog.Info("Starting application...")

	// Start the Telegram bot
	go a.telegramBot.Start(ctx)

	// Wait for context cancellation (termination signal)
	<-ctx.Done()
	slog.Info("Received termination signal")

	return nil
}

// Shutdown performs a graceful shutdown of all application components
func (a *App) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down application...")

	// Explicitly stop the Telegram bot
	if err := a.telegramBot.Stop(ctx); err != nil {
		slog.Warn("error stopping Telegram bot", "error", err)
		// Continue with shutdown even if there's an error
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	slog.Info("Application shutdown complete")
	return nil
}

// Run is kept for backward compatibility
func Run(ctx context.Context, apiToken string, dbPath string) {
	app, err := NewApp(apiToken, dbPath)
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start the app
	go func() {
		if err := app.Start(ctx); err != nil {
			slog.Error("Error running app", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Perform graceful shutdown
	if err := app.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error during shutdown", "error", err)
	}
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

	slog.Debug("Migrations applied successfully!")
	return nil
}
