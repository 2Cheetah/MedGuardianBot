package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handleArbitraryText(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Check dialog state, if any is not in IDLE continue

	// If all dialogs in IDLE state, send list of available actions
}
