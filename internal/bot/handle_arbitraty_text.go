package bot

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleArbitraryText(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Check dialogs states, if any is in STARTED state
	userID := update.Message.From.ID
	slog.Info("checking if user has STARTED dialogs", "userID", userID)
	dialog, err := tb.DialogService.GetActiveDialogByUserId(userID)
	if err != nil {
		slog.Error("couldn't GetActiveDialogByUserId", "error", err)
		return
	}
	slog.Info("received from GetActiveDialogByUserId", "dialog", dialog)
	if dialog != nil {
		switch dialog.Command {
		case "create_notification":
			dialog.Context = update.Message.Text
			dialog.UpdatedAt = time.Now().UTC()
			if err := tb.DialogService.UpdateActiveDialog(dialog); err != nil {
				slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
			}
			return
		}
	}
	// If no dialogs in STARTED state, respond with help message
	msg := "No active dialogs found. Select a command."
	sendMsg(ctx, b, update.Message.Chat.ID, msg)
}
