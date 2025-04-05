package bot

import (
	"context"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleCreateNotification(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info("CreateNotification handler triggered")
	// Stop any non-finished dialogs (non-IDLE state)

	msg := "Alright! Let's create a notification. What is the schedule? For instance, \"daily 9am\""
	sendMsg(ctx, b, update.Message.Chat.ID, msg)

	slog.Info("calling dialog service to create a dialog")
	d := domain.Dialog{
		UserID:  update.Message.From.ID,
		ChatID:  update.Message.Chat.ID,
		Command: "create_notification",
	}
	if err := tb.DialogService.CreateDialog(d); err != nil {
		slog.Error("coulnd't create dialog", "dialog", d, "error", err)
	}
}
