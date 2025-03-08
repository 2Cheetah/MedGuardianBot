// Package bot provides functionality to create and manage a Telegram bot
// that interacts with users and handles messages using the go-telegram-bot library.
package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramBot struct {
	bot *bot.Bot
}

func NewTelegramBot(apiToken string) (*TelegramBot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handleEcho),
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
	// // Register handler here
	// tb.bot.RegisterHandler(
	// 	bot.HandlerTypeMessageText,
	// 	"",
	// 	bot.MatchTypeContains,
	// 	handleEcho,
	// )
	tb.bot.Start(ctx)
}

func handleEcho(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil {
		log.Fatalf("couldn't send message, err %v", err)
	}
}
