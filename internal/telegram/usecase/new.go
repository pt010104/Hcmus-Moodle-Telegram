package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pt010104/Hcmus-Moodle-Telegram/internal/telegram"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/log"
)

type implUseCase struct {
	l      log.Logger
	bot    *tgbotapi.BotAPI
	chatID int64
}

func New(
	l log.Logger,
	bot *tgbotapi.BotAPI,
	chatID int64,
) telegram.UseCase {
	return &implUseCase{
		l:      l,
		bot:    bot,
		chatID: chatID,
	}
}
