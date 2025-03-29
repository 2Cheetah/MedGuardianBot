package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tb *TelegramBot) handleArbitraryText(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	context := update.Message.Text
	dialog := &domain.Dialog{
		UserID:  userID,
		Context: context,
	}
	msg, err := tb.DialogService.HandleDialog(dialog)
	if err != nil {
		slog.Error("couldn't handleArbitraryText", "error", err)
		msg := fmt.Sprintf("Error while handling text:\n%s", update.Message.Text)
		sendMsg(ctx, b, update.Message.Chat.ID, msg)
		return
	}
	sendMsg(ctx, b, update.Message.Chat.ID, msg)
}
