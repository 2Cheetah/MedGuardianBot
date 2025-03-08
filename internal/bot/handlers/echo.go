package handlers

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleEcho(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil {
		log.Fatalf("couldn't send message, err %v", err)
	}
}
