package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/app"
)

// Send any text message to the bot after the bot has been started

func main() {
	// Set up signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Get configuration
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		slog.Error("API_TOKEN environment variable is required")
		os.Exit(1)
	}
	dbPath := "medguardian.db"

	// Initialize the application
	application, err := app.NewApp(apiToken, dbPath)
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start the application in a goroutine
	go func() {
		if err := application.Start(ctx); err != nil {
			slog.Error("error running application", "error", err)
			cancel() // Cancel context to trigger shutdown
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	slog.Info("Termination signal received, initiating shutdown...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Perform graceful shutdown
	if err := application.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error during shutdown", "error", err)
	}

	slog.Info("Application has been shut down gracefully")
}
