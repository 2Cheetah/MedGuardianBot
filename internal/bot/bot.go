// Package bot provides functionality to create and manage a Telegram bot
// that interacts with users and handles messages using the go-telegram-bot library.
package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
	"github.com/2Cheetah/MedGuardianBot/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramBot struct {
	bot         *bot.Bot
	UserService *service.UserService
}

func NewTelegramBot(apiToken string, us *service.UserService) (*TelegramBot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handleEcho),
	}
	bot, err := bot.New(apiToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create bot, error %w", err)
	}
	return &TelegramBot{
		bot:         bot,
		UserService: us,
	}, nil
}

func (tb *TelegramBot) Start(ctx context.Context) {
	slog.Debug("registering handlers...")
	tb.RegisterHandler("/start", tb.handleStart)
	slog.Debug("starting bot...")
	tb.bot.Start(ctx)
}

// Stop gracefully stops the bot
func (tb *TelegramBot) Stop(ctx context.Context) error {
	slog.Debug("stopping bot...")
	// The bot will stop when its context is canceled
	// Any additional cleanup logic can be added here if needed
	return nil
}

func (tb *TelegramBot) RegisterHandler(pattern string, h bot.HandlerFunc) {
	id := tb.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		h,
	)
	slog.Debug("registered handler", "pattern", pattern, "id", id)
}

func handleEcho(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil {
		slog.Warn("couldn't send message", "error", err)
	}
}

func (tb *TelegramBot) handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	// get user info
	id := update.Message.From.ID
	user := &domain.User{
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Username:  update.Message.From.Username,
		ID:        id,
	}

	slog.Info("message from", "user", user)

	// check if user exists in DB and save user info
	u, err := tb.UserService.GetUser(id)
	if err != nil {
		slog.Warn("there was an error while trying to get a user", "error", err)
	}
	if u == nil {
		if err := tb.UserService.CreateUser(user); err != nil {
			slog.Warn("couldn't create a user", "user data", user)
		}
	}

	// send message to the user
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Congrats!111 You selected \"/start\" command.",
	}); err != nil {
		slog.Warn("couldn't send message", "error", err)
	}
}
