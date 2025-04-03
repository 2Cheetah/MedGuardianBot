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
	bot           *bot.Bot
	UserService   *service.UserService
	DialogService *service.DialogService
}

func NewTelegramBot(apiToken string, us *service.UserService, ds *service.DialogService) (*TelegramBot, error) {
	tb := &TelegramBot{
		UserService:   us,
		DialogService: ds,
	}
	opts := []bot.Option{
		bot.WithDefaultHandler(tb.handleArbitraryText),
	}
	bot, err := bot.New(apiToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create bot, error %w", err)
	}
	tb.bot = bot
	return tb, nil
}

func (tb *TelegramBot) Start(ctx context.Context) {
	slog.Debug("registering handlers...")
	tb.RegisterHandlerExactMatch("/start", tb.handleStart)
	tb.RegisterHandlerExactMatch("/create_notification", tb.handleCreateNotification)
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

func (tb *TelegramBot) RegisterHandlerExactMatch(pattern string, h bot.HandlerFunc) {
	id := tb.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		h,
	)
	slog.Debug("registered handler", "pattern", pattern, "id", id)
}

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
