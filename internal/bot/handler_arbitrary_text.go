package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleArbitraryText(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	input := update.Message.Text
	msg, err := tb.NotificationFSMService.HandleInput(userID, input)
	if err != nil {
		sendMsg(ctx, tb.bot, chatID, "Something went wrong")
		return
	}
	sendMsg(ctx, tb.bot, chatID, msg)
}
