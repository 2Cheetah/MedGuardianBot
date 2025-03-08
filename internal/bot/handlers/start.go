package handlers

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Congrats!111 You selected \"start\" command.",
	})
	if err != nil {
		log.Fatalf("couldn't send message, err %v", err)
	}
}
