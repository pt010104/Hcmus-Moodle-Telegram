package usecase

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (uc implUseCase) CommandHandler(ctx context.Context, message *tgbotapi.Message) error {
	if message == nil || message.Text == "" {
		return nil
	}

	msgText := message.Text
	uc.l.Info(ctx, "Received message", msgText)

	switch {
	case msgText == "/ls":
		return uc.handleListCourses(ctx, message)
	case msgText == "/ld":
		return uc.handleListDeadlines(ctx, message)
	case strings.HasPrefix(msgText, "/cd"):
		return uc.handleCourseDeadlines(ctx, message)
	default:
		return uc.sendTextMessage(ctx, message.Chat.ID, "Unknown command")
	}
}
