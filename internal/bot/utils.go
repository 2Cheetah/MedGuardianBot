package bot

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func sendMsg(ctx context.Context, b *bot.Bot, id int64, msg string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: id,
		Text:   msg,
	}); err != nil {
		slog.Warn("couldn't send message", "error", err)
	}
}

func sendTypingAction(ctx context.Context, b *bot.Bot, id int64) {
	if _, err := b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: id,
		Action: models.ChatActionTyping,
	}); err != nil {
		slog.Warn("couldn't send aciton \"typing\"", "error", err)
	}
}
