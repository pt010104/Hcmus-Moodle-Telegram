package telegram

import (
	"context"

	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/calendar"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UseCase interface {
	SendMessage(ctx context.Context, message string) error
	CommandHandler(ctx context.Context, message *tgbotapi.Message) error

	SetCalendarUC(calendarUC calendar.UseCase)
}
