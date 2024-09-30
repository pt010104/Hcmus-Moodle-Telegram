package usecase

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (uc implUseCase) SendMessage(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(uc.chatID, message)
	msg.ParseMode = "HTML"
	_, err := uc.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
