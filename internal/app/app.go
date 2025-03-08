// internal/app/app.go
package app

import (
	"context"
	"log"

	"github.com/2Cheetah/MedGuardianBot/internal/bot"
)

func Run(apiToken string, ctx context.Context) {
	telegramBot, err := bot.NewTelegramBot(apiToken)
	if err != nil {
		log.Fatal("Failed to initialize bot:", err)
	}

	telegramBot.Start(ctx)
}
