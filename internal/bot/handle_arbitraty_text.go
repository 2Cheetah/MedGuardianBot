package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleArbitraryText(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Check dialogs states, if any is in STARTED state

	// If all dialogs in S state, send list of available actions
}
