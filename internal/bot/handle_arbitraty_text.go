package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
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
			if dialog.Context == "" {
				dialog.Context = "schedule: " + update.Message.Text
				dialog.UpdatedAt = time.Now().UTC()
				if err := tb.DialogService.UpdateActiveDialog(dialog); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}
				msg := "What do you want me notify you about? What is the notificaiton text?"
				sendMsg(ctx, b, update.Message.Chat.ID, msg)
			} else {
				dialog.Context += " text: " + update.Message.Text
				dialog.UpdatedAt = time.Now().UTC()
				dialog.State = domain.DialogStatusFinished
				if err := tb.DialogService.UpdateActiveDialog(dialog); err != nil {
					slog.Error("couldn't UpdateActiveDialog from handle_arbitraty_text.go", "error", err)
				}
				msg := fmt.Sprintf("Success! Notification created! %s", dialog.Context)
				sendMsg(ctx, b, update.Message.Chat.ID, msg)
			}
			return
		}
	}
	// If no dialogs in STARTED state, respond with help message
	msg := "No active dialogs found. Select a command."
	sendMsg(ctx, b, update.Message.Chat.ID, msg)
}
