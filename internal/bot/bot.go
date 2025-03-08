// Package bot provides functionality to create and manage a Telegram bot
// that interacts with users and handles messages using the go-telegram-bot library.
package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/bot/handlers"
	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot *bot.Bot
}

func NewTelegramBot(apiToken string) (*TelegramBot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.HandleEcho),
	}
	bot, err := bot.New(apiToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create bot, error %w", err)
	}
	return &TelegramBot{
		bot: bot,
	}, nil
}

func (tb *TelegramBot) Start(ctx context.Context) {
	tb.RegisterHandler("/start", handlers.HandleStart)
	tb.bot.Start(ctx)
}

func (tb *TelegramBot) RegisterHandler(pattern string, h bot.HandlerFunc) {
	id := tb.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		h,
	)
	slog.Info("registered /start handler", "id", id)
}
