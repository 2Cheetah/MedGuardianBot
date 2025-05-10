package bot

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleCreateNotification(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info("CreateNotification handler triggered")
	// Stop any non-finished dialogs (non-IDLE state)

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	tb.NotificationFSMService.StartSession(userID, chatID)
	msg := "Alright! Let's create a notification. What is the schedule? For instance, \"daily 9am\""
	sendMsg(ctx, b, update.Message.Chat.ID, msg)
}
