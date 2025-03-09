package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/2Cheetah/MedGuardianBot/internal/app"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiToken := os.Getenv("API_TOKEN")
	dbPath := "users.db"

	app.Run(ctx, apiToken, dbPath)

}
