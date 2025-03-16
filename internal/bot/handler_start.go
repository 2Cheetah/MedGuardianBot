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

func (tb *TelegramBot) handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	id := update.Message.From.ID

	// check if user exists in DB and save user info
	u, err := tb.UserService.GetUser(id)
	if err != nil {
		slog.Warn("there was an error while trying to get a user", "error", err)
	}
	if u == nil {
		// send welcome message and action "typing"
		msg := "Welcome to the MedGuardian Bot! Let me register you first."
		sendMsg(ctx, b, update.Message.Chat.ID, msg)
		sendTypingAction(ctx, b, update.Message.Chat.ID)
		time.Sleep(1 * time.Second)

		// create a user in data store
		user := &domain.User{
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			Username:  update.Message.From.Username,
			ID:        id,
		}
		if err := tb.UserService.CreateUser(user); err != nil {
			slog.Error("couldn't create a user", "user data", user)
			msg := "There was a problem while registering you. Please, try again later."
			sendMsg(ctx, b, update.Message.Chat.ID, msg)
			return
		}

		// greet the user
		msg = fmt.Sprintf("Hi, %s! Was darf es sein?", user.Username)
		sendMsg(ctx, b, update.Message.Chat.ID, msg)
	} else {
		slog.Info("message from a user", "username", update.Message.From.Username)

		// send message to the user
		msg := fmt.Sprintf("Welcome back, %s! Was darf es sein?", u.Username)
		sendMsg(ctx, b, update.Message.Chat.ID, msg)
	}
}
